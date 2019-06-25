/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package monitoring

import (
	"sort"
	"time"

	"github.com/nalej/derrors"

	"github.com/nalej/grpc-edge-inventory-proxy-go"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/conversions"

	"github.com/nalej/inventory-manager/internal/pkg/server/contexts"

	"github.com/rs/zerolog/log"
)

type Manager struct {
	proxyClient grpc_edge_inventory_proxy_go.EdgeControllerProxyClient
	assetsClient grpc_inventory_go.AssetsClient
	controllersClient grpc_inventory_go.ControllersClient
}

func NewManager(proxyClient grpc_edge_inventory_proxy_go.EdgeControllerProxyClient, assetsClient grpc_inventory_go.AssetsClient, controllersClient grpc_inventory_go.ControllersClient) *Manager {
	return &Manager{
		proxyClient: proxyClient,
		assetsClient: assetsClient,
		controllersClient: controllersClient,
	}
}

const edgeControllerAliveTimeout = 600

func (m *Manager) ListMetrics(selector *grpc_inventory_manager_go.AssetSelector) (*grpc_inventory_manager_go.MetricsList, error) {
	// Get a selector for each relevant Edge Controller
	selectors, derr := m.getSelectors(selector)
	if derr != nil {
		return nil, conversions.ToGRPCError(derr)
	}

	metrics := make(map[string]bool)

	// Create a request for each Edge Controller and execute
	for _, proxyRequest := range(selectors) {
		ecId := proxyRequest.GetEdgeControllerId()
		log.Debug().Interface("request", proxyRequest).Msg("proxy request for ListMetrics")
		ctx, _ := contexts.ProxyContext()
		list, err := m.proxyClient.ListMetrics(ctx, proxyRequest)
		if err != nil {
			// We still want to query to working edge controllers
			log.Warn().Str("edge-controller-id", ecId).Err(err).Msg("failed calling ListMetrics")
			continue
		}
		for _, metric := range(list.GetMetrics()) {
			metrics[metric] = true
		}
	}

	// Unify the results
	metricsList := make([]string, 0, len(metrics))
	for metric := range(metrics) {
		metricsList = append(metricsList, metric)
	}

	return &grpc_inventory_manager_go.MetricsList{
		Metrics: metricsList,
	}, nil
}

