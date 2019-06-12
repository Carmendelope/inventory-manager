/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package monitoring

import (
	"github.com/nalej/derrors"

	"github.com/nalej/grpc-edge-inventory-proxy-go"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/grpc-organization-go"

	"github.com/nalej/inventory-manager/internal/pkg/server/contexts"
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

func (m *Manager) ListMetrics(selector *grpc_inventory_manager_go.AssetSelector) (*grpc_inventory_manager_go.MetricsList, error) {
	// Get a selector for each relevant Edge Controller

	// Create a request for each Edge Controller and execute

	// Unify the results

	return nil, derrors.NewUnimplementedError("ListMetrics is not implemented")
}

func (m *Manager) QueryMetrics(request *grpc_inventory_manager_go.QueryMetricsRequest) (*grpc_inventory_manager_go.QueryMetricsResult, error) {
	// Get a selector for each relevant Edge Controller

	// Create a request for each Edge Controller and execute

	// Unify the results

	return nil, derrors.NewUnimplementedError("QueryMetrics is not implemented")
}

// GetSelectors will turn a single asset selector into a selector per
// edge controller, so we can sent the request to each edge controller
// that needs it. This is done by creating a list of assets from either
// AssetIds in the source selector or, if none are provided, by retrieving
// all Assets for an OrganizationId or EdgeControllerId. We then filter
// that list to remove the Assets that don't match the groups and labels and
// sort them out by EdgeControllerId.
func (m *Manager) getSelectors(selector *grpc_inventory_manager_go.AssetSelector) (map[string]*grpc_inventory_manager_go.AssetSelector, derrors.Error) {
	selectors := make(map[string]*grpc_inventory_manager_go.AssetSelector)

	orgId := selector.GetOrganizationId()
	ecId := selector.GetEdgeControllerId()
	assetIds := selector.GetAssetIds()

	// If we have explicit assets, that's the minimum set we start from
	if len(assetIds) > 0 {
		for _, id := range(assetIds) {
			ctx, _ := contexts.InventoryContext()
			asset, err := m.assetsClient.Get(ctx, &grpc_inventory_go.AssetId{
				OrganizationId: orgId,
				AssetId: id,
			})
			if err != nil {
				return nil, derrors.NewUnavailableError("unable to retrieve asset information", err).WithParams(id)
			}
			if selectedAsset(asset, selector) {
				addAsset(selectors, asset)
			}
		}
		return selectors, nil
	}

	// If we have no assets, no labels, no groups, we make selectors just
	// for edge controllers
	if len(selector.GetLabels()) == 0 && len(selector.GetGroupIds()) == 0 {
		// With Edge Controller set, we actually just need the original
		if ecId != "" {
			selectors[ecId] = selector
			return selectors, nil
		}

		// If not, we need a selector per edge controller
		ctx, _ := contexts.InventoryContext()
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

		return selectors, nil
	}


	// If we have labels or groups, we need to make the assets explicit to
	// be able to filter
	ctx, _ := contexts.InventoryContext()
	var assetList *grpc_inventory_go.AssetList
	var err error

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

	return selectors, nil
}

func selectedAsset(asset *grpc_inventory_go.Asset, selector *grpc_inventory_manager_go.AssetSelector) bool {
	// Don't count hidden assets
	if !asset.GetShow() {
		return false
	}

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
