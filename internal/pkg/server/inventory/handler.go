/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package inventory

import (
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/inventory-manager/internal/pkg/entities"
	"golang.org/x/net/context"
)

type Handler struct {
	manager Manager
}

func NewHandler(manager Manager) *Handler {
	return &Handler{
		manager: manager,
	}
}

func (h *Handler) List(_ context.Context, orgID *grpc_organization_go.OrganizationId) (*grpc_inventory_manager_go.InventoryList, error) {
	vErr := entities.ValidOrganizationID(orgID)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	list, err := h.manager.List(orgID)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (h *Handler) GetControllerExtendedInfo(_ context.Context, edgeControllerID *grpc_inventory_go.EdgeControllerId) (*grpc_inventory_manager_go.EdgeControllerExtendedInfo, error) {
	vErr := entities.ValidEdgeControllerId(edgeControllerID)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	controller, assets, err := h.manager.GetControllerExtendedInfo(edgeControllerID)
	if err != nil {
		return nil, err
	}
	return &grpc_inventory_manager_go.EdgeControllerExtendedInfo{
		Controller:    controller,
		ManagedAssets: assets,
	}, nil
}

func (h *Handler) GetAssetInfo(_ context.Context, assetID *grpc_inventory_go.AssetId) (*grpc_inventory_manager_go.Asset, error) {
	vErr := entities.ValidAssetID(assetID)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	asset, err := h.manager.GetAssetInfo(assetID)
	if err != nil {
		return nil, err
	}
	return asset, nil
}

// GetAssetInfo returns the information of a given device
func (h *Handler)GetDeviceInfo(ctx context.Context, deviceID *grpc_inventory_manager_go.DeviceId) (*grpc_inventory_manager_go.Device, error) {
	vErr := entities.ValidDeviceId(deviceID)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	return h.manager.GetDeviceInfo(deviceID)

}

func (h *Handler) Summary(_ context.Context, orgID *grpc_organization_go.OrganizationId) (*grpc_inventory_manager_go.InventorySummary, error) {
	vErr := entities.ValidOrganizationID(orgID)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	return h.manager.Summary(orgID)
}

// UpdateAsset updates an asset in the inventory.
func (h *Handler) UpdateAsset(ctx context.Context, in *grpc_inventory_go.UpdateAssetRequest) (*grpc_inventory_go.Asset, error){
	vErr := entities.ValidUpdateAssetRequest(in)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}

	return h.manager.UpdateAssetLocation(in)
}

// UpdateDevice updates a device in the inventory.
func (h *Handler) UpdateDevice(ctx context.Context, in *grpc_inventory_manager_go.UpdateDeviceLocationRequest) (*grpc_inventory_manager_go.Device, error){
	vErr := entities.ValidUpdateDeviceLocationRequest(in)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}

	return h.manager.UpdateDeviceLocation(in)
}
