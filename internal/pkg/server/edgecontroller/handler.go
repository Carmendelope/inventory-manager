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

type Handler struct{
	manager Manager
}

func NewHandler(manager Manager) *Handler{
	return &Handler{
		manager:manager,
	}
}

func (h *Handler) CreateEICToken(_ context.Context, orgID *grpc_organization_go.OrganizationId) (*grpc_inventory_manager_go.EICJoinToken, error) {
	verr := entities.ValidOrganizationID(orgID)
	if verr != nil {
		return nil, conversions.ToGRPCError(verr)
	}
	token, err := h.manager.CreateEICToken(orgID)
	if err != nil{
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
	log.Debug().Interface("info", info).Msg("EIC start")
	// TODO Implement start logic
	return &grpc_common_go.Success{}, nil
}

func (h *Handler) UnlinkEIC(_ context.Context, edgeControllerID *grpc_inventory_go.EdgeControllerId) (*grpc_common_go.Success, error) {
	verr := entities.ValidEdgeControllerId(edgeControllerID)
	if verr != nil {
		return nil, conversions.ToGRPCError(verr)
	}
	err := h.manager.UnlinkEIC(edgeControllerID)
	if err != nil{
		return nil, err
	}
	return &grpc_common_go.Success{}, nil
}

func (h *Handler) ConfigureEIC(context.Context, *grpc_inventory_manager_go.ConfigureEICRequest) (*grpc_common_go.Success, error) {
	return nil, derrors.NewUnimplementedError("ConfigureEIC is not implemented")
}
