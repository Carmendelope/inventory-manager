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

package entities

import (
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-device-manager-go"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/grpc-organization-go"
)


func ValidEICJoinToken(request *grpc_inventory_manager_go.EICJoinToken) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("organization_id must not be empty")
	}
	if request.Token == "" {
		return derrors.NewInvalidArgumentError("token must not be empty")
	}
	return nil
}

func ValidOrganizationID(orgID *grpc_organization_go.OrganizationId) derrors.Error {
	if orgID.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("organization_id must not be empty")
	}
	return nil
}

func ValidEICJoinRequest(request *grpc_inventory_manager_go.EICJoinRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("organization_id must not be empty")
	}
	if request.Name == "" {
		return derrors.NewInvalidArgumentError("name must not be empty")
	}
	return nil
}

func GetEdgeControllerName(organizationID string, edgeControllerID string) string {
	return fmt.Sprintf("%s-%s.eic", organizationID, edgeControllerID)
}

func ValidEdgeControllerOpResponse(response * grpc_inventory_manager_go.EdgeControllerOpResponse) derrors.Error{
	if response.OrganizationId == ""{
		return derrors.NewInvalidArgumentError("organization_id cannot be empty")
	}
	if response.EdgeControllerId == ""{
		return derrors.NewInvalidArgumentError("edge_controller_id cannot be empty")
	}
	if response.OperationId == ""{
		return derrors.NewInvalidArgumentError("operation_id cannot be empty")
	}
	return nil
}

func ValidEdgeControllerId(edgeControllerID *grpc_inventory_go.EdgeControllerId) derrors.Error {
	if edgeControllerID.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("organization_id must not be empty")
	}
	if edgeControllerID.EdgeControllerId == "" {
		return derrors.NewInvalidArgumentError("edge_controller_id must not be empty")
	}
	return nil
}

func ValidUnlinkECRequest(request *grpc_inventory_manager_go.UnlinkECRequest) derrors.Error{
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("organization_id must not be empty")
	}
	if request.EdgeControllerId == "" {
		return derrors.NewInvalidArgumentError("edge_controller_id must not be empty")
	}
	return nil
}

func ValidEICStartInfo(info *grpc_inventory_manager_go.EICStartInfo) derrors.Error {
	if info.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("organization_id must not be empty")
	}
	if info.EdgeControllerId == "" {
		return derrors.NewInvalidArgumentError("edge_controller_id must not be empty")
	}
	// TODO Validate IP regex
	if info.Ip == "" {
		return derrors.NewInvalidArgumentError("ip must not be empty")
	}
	return nil
}

func ValidInstallAgentRequest(request *grpc_inventory_manager_go.InstallAgentRequest) derrors.Error{
	if request.OrganizationId == ""{
		return derrors.NewInvalidArgumentError("organization_id cannot be empty")
	}
	if request.EdgeControllerId == ""{
		return derrors.NewInvalidArgumentError("edge_controller_id cannot be empty")
	}
	if request.TargetHost == ""{
		return derrors.NewInvalidArgumentError("target_host cannot be empty")
	}
	return nil
}

func ValidAgentJoinRequest(request *grpc_inventory_manager_go.AgentJoinRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("organization_id cannot be empty")
	}
	if request.EdgeControllerId == "" {
		return derrors.NewInvalidArgumentError("edge_controller_id cannot be empty")
	}
	if request.AgentId == "" {
		return derrors.NewInvalidArgumentError("agent_id cannot be empty")
	}
	return nil
}

func ValidAgentsAlive(request *grpc_inventory_manager_go.AgentsAlive) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("organization_id cannot be empty")
	}
	if request.EdgeControllerId == "" {
		return derrors.NewInvalidArgumentError("edge_controller_id cannot be empty")
	}
	if request.Agents == nil || len(request.Agents) <= 0 {
		return derrors.NewInvalidArgumentError("agents cannot be empty")
	}
	return nil
}

