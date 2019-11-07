/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package agent

import (
	"context"
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-edge-inventory-proxy-go"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/inventory-manager/internal/pkg/server/contexts"
	"github.com/rs/zerolog/log"
	"github.com/satori/go.uuid"
	"time"
)

const ProxyTimeout = time.Second * 60

type Manager struct {
	proxyClient grpc_edge_inventory_proxy_go.EdgeControllerProxyClient
	assetClient grpc_inventory_go.AssetsClient
	controllersClient 	grpc_inventory_go.ControllersClient
	CACert      string
}

func NewManager(proxyClient grpc_edge_inventory_proxy_go.EdgeControllerProxyClient, assetClient grpc_inventory_go.AssetsClient,
	controllersClient grpc_inventory_go.ControllersClient,	caCert string) Manager {
	return Manager{
		proxyClient: proxyClient,
		assetClient: assetClient,
		controllersClient: controllersClient,
		CACert:      caCert,
	}
}

func (m *Manager) generateToken() string {
	return uuid.NewV4().String()
}

func (m *Manager) InstallAgent(request *grpc_inventory_manager_go.InstallAgentRequest) (*grpc_inventory_manager_go.EdgeControllerOpResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ProxyTimeout)
	defer cancel()
	// propagate the certificate to the agent.
	request.CaCert = m.CACert
	response, err := m.proxyClient.InstallAgent(ctx, request)
	if err != nil{
		return nil, err
	}

	// update the last operation result in EC
	err = m.updateLastECOpResponse(response)
	if err != nil {
		log.Warn().Str("operation_id", response.OperationId).Str("status", response.Status.String()).Str("info", response.Info).
			Str("error", conversions.ToDerror(err).DebugReport()).Msg("error updating install agent response")
	}

	return response, nil
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

	location := &grpc_inventory_go.InventoryLocation{
		Geolocation: request.Geolocation,
	}

	// send a message to system model to add the agent
	asset, err := m.assetClient.Add(ctx, &grpc_inventory_go.AddAssetRequest{
		OrganizationId:   request.OrganizationId,
		EdgeControllerId: request.EdgeControllerId,
		AgentId:          request.AgentId,
		Labels:           request.Labels,
		Os:               request.Os,
		Hardware:         request.Hardware,
		Storage:          request.Storage,
		Location:      	  location,
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
		Status: response.Status,
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
		log.Error().Err(err).Interface("response", response).Msg("cannot store last op summary")
		return nil, err
	}

	return &grpc_common_go.Success{}, nil
}

func (m *Manager) ListAgentOperations(agentID *grpc_inventory_manager_go.AgentId) (*grpc_inventory_manager_go.AgentOpResponseList, error) {
	panic("implement me")
}

func (m *Manager) DeleteAgentOperation(operationID *grpc_inventory_manager_go.AgentOperationId) (*grpc_common_go.Success, error) {
	panic("implement me")
}

func (m *Manager) updateLastECOpResponse(request *grpc_inventory_manager_go.EdgeControllerOpResponse) error {

	// updates the EC with last operation result
	ctxSMUpdate, cancelSMUpdate := contexts.SMContext()
	defer cancelSMUpdate()

	_, err := m.controllersClient.Update(ctxSMUpdate, &grpc_inventory_go.UpdateEdgeControllerRequest{
		OrganizationId: request.OrganizationId,
		EdgeControllerId: request.EdgeControllerId,
		AddLabels: false,
		RemoveLabels: false,
		UpdateLastAlive: false,
		UpdateGeolocation: false,
		UpdateLastOpSummary: true,
		LastOpSummary: &grpc_inventory_go.ECOpSummary{
			OperationId: request.OperationId,
			Timestamp: request.Timestamp,
			Status: request.Status,
			Info:request.Info,
		},
	})
	return err
}

func (m *Manager) UninstallAgent( request *grpc_inventory_manager_go.UninstallAgentRequest) (*grpc_inventory_manager_go.EdgeControllerOpResponse, error){

	// Check if the asset_id is correct, (exist and is managed by edge_controller_id)
	ctxSM, cancelSM := contexts.SMContext()
	defer cancelSM()

	asset, err := m.assetClient.Get(ctxSM, &grpc_inventory_go.AssetId{
		OrganizationId: request.OrganizationId,
		AssetId: request.AssetId,
	} )
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), ProxyTimeout)
	defer cancel()

	res, err :=  m.proxyClient.UninstallAgent(ctx, &grpc_inventory_manager_go.FullUninstallAgentRequest{
		OrganizationId: request.OrganizationId,
		EdgeControllerId: asset.EdgeControllerId,
		AssetId: request.AssetId,
		Force: request.Force,
	})
	if err != nil {
		return nil, err
	}

	// update the last operation result in system-model
	err = m.updateLastECOpResponse(res)
	if err != nil {
		log.Warn().Str("operation_id", res.OperationId).Str("status", res.Status.String()).Str("info", res.Info).
			Str("error", conversions.ToDerror(err).DebugReport()).Msg("error updating uninstall agent response")
	}

	return res, nil

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

	// update last_operation_result

	err = m.updateLastECOpResponse(&grpc_inventory_manager_go.EdgeControllerOpResponse{
		OrganizationId: assetID.OrganizationId,
		EdgeControllerId: assetID.EdgeControllerId,
		OperationId: assetID.OperationId,
		Timestamp: time.Now().Unix(), // no timestamp received
		Status: grpc_inventory_go.OpStatus_SUCCESS,
		Info: fmt.Sprintf("agent %s uninstalled", assetID.AssetId),
	})
	if err != nil {
		log.Warn().Str("edge_controller_id", assetID.EdgeControllerId).Str("operation_id", assetID.OperationId).
			Msg("unable to update last operation result")
	}

	return &grpc_common_go.Success{}, nil
}