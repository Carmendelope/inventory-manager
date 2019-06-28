/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package agent

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/inventory-manager/internal/pkg/entities"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

type Handler struct {
	manager Manager
}

func NewHandler(manager Manager) *Handler {
	return &Handler{
		manager: manager,
	}
}

func (h *Handler) InstallAgent(_ context.Context, request *grpc_inventory_manager_go.InstallAgentRequest) (*grpc_inventory_manager_go.InstallAgentResponse, error) {
	vErr := entities.ValidInstallAgentRequest(request)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	log.Debug().Str("organizationID", request.OrganizationId).Str("edgeControllerId", request.EdgeControllerId).Str("target", request.TargetHost).Msg("Installing an agent")
	return h.manager.InstallAgent(request)
}

func (h *Handler) CreateAgentJoinToken(_ context.Context, edgeControllerID *grpc_inventory_go.EdgeControllerId) (*grpc_inventory_manager_go.AgentJoinToken, error) {
	verr := entities.ValidEdgeControllerId(edgeControllerID)
	if verr != nil {
		return nil, conversions.ToGRPCError(verr)
	}
	log.Debug().Interface("edgeControllerID", edgeControllerID).Msg("triggering agent join token creation")
	return h.manager.CreateAgentJoinToken(edgeControllerID)
}

func (h *Handler) AgentJoin(_ context.Context, request *grpc_inventory_manager_go.AgentJoinRequest) (*grpc_inventory_manager_go.AgentJoinResponse, error) {
	vErr := entities.ValidAgentJoinRequest(request)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}

	log.Debug().Str("organization_id", request.OrganizationId).Str("edge_controller_id", request.EdgeControllerId).
		Str("agent_id", request.AgentId).Msg("Agent join")

	return h.manager.AgentJoin(request)
}

func (h *Handler) LogAgentAlive(_ context.Context, request *grpc_inventory_manager_go.AgentsAlive) (*grpc_common_go.Success, error) {
	vErr := entities.ValidAgentsAlive(request)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}

	log.Debug().Str("organization_id", request.OrganizationId).Str("edge_controller_id", request.EdgeControllerId).
		Int("agents len", len(request.Agents)).Msg("Agents alive")

	err := h.manager.LogAgentAlive(request)
	if err != nil {
		return nil, err
	}
	return &grpc_common_go.Success{}, nil

}

func (h *Handler) TriggerAgentOperation(_ context.Context, opRequest *grpc_inventory_manager_go.AgentOpRequest) (*grpc_inventory_manager_go.AgentOpResponse, error) {
	vErr := entities.ValidAgentOpRequest(opRequest)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	log.Debug().Interface("operation", opRequest).Msg("Trigger Agent operation")

	return h.manager.TriggerAgentOperation(opRequest)
}

func (h *Handler) CallbackAgentOperation(_ context.Context, response *grpc_inventory_manager_go.AgentOpResponse) (*grpc_common_go.Success, error) {
	vErr := entities.ValidAgentOpResponse(response)

	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}

	return h.manager.CallbackAgentOperation(response)
}

func (h *Handler) ListAgentOperations(context.Context, *grpc_inventory_manager_go.AgentId) (*grpc_inventory_manager_go.AgentOpResponseList, error) {
	return nil, conversions.ToGRPCError(derrors.NewUnimplementedError("not implemented yet"))
}

func (h *Handler) DeleteAgentOperation(context.Context, *grpc_inventory_manager_go.AgentOperationId) (*grpc_common_go.Success, error) {
	return nil, conversions.ToGRPCError(derrors.NewUnimplementedError("not implemented yet"))
}

// UninstallAgent operation to uninstall an agent
func (h *Handler) UninstallAgent(_ context.Context, request *grpc_inventory_manager_go.UninstallAgentRequest) (*grpc_common_go.Success, error){

	vErr := entities.ValidUninstallAgentRequest(request)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}

	return h.manager.UninstallAgent(request)

}
func (h *Handler) UninstalledAgent(_ context.Context,  assetID *grpc_inventory_go.AssetUninstalledId) (*grpc_common_go.Success, error){
	vErr := entities.ValidAssetUninstalledId(assetID)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}

	return h.manager.UninstalledAgent(assetID)
}

