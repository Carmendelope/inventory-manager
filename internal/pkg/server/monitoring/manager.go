/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package monitoring

import (
	"time"

	"github.com/nalej/derrors"

	"github.com/nalej/grpc-edge-inventory-proxy-go"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-inventory-manager-go"
)

const ProxyTimeout = time.Second * 60

type Manager struct {
	proxyClient grpc_edge_inventory_proxy_go.EdgeControllerProxyClient
	assetClient grpc_inventory_go.AssetsClient
}

func NewManager(proxyClient grpc_edge_inventory_proxy_go.EdgeControllerProxyClient, assetClient grpc_inventory_go.AssetsClient) *Manager {
	return &Manager{
		proxyClient: proxyClient,
		assetClient: assetClient,
	}
}

func (m *Manager) ListMetrics(selector *grpc_inventory_manager_go.AssetSelector) (*grpc_inventory_manager_go.MetricsList, error) {
	return nil, derrors.NewUnimplementedError("ListMetrics is not implemented")
}

func (m *Manager) QueryMetrics(request *grpc_inventory_manager_go.QueryMetricsRequest) (*grpc_inventory_manager_go.QueryMetricsResult, error) {
	return nil, derrors.NewUnimplementedError("QueryMetrics is not implemented")
}
