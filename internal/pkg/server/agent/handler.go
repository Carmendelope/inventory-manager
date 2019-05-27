/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package agent

import (
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-inventory-manager-go"
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

func (h*Handler) InstallAgent(context.Context, *grpc_inventory_manager_go.InstallAgentRequest) (*grpc_inventory_manager_go.InstallAgentResponse, error) {
	panic("implement me")
}

func (h*Handler) CreateAgentJoinToken(_ context.Context, edgeControllerID *grpc_inventory_go.EdgeControllerId) (*grpc_inventory_manager_go.AgentJoinToken, error) {
	verr := entities.ValidEdgeControllerId(edgeControllerID)
	if verr != nil {
		return nil, conversions.ToGRPCError(verr)
	}
	log.Debug().Interface("edgeControllerID", edgeControllerID).Msg("triggering agent join token creation")
	return h.manager.CreateAgentJoinToken(edgeControllerID)
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


