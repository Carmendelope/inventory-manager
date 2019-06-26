/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package monitoring

import (
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
		ctx, cancel := contexts.ProxyContext() // Manual calling cancel to avoid big list of defers
		list, err := m.proxyClient.ListMetrics(ctx, proxyRequest)
		cancel()
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

	// TODO [NP-1520]
	// For now, we only allow metrics from a single edge controller
	// - Create a request for each Edge Controller and execute
	// - Unify the results

	if len(selectors) > 1 {
		return nil, derrors.NewUnimplementedError("Querying metrics from more than one edge controller not supported")
	}

	results := make([]*grpc_inventory_manager_go.QueryMetricsResult, 0, len(selectors))

	for _, selector := range(selectors) {
		proxyRequest := &grpc_inventory_manager_go.QueryMetricsRequest{
			Assets: selector,
			Metrics: request.GetMetrics(),
			TimeRange: request.GetTimeRange(),
			Aggregation: request.GetAggregation(),
		}
		ecId := selector.GetEdgeControllerId()
		log.Debug().Interface("request", proxyRequest).Msg("proxy request for QueryMetrics")
		ctx, cancel := contexts.ProxyContext() // Manual calling cancel to avoid big list of defers
		result, err := m.proxyClient.QueryMetrics(ctx, proxyRequest)
		cancel()
		if err != nil {
			// We still want to query to working edge controllers
			log.Warn().Str("edge-controller-id", ecId).Err(err).Msg("failed calling QueryMetrics")
			continue
		}
		results = append(results, result)
	}

	// We only have one - if we have anything
	if len(results) == 0 {
		return &grpc_inventory_manager_go.QueryMetricsResult{}, nil
	}
	return results[0], nil
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

	// If we have a set of assets, we'll go from there
	if len(assetIds) > 0 {
		// If we have explicit assets, that's the minimum set we start from
		for _, id := range(assetIds) {
			ctx, cancel := contexts.InventoryContext()
			// Calling cancel manually to avoid stacking up a lot of defers
			asset, err := m.assetsClient.Get(ctx, &grpc_inventory_go.AssetId{
				OrganizationId: orgId,
				AssetId: id,
			})
			cancel()
			if err != nil {
				return nil, derrors.NewUnavailableError("unable to retrieve asset information", err).WithParams(id)
			}
			if selectedAsset(asset, selector) {
				addAsset(selectors, asset)
			}
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
		// If we have more filters to apply (labels, groups), we need to get
		// a set of matching assets to filter. The Edge Controller doesn't
		// have this info so we need to do it here and provide an exhaustive
		// list of assets to query.
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
