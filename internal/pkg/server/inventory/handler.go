/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package inventory

import (
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

func (h*Handler) List(context.Context, *grpc_organization_go.OrganizationId) (*grpc_inventory_manager_go.AssetList, error) {
	panic("implement me")
}

func (h*Handler) Summary(context.Context, *grpc_organization_go.OrganizationId) (*grpc_inventory_manager_go.InventorySummary, error) {
	panic("implement me")
}
