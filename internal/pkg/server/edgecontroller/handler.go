/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package edgecontroller

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-common-go"
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

func (h *Handler) CreateEICToken(_ context.Context, orgID *grpc_organization_go.OrganizationId) (*grpc_inventory_manager_go.EICJoinToken, error) {
	verr := entities.ValidOrganizationID(orgID)
	if verr != nil {
		return nil, conversions.ToGRPCError(verr)
	}
	token, err := h.manager.CreateEICToken(orgID)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (h *Handler) EICJoin(_ context.Context, request *grpc_inventory_manager_go.EICJoinRequest) (*grpc_inventory_manager_go.EICJoinResponse, error) {
	verr := entities.ValidEICJoinRequest(request)
	if verr != nil {
		return nil, conversions.ToGRPCError(verr)
	}
	return h.manager.EICJoin(request)
}

func (h *Handler) EICStart(_ context.Context, info *grpc_inventory_manager_go.EICStartInfo) (*grpc_common_go.Success, error) {
	verr := entities.ValidEICStartInfo(info)
	if verr != nil {
		return nil, conversions.ToGRPCError(verr)
	}
	err := h.manager.EICStart(info)
	if err != nil {
		return nil, err
	}
	log.Info().Str("organization_id", info.OrganizationId).Str("edge_controller_id", info.EdgeControllerId).Str("IP", info.Ip).Msg("EIC start has been processed")
	return &grpc_common_go.Success{}, nil
}

func (h *Handler) UnlinkEIC(_ context.Context, edgeControllerID *grpc_inventory_go.EdgeControllerId) (*grpc_common_go.Success, error) {
	vErr := entities.ValidEdgeControllerId(edgeControllerID)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	return h.manager.UnlinkEIC(edgeControllerID)

}

func (h *Handler) ConfigureEIC(context.Context, *grpc_inventory_manager_go.ConfigureEICRequest) (*grpc_common_go.Success, error) {
	return nil, derrors.NewUnimplementedError("ConfigureEIC is not implemented")
}

func (h *Handler) EICAlive(_ context.Context, edgeControllerID *grpc_inventory_go.EdgeControllerId) (*grpc_common_go.Success, error) {
	vErr := entities.ValidEdgeControllerId(edgeControllerID)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}

	err := h.manager.EICAlive(edgeControllerID)
	if err != nil {
		return nil, err
	}

	return &grpc_common_go.Success{}, nil
}

func (h *Handler) CallbackECOperation(_ context.Context, response *grpc_inventory_manager_go.EdgeControllerOpResponse) (*grpc_common_go.Success, error) {
	vErr := entities.ValidEdgeControllerOpResponse(response)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}

	return h.manager.CallbackECOperation(response)
}

// UpdateECLocation operation to update the geolocation
func (h *Handler) UpdateECGeolocation(_ context.Context, in *grpc_inventory_manager_go.UpdateGeolocationRequest) (*grpc_inventory_go.EdgeController, error){
	log.Info().Msg("UpdateECGeolocation")
	vErr := entities.ValidUpdateGeolocationRequest(in)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}

	return h.manager.UpdateECGeolocation(in)

}

// UpdateEC updates an Edge Controller
func (h *Handler) UpdateEC(ctx context.Context, request *grpc_inventory_go.UpdateEdgeControllerRequest) (*grpc_inventory_go.EdgeController, error){
	log.Info().Msg("UpdateEC")
	vErr := entities.ValidUpdateECRequest(request)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}

	return h.manager.UpdateEC(request)

}