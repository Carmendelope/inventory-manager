/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package agent

import (
	"context"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-edge-inventory-proxy-go"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/inventory-manager/internal/pkg/entities"
	"github.com/nalej/inventory-manager/internal/pkg/server/contexts"
	"github.com/rs/zerolog/log"
	"github.com/satori/go.uuid"
	"time"
)

const ProxyTimeout = time.Second * 60

type Manager struct {
	proxyClient grpc_edge_inventory_proxy_go.EdgeControllerProxyClient
	assetClient grpc_inventory_go.AssetsClient
	CACert      string
}

func NewManager(proxyClient grpc_edge_inventory_proxy_go.EdgeControllerProxyClient, assetClient grpc_inventory_go.AssetsClient, caCert string) Manager {
	return Manager{
		proxyClient: proxyClient,
		assetClient: assetClient,
		CACert:      caCert,
	}
}

func (m *Manager) generateToken() string {
	return uuid.NewV4().String()
}

func (m *Manager) InstallAgent(request *grpc_inventory_manager_go.InstallAgentRequest) (*grpc_inventory_manager_go.InstallAgentResponse, error) {
	panic("implement me")
}

func (m *Manager) CreateAgentJoinToken(edgeControllerID *grpc_inventory_go.EdgeControllerId) (*grpc_inventory_manager_go.AgentJoinToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ProxyTimeout)
	defer cancel()
	token, err := m.proxyClient.CreateAgentJoinToken(ctx, edgeControllerID)
	if err != nil{
		return nil, err
	}
	token.CaCert = m.CACert
	return token, nil
}

func (m *Manager) AgentJoin(request *grpc_inventory_manager_go.AgentJoinRequest) (*grpc_inventory_manager_go.AgentJoinResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ProxyTimeout)
	defer cancel()

	// send a message to system model to add the agent
	asset, err := m.assetClient.Add(ctx, &grpc_inventory_go.AddAssetRequest{
		OrganizationId:   request.OrganizationId,
		EdgeControllerId: request.EdgeControllerId,
		AgentId:          request.AgentId,
		Labels:           request.Labels,
		Os:               request.Os,
		Hardware:         request.Hardware,
		Storage:          request.Storage,
	})
	if err != nil {
		return nil, err
	}

	// generate a token return it
	return &grpc_inventory_manager_go.AgentJoinResponse{
		OrganizationId: asset.OrganizationId,
		AssetId:        asset.AssetId,
		Token:          m.generateToken(),
		CaCert:         m.CACert,
	}, nil
}

// TODO: add a message in system-model to add many alive messages at once
func (m *Manager) LogAgentAlive(agents *grpc_inventory_manager_go.AgentsAlive) error {

	for agent, timestamp := range agents.Agents {
		// send to system model a message to update the timestamp
		ctx, cancel := contexts.SMContext()
		// set timestamp
		// send a message to system-model to update the timestamp
		// check if the IP is changed reviewing AgentsIp map
		updateIp := false
		ip, exists := agents.AgentsIp[agent]
		if !exists {
			ip = ""
		} else {
			updateIp = true
		}
		_, err := m.assetClient.Update(ctx, &grpc_inventory_go.UpdateAssetRequest{
			OrganizationId:      agents.OrganizationId,
			AssetId:             agent,
			AddLabels:           false,
			RemoveLabels:        false,
			UpdateLastOpSummary: false,
			UpdateLastAlive:     true,
			LastAliveTimestamp:  timestamp,
			UpdateIp:            updateIp,
			EicNetIp:            ip,
		})
		if err != nil {
			log.Warn().Str("organizationID", agents.OrganizationId).Str("assetID", agent).Msg("enable to send alive message to sytem-model")
		}

		cancel()
	}

	return nil

}

func (m *Manager) TriggerAgentOperation(request *grpc_inventory_manager_go.AgentOpRequest) (*grpc_inventory_manager_go.AgentOpResponse, error) {


	// Check if the asset_id is correct, (exist and is managed by edge_controller_id)
	ctxSM, cancelSM := contexts.SMContext()
	defer cancelSM()

	asset, err := m.assetClient.Get(ctxSM, &grpc_inventory_go.AssetId{
		OrganizationId: request.OrganizationId,
		AssetId: request.AssetId,
	})
	if err != nil {
		return nil, err
	}

	if asset.EdgeControllerId != request.EdgeControllerId {

		return nil, conversions.ToDerror(derrors.NewInvalidArgumentError("this asset is not managed by the EC").WithParams("edge_controller_id", request.EdgeControllerId).WithParams("asset_id", request.AssetId))
	}

	ctx, cancel := context.WithTimeout(context.Background(), ProxyTimeout)
	defer cancel()

	return m.proxyClient.TriggerAgentOperation(ctx, request)

}

func (m *Manager) CallbackAgentOperation(response *grpc_inventory_manager_go.AgentOpResponse) (*grpc_common_go.Success, error) {
	// Check if the asset_id is correct, (exist and is managed by edge_controller_id)
	ctxSM, cancelSM := contexts.SMContext()
	defer cancelSM()

	opSummary := &grpc_inventory_go.AgentOpSummary{
		OperationId: response.OperationId,
		Timestamp: response.Timestamp,
		Status: entities.AgentOpResponseFromGRPC[response.Status],
		Info: response.Info,
	}

	_, err := m.assetClient.Update(ctxSM, &grpc_inventory_go.UpdateAssetRequest{
		OrganizationId:      response.OrganizationId,
		AssetId:             response.AssetId,
		AddLabels:           false,
		RemoveLabels:        false,
		UpdateLastOpSummary: true,
		LastOpSummary:  opSummary,
		UpdateLastAlive:     false,
		UpdateIp:            false,
	})
	if err != nil {
		log.Warn().Str("organizationID", response.OrganizationId).Str("assetID", response.AssetId).Msg("enable to send alive message to sytem-model to store the op summary")
	}

	return &grpc_common_go.Success{}, nil
}

func (m *Manager) ListAgentOperations(agentID *grpc_inventory_manager_go.AgentId) (*grpc_inventory_manager_go.AgentOpResponseList, error) {
	panic("implement me")
}

func (m *Manager) DeleteAgentOperation(operationID *grpc_inventory_manager_go.AgentOperationId) (*grpc_common_go.Success, error) {
	panic("implement me")
}


func (m *Manager) UninstallAgent( assetID *grpc_inventory_go.AssetId) (*grpc_common_go.Success, error){

	// Check if the asset_id is correct, (exist and is managed by edge_controller_id)
	ctxSM, cancelSM := contexts.SMContext()
	defer cancelSM()

	asset, err := m.assetClient.Get(ctxSM, assetID)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), ProxyTimeout)
	defer cancel()

	return m.proxyClient.UninstallAgent(ctx, &grpc_inventory_manager_go.FullAssetId{
		OrganizationId: assetID.OrganizationId,
		EdgeControllerId: asset.EdgeControllerId,
		AssetId: assetID.AssetId,
	})

}

// UninstalledAgent method to delete an agent when it was uninstalled
func (m *Manager) UninstalledAgent( assetID *grpc_inventory_go.AssetUninstalledId) (*grpc_common_go.Success, error) {

	ctxSM, cancelSM := contexts.SMContext()
	defer cancelSM()

	_, err := m.assetClient.Remove(ctxSM, &grpc_inventory_go.AssetId{
		OrganizationId: assetID.OrganizationId,
		AssetId: assetID.AssetId,
	} )
	if err != nil {
		return nil, err
	}

	log.Debug().Str("asset", assetID.AssetId).Msg("removed from the system")

	return &grpc_common_go.Success{}, nil
}