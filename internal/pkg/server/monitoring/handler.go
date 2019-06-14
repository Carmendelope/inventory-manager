/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package monitoring

import (
	"context"

	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/inventory-manager/internal/pkg/entities"
)

type Handler struct {
	manager *Manager
}

func NewHandler(manager *Manager) *Handler {
	return &Handler{
		manager: manager,
	}
}

func (h *Handler) ListMetrics(ctx context.Context, selector *grpc_inventory_manager_go.AssetSelector) (*grpc_inventory_manager_go.MetricsList, error) {
	derr := entities.ValidAssetSelector(selector)
	if derr != nil {
		return nil, conversions.ToGRPCError(derr)
	}

	return h.manager.ListMetrics(selector)
}

func (h *Handler) QueryMetrics(ctx context.Context, request *grpc_inventory_manager_go.QueryMetricsRequest) (*grpc_inventory_manager_go.QueryMetricsResult, error) {
	derr := entities.ValidQueryMetricsRequest(request)
	if derr != nil {
		return nil, conversions.ToGRPCError(derr)
	}

	return h.manager.QueryMetrics(request)
}
