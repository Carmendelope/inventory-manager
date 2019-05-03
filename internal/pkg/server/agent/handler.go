/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package agent

import (
	"github.com/nalej/grpc-common-go"
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

func (h*Handler) InstallAgent(context.Context, *grpc_inventory_manager_go.InstallAgentRequest) (*grpc_inventory_manager_go.InstallAgentResponse, error) {
	panic("implement me")
}

func (h*Handler) CreateAgentJoinToken(context.Context, *grpc_organization_go.OrganizationId) (*grpc_inventory_manager_go.AgentJoinToken, error) {
	panic("implement me")
}

func (h*Handler) AgentJoin(context.Context, *grpc_inventory_manager_go.AgentJoinRequest) (*grpc_inventory_manager_go.AgentJoinResponse, error) {
	panic("implement me")
}

func (h*Handler) LogAgentAlive(context.Context, *grpc_inventory_manager_go.AgentIds) (*grpc_common_go.Success, error) {
	panic("implement me")
}

func (h*Handler) TriggerAgentOperation(context.Context, *grpc_inventory_manager_go.AgentOpRequest) (*grpc_inventory_manager_go.AgentOpResponse, error) {
	panic("implement me")
}

func (h*Handler) CallbackAgentOperation(context.Context, *grpc_inventory_manager_go.AgentOpResponse) (*grpc_common_go.Success, error) {
	panic("implement me")
}

func (h*Handler) ListAgentOperations(context.Context, *grpc_inventory_manager_go.AgentId) (*grpc_inventory_manager_go.AgentOpResponseList, error) {
	panic("implement me")
}

func (h*Handler) DeleteAgentOperation(context.Context, *grpc_inventory_manager_go.AgentOperationId) (*grpc_common_go.Success, error) {
	panic("implement me")
}