func (m *Manager) QueryMetrics(request *grpc_inventory_manager_go.QueryMetricsRequest) (*grpc_inventory_manager_go.QueryMetricsResult, error) {
	// Get a selector for each relevant Edge Controller
	selectors, derr := m.getSelectors(request.GetAssets())
	if derr != nil {
		return nil, conversions.ToGRPCError(derr)
	}

	aggregationType := request.GetAggregation()

	// Results is a mapping from metric to values, where values is a mapping
	// from timestamp to value and count. This last mapping is needed for merging
	// results from multiple edge controllers. We will convert to one
	// QueryMetricsResult to return afterwards
	results := make(map[string]map[int64]*grpc_inventory_manager_go.QueryMetricsResult_Value)

	// Request for each Edge Controller and execute
	for _, selector := range(selectors) {
		proxyRequest := &grpc_inventory_manager_go.QueryMetricsRequest{
			Assets: selector,
			Metrics: request.GetMetrics(),
			TimeRange: request.GetTimeRange(),
			Aggregation: request.GetAggregation(),
		}
		ecId := selector.GetEdgeControllerId()
		log.Debug().Interface("request", proxyRequest).Msg("proxy request for QueryMetrics")
		ctx, _ := contexts.ProxyContext()
		result, err := m.proxyClient.QueryMetrics(ctx, proxyRequest)

		// Optimization when we're querying only a single EC, in which
		// case we can also return any errors
		if len(selectors) == 1 {
			log.Debug().Str("ecid", ecId).Msg("querying single edge controller - skipping merging")
			return result, err
		}

		if err != nil {
			// We still want to query to working edge controllers
			log.Warn().Str("edge-controller-id", ecId).Err(err).Msg("failed calling QueryMetrics")
			continue
		}

		// Loop over all returned metrics
		for metric, assetMetrics := range(result.GetMetrics()) {
			// Currently, we always have a single assetMetric
			assetMetricList := assetMetrics.GetMetrics()
			if len(assetMetricList) == 0 {
				continue
			}
			if len(assetMetricList) > 1 {
				log.Warn().Msg("received query result for more than one individual asset - not supported")
			}

			assetMetric := assetMetricList[0]
			assetMetricValues := assetMetric.GetValues()

			// Create or retrieve map for this metric
			values, found := results[metric]
			if !found {
				values = make(map[int64]*grpc_inventory_manager_go.QueryMetricsResult_Value, len(assetMetricValues))
				results[metric] = values
			}

			log.Debug().Str("ecid", ecId).Str("metric", metric).Int("count", len(assetMetricValues)).Msg("storing metrics")
			for _, assetMetricValue := range assetMetricValues {
				// Create or add to value for this timestamp
				timestamp := assetMetricValue.GetTimestamp()
				value, found := values[timestamp]
				if !found {
					values[timestamp] = assetMetricValue
				} else {
					value.Value += assetMetricValue.Value
					value.AssetCount += assetMetricValue.AssetCount
				}
			}
		}
	}

	// Unify the results
	metricResults := make(map[string]*grpc_inventory_manager_go.QueryMetricsResult_AssetMetrics, len(results))
	for metric, valueMap := range(results) {
		// Make a list sorted timestamps
		keys := make([]int64, len(valueMap))
		for key := range(valueMap) {
			keys = append(keys, key)
		}
		// Sort int64
		sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

		log.Debug().Str("metric", metric).Int("count", len(keys)).Msg("aggregating metrics")
		// Create final value list, applying aggregation if needed
		values := make([]*grpc_inventory_manager_go.QueryMetricsResult_Value, 0, len(keys))
		for _, key := range(keys) {
			value := valueMap[key]
			switch aggregationType {
			case grpc_inventory_manager_go.AggregationType_SUM:
				// Nothing to do - it's already a sum
			case grpc_inventory_manager_go.AggregationType_AVG:
				value.Value = value.Value / value.AssetCount
			default:
				return nil, derrors.NewInvalidArgumentError("unknown aggregation type").WithParams(aggregationType)
			}
			values = append(values, value)
		}

		metricResults[metric] = &grpc_inventory_manager_go.QueryMetricsResult_AssetMetrics{
			Metrics: []*grpc_inventory_manager_go.QueryMetricsResult_AssetMetricValues{
				&grpc_inventory_manager_go.QueryMetricsResult_AssetMetricValues{
					// We don't have an asset id - if we had a single
					// asset, we would have had a single EC, and we
					// would have returned already.
					Values: values,
					Aggregation: aggregationType,
				},
			},
		}
	}

	result := &grpc_inventory_manager_go.QueryMetricsResult{
		Metrics: metricResults,
	}
	return result, nil
}

// GetSelectors will turn a single asset selector into a selector per
// edge controller, so we can sent the request to each edge controller
// that needs it. This is done by creating a list of assets from either
// AssetIds in the source selector or, if none are provided, by retrieving
// all Assets for an OrganizationId or EdgeControllerId. We then filter
// that list to remove the Assets that don't match the groups and labels and
// sort them out by EdgeControllerId.
// If there are no group/label filters, we can just create a selector
// for each edge controller without specific assets, as that will select
// all assets available on an Edge Controller without having to communicate
// a long list.
// We do not filter out disabled assets, as we assume that disabled assets
// ("show" is false) are not sending monitoring data anymore anyway. We still
// want to include when retrieving historic data.
func (m *Manager) getSelectors(selector *grpc_inventory_manager_go.AssetSelector) (map[string]*grpc_inventory_manager_go.AssetSelector, derrors.Error) {
	selectors := make(map[string]*grpc_inventory_manager_go.AssetSelector)

	orgId := selector.GetOrganizationId()
	ecId := selector.GetEdgeControllerId()
	assetIds := selector.GetAssetIds()

	// We'll always create a set of asset ids to filter, as we need to
	// filter out assets that are deleted/disabled. If we already have
	// a set of assets, we still need to filter out disabled assets and
	// disabled edge controllers.

	if len(assetIds) > 0 {
		// If we have explicit assets, that's the minimum set we start from
		for _, id := range(assetIds) {
			ctx, cancel := contexts.InventoryContext()
			// Calling cancel manually to avoid stacking up a lot of defers

			asset, err := m.assetsClient.Get(ctx, &grpc_inventory_go.AssetId{
				OrganizationId: orgId,
				AssetId: id,
			})
			if err != nil {
				cancel()
				return nil, derrors.NewUnavailableError("unable to retrieve asset information", err).WithParams(id)
			}
			if selectedAsset(asset, selector) {
				addAsset(selectors, asset)
			}
			cancel()
		}
	} else if len(selector.GetLabels()) == 0 && len(selector.GetGroupIds()) == 0 {
		// Make a selector for each Edge Controller, without explicit assets
		if ecId != "" {
			// No further selectors and ecId means we just need the
			// already existing selector
			selectors[ecId] = selector
		} else {
			// Selector for each Edge Controller in Organization
			ctx, cancel := contexts.InventoryContext()
			defer cancel()

			ecList, err := m.controllersClient.List(ctx, &grpc_organization_go.OrganizationId{
				OrganizationId: orgId,
			})
			if err != nil {
				return nil, derrors.NewUnavailableError("unable to retrieve edge controllers", err).WithParams(orgId)
			}

			for _, ec := range(ecList.GetControllers()) {
				id := ec.GetEdgeControllerId()
				selectors[id] = &grpc_inventory_manager_go.AssetSelector{
					OrganizationId: orgId,
					EdgeControllerId: id,
				}
			}
		}
	} else {
		var assetList *grpc_inventory_go.AssetList
		var err error

		ctx, cancel := contexts.InventoryContext()
		defer cancel()

		if ecId != "" {
			// We start with all assets for an Edge Controller
			assetList, err = m.assetsClient.ListControllerAssets(ctx, &grpc_inventory_go.EdgeControllerId{
				OrganizationId: orgId,
				EdgeControllerId: ecId,
			})
			if err != nil {
				return nil, derrors.NewUnavailableError("unable to retrieve assets for edge controller", err).WithParams(ecId)
			}
		} else {
			// If there's no Edge Controller, we start with all assets for
			// an organization
			assetList, err = m.assetsClient.List(ctx, &grpc_organization_go.OrganizationId{
				OrganizationId: orgId,
			})
			if err != nil {
				return nil, derrors.NewUnavailableError("unable to retrieve assets for organization", err).WithParams(orgId)
			}
		}

		for _, asset := range(assetList.GetAssets()) {
			if selectedAsset(asset, selector) {
				addAsset(selectors, asset)
			}
		}
	}

	// Lastly, we filter out disabled and unavailable Edge Controllers
	ctx, cancel := contexts.InventoryContext()
	defer cancel()

	ecList, err := m.controllersClient.List(ctx, &grpc_organization_go.OrganizationId{
		OrganizationId: orgId,
	})
	if err != nil {
		return nil, derrors.NewUnavailableError("unable to retrieve edge controllers", err).WithParams(orgId)
	}
	for _, ec := range(ecList.GetControllers()) {
		ecId := ec.GetEdgeControllerId()
		lastAlive := ec.GetLastAliveTimestamp()

		if !ec.GetShow() {
			log.Debug().Str("edge-controller-id", ecId).Msg("removing disabled edge controller from selectors")
			delete(selectors, ecId)
		} else if time.Now().UTC().Unix() - lastAlive > edgeControllerAliveTimeout {
			log.Debug().Str("edge-controller-id", ecId).Int64("last-alive", lastAlive).Msg("removing unavailable edge controller from selectors")
			delete(selectors, ec.GetEdgeControllerId())
		}
	}

	return selectors, nil
}

func selectedAsset(asset *grpc_inventory_go.Asset, selector *grpc_inventory_manager_go.AssetSelector) bool {
	// Check org
	orgId := selector.GetOrganizationId()
	if asset.GetOrganizationId() != orgId {
		return false
	}

	// Check Edge Controller
	ecId := selector.GetEdgeControllerId()
	if ecId != "" && asset.GetEdgeControllerId() != ecId {
		return false
	}

	// Check labels
	labels := selector.GetLabels()
	if labels != nil {
		assetLabels := asset.GetLabels()
		if assetLabels == nil {
			return false
		}
		for k, v := range(labels) {
			if assetLabels[k] != v {
				return false
			}
		}
	}

	// All checks succeeded
	return true
}

func addAsset(selectors map[string]*grpc_inventory_manager_go.AssetSelector, asset *grpc_inventory_go.Asset) {
	ecId := asset.GetEdgeControllerId()
	assetId := asset.GetAssetId()
	selector, found := selectors[ecId]
	if found {
		selector.AssetIds = append(selector.AssetIds, assetId)
	} else {
		selectors[ecId] = &grpc_inventory_manager_go.AssetSelector{
			OrganizationId: asset.GetOrganizationId(),
			EdgeControllerId: ecId,
			AssetIds: []string{assetId},
		}
	}
}
