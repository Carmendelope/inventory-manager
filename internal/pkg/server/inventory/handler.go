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
	"github.com/rs/zerolog/log"
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
	verr := entities.ValidOrganizationID(orgID)
	if verr != nil {
		return nil, conversions.ToGRPCError(verr)
	}
	list, err := h.manager.List(orgID)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (h *Handler) GetControllerExtendedInfo(_ context.Context, edgeControllerID *grpc_inventory_go.EdgeControllerId) (*grpc_inventory_manager_go.EdgeControllerExtendedInfo, error) {
	verr := entities.ValidEdgeControllerId(edgeControllerID)
	if verr != nil {
		return nil, conversions.ToGRPCError(verr)
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
	verr := entities.ValidAssetID(assetID)
	if verr != nil {
		return nil, conversions.ToGRPCError(verr)
	}
	asset, err := h.manager.GetAssetInfo(assetID)
	if err != nil {
		return nil, err
	}
	return asset, nil
}

func (h *Handler) Summary(_ context.Context, orgID *grpc_organization_go.OrganizationId) (*grpc_inventory_manager_go.InventorySummary, error) {
	verr := entities.ValidOrganizationID(orgID)
	if verr != nil {
		return nil, conversions.ToGRPCError(verr)
	}
	panic("implement me")
}

// UpdateAsset updates an asset in the inventory.
func (h *Handler) UpdateAsset(ctx context.Context, in *grpc_inventory_go.UpdateAssetRequest) (*grpc_inventory_go.Asset, error){
	log.Info().Msg("Update asset")
	vErr := entities.ValidUpdateAssetRequest(in)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}

	return h.manager.UpdateAssetLocation(in)
}

// UpdateDevice updates a device in the inventory.
func (h *Handler) UpdateDevice(ctx context.Context, in *grpc_inventory_manager_go.UpdateDeviceLocationRequest) (*grpc_inventory_manager_go.Device, error){
	log.Info().Msg("Update device")
	vErr := entities.ValidUpdateDeviceLocationRequest(in)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}

	return h.manager.UpdateDeviceLocation(in)
}

func (h *Handler) GetDeviceInfo (ctx context.Context, request *grpc_inventory_manager_go.DeviceId) (*grpc_inventory_manager_go.Device, error) {
	return nil, nil
}