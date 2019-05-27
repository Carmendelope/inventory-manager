/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package agent

import (
	"context"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/grpc-edge-inventory-proxy-go"
	"time"
)

const ProxyTimeout = time.Second * 60

type Manager struct{
	proxyClient grpc_edge_inventory_proxy_go.EdgeControllerProxyClient
}

func NewManager(proxyClient grpc_edge_inventory_proxy_go.EdgeControllerProxyClient) Manager{
	return Manager{
		proxyClient: proxyClient,
	}
}

func (m * Manager) InstallAgent(request *grpc_inventory_manager_go.InstallAgentRequest) (*grpc_inventory_manager_go.InstallAgentResponse, error) {
	panic("implement me")
}

func (m * Manager) CreateAgentJoinToken(edgeControllerID *grpc_inventory_go.EdgeControllerId) (*grpc_inventory_manager_go.AgentJoinToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ProxyTimeout)
	defer cancel()
	return m.proxyClient.CreateAgentJoinToken(ctx, edgeControllerID)
}

func (m * Manager) AgentJoin(request *grpc_inventory_manager_go.AgentJoinRequest) (*grpc_inventory_manager_go.AgentJoinResponse, error) {
	panic("implement me")
}

func (m * Manager) LogAgentAlive(agentIds *grpc_inventory_manager_go.AgentIds) (*grpc_common_go.Success, error) {
	panic("implement me")
}

func (m * Manager) TriggerAgentOperation(request *grpc_inventory_manager_go.AgentOpRequest) (*grpc_inventory_manager_go.AgentOpResponse, error) {
	panic("implement me")
}

func (m * Manager) CallbackAgentOperation(response *grpc_inventory_manager_go.AgentOpResponse) (*grpc_common_go.Success, error) {
	panic("implement me")
}

func (m * Manager) ListAgentOperations(agentID *grpc_inventory_manager_go.AgentId) (*grpc_inventory_manager_go.AgentOpResponseList, error) {
	panic("implement me")
}

func (m * Manager) DeleteAgentOperation(operationID *grpc_inventory_manager_go.AgentOperationId) (*grpc_common_go.Success, error) {
	panic("implement me")
}