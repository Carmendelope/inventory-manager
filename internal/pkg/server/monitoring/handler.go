/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package monitoring

import (
	"context"

	"github.com/nalej/derrors"

	"github.com/nalej/grpc-inventory-manager-go"
	_ "github.com/nalej/grpc-utils/pkg/conversions"
	_ "github.com/nalej/inventory-manager/internal/pkg/entities"
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
	return nil, derrors.NewUnimplementedError("ListMetrics is not implemented")
}

func (h *Handler) QueryMetrics(ctx context.Context, request *grpc_inventory_manager_go.QueryMetricsRequest) (*grpc_inventory_manager_go.QueryMetricsResult, error) {
	return nil, derrors.NewUnimplementedError("QueryMetrics is not implemented")
}
