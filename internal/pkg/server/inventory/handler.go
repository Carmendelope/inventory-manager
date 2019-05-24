/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package inventory

import (
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/inventory-manager/internal/pkg/entities"
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

func (h*Handler) List(_ context.Context, orgID *grpc_organization_go.OrganizationId) (*grpc_inventory_manager_go.InventoryList, error) {
	verr := entities.ValidOrganizationID(orgID)
	if verr != nil {
		return nil, conversions.ToGRPCError(verr)
	}
	list, err := h.manager.List(orgID)
	if err != nil{
		return nil, err
	}
	return list, nil
}

func (h*Handler) Summary(_ context.Context, orgID *grpc_organization_go.OrganizationId) (*grpc_inventory_manager_go.InventorySummary, error) {
	verr := entities.ValidOrganizationID(orgID)
	if verr != nil {
		return nil, conversions.ToGRPCError(verr)
	}
	panic("implement me")
}
