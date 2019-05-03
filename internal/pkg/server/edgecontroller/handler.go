/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package edgecontroller

import (
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/grpc-organization-go"
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

func (h *Handler) CreateEICToken(context.Context, *grpc_organization_go.OrganizationId) (*grpc_inventory_manager_go.EICJoinToken, error) {
	panic("implement me")
}

func (h *Handler) EICJoin(context.Context, *grpc_authx_go.EICJoinRequest) (*grpc_inventory_manager_go.EICJoinResponse, error) {
	panic("implement me")
}

func (h *Handler) EICStart(context.Context, *grpc_inventory_manager_go.EICStartInfo) (*grpc_common_go.Success, error) {
	panic("implement me")
}

func (h *Handler) UnlinkEIC(context.Context, *grpc_inventory_go.EdgeControllerId) (*grpc_common_go.Success, error) {
	panic("implement me")
}