func ValidAssetID(id *grpc_inventory_go.AssetId) derrors.Error {
	if id.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("organization_id cannot be empty")
	}
	if id.AssetId == "" {
		return derrors.NewInvalidArgumentError("asset_id cannot be empty")
	}
	return nil
}


func ValidUninstallAgentRequest(request *grpc_inventory_manager_go.UninstallAgentRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("organization_id cannot be empty")
	}
	if request.AssetId == "" {
		return derrors.NewInvalidArgumentError("asset_id cannot be empty")
	}
	return nil
}

func ValidUpdateGeolocationRequest(request *grpc_inventory_manager_go.UpdateGeolocationRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("organization_id cannot be empty")
	}
	if request.EdgeControllerId == "" {
		return derrors.NewInvalidArgumentError("edge_controller_id cannot be empty")
	}
	return nil
}

func ValidUpdateECRequest (request *grpc_inventory_go.UpdateEdgeControllerRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("organization_id cannot be empty")
	}
	if request.EdgeControllerId == "" {
		return derrors.NewInvalidArgumentError("edge_controller_id cannot be empty")
	}
	return nil
}

func ValidAgentOpRequest(request *grpc_inventory_manager_go.AgentOpRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("organization_id cannot be empty")
	}
	if request.EdgeControllerId == "" {
		return derrors.NewInvalidArgumentError("edge_controller_id cannot be empty")
	}
	if request.AssetId == "" {
		return derrors.NewInvalidArgumentError("asset_id cannot be empty")
	}
	if request.OperationId == "" {
		return derrors.NewInvalidArgumentError("operation_id cannot be empty")
	}
	if request.Plugin == "" {
		return derrors.NewInvalidArgumentError("plugin cannot be empty")
	}
	return nil
}

func ValidAgentOpResponse (request *grpc_inventory_manager_go.AgentOpResponse) derrors.Error {

	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("organization_id cannot be empty")
	}
	if request.EdgeControllerId == "" {
		return derrors.NewInvalidArgumentError("edge_controller_id cannot be empty")
	}
	if request.AssetId == "" {
		return derrors.NewInvalidArgumentError("asset_id cannot be empty")
	}
	if request.OperationId == "" {
		return derrors.NewInvalidArgumentError("operation_id cannot be empty")
	}
	return nil
}

func NewDeviceFromGRPC (device *grpc_device_manager_go.Device) *grpc_inventory_manager_go.Device {
	return &grpc_inventory_manager_go.Device{
		OrganizationId:       device.OrganizationId,
		DeviceGroupId:        device.DeviceGroupId,
		DeviceId:             device.DeviceId,
		AssetDeviceId:        fmt.Sprintf("%s#%s", device.DeviceGroupId, device.DeviceId),
		RegisterSince:        device.RegisterSince,
		Labels:               device.Labels,
		Enabled:              device.Enabled,
		DeviceApiKey:         device.DeviceApiKey,
		DeviceStatus:         device.DeviceStatus,
		Location:             device.Location,
		AssetInfo:            device.AssetInfo,
	}
}

func ValidAssetUninstalledId (request  *grpc_inventory_go.AssetUninstalledId) derrors.Error {
	return nil
}

func ValidUpdateAssetRequest (request *grpc_inventory_go.UpdateAssetRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("organization_id cannot be empty")
	}
	if request.AssetId == "" {
		return derrors.NewInvalidArgumentError("asset_id cannot be empty")
	}
	return nil
}

func ValidDeviceId (request *grpc_inventory_manager_go.DeviceId) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("organization_id cannot be empty")
	}
	if request.AssetDeviceId == "" {
		return derrors.NewInvalidArgumentError("asset_device_id cannot be empty")
	}
	return nil
}

func ValidUpdateDeviceLocationRequest (request *grpc_inventory_manager_go.UpdateDeviceLocationRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("organization_id cannot be empty")
	}
	if request.AssetDeviceId == "" {
		return derrors.NewInvalidArgumentError("asset_device_id cannot be empty")
	}
	if request.Location != nil && request.Location.Geolocation == "" {
		return derrors.NewInvalidArgumentError("location cannot be empty")
	}
	return nil
}
