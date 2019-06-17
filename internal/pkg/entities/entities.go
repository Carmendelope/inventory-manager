/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
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



var AgentOpResponseFromGRPC = map[grpc_inventory_manager_go.AgentOpStatus]grpc_inventory_go.AgentOpStatus{
	grpc_inventory_manager_go.AgentOpStatus_SCHEDULED:grpc_inventory_go.AgentOpStatus_SCHEDULED,
	grpc_inventory_manager_go.AgentOpStatus_SUCCESS:grpc_inventory_go.AgentOpStatus_SUCCESS,
	grpc_inventory_manager_go.AgentOpStatus_FAIL:grpc_inventory_go.AgentOpStatus_FAIL,

}

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

func ValidEdgeControllerId(edgeControllerID *grpc_inventory_go.EdgeControllerId) derrors.Error {
	if edgeControllerID.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("organization_id must not be empty")
	}
	if edgeControllerID.EdgeControllerId == "" {
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

func ValidUpdateGeolocationRequest(request *grpc_inventory_manager_go.UpdateGeolocationRequest) derrors.Error {
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
		OrganizationId: device.OrganizationId,
		DeviceId: device.DeviceId,
		DeviceGroupId: device.DeviceGroupId,
		AssetDeviceId: fmt.Sprintf("%s#%s", device.DeviceGroupId, device.DeviceId),
		RegisterSince: device.RegisterSince,
		Labels: device.Labels,
		Location: device.Location,
		DeviceApiKey: device.DeviceApiKey,
		DeviceStatus: device.DeviceStatus,
		Enabled: device.Enabled,
	}
}

func ValidAssetSelector(selector *grpc_inventory_manager_go.AssetSelector) derrors.Error {
	if selector == nil {
		return derrors.NewInvalidArgumentError("empty asset selector")
	}
	if selector.GetOrganizationId() == "" {
		return derrors.NewInvalidArgumentError("organization_id cannot be empty")
	}
	return nil
}

func ValidTimeRange(timeRange *grpc_inventory_manager_go.QueryMetricsRequest_TimeRange) derrors.Error {
	if !(timeRange.GetTimestamp() == 0) {
		if timeRange.GetTimeStart() != 0 || timeRange.GetTimeEnd() != 0 || timeRange.GetResolution() != 0 {
			return derrors.NewInvalidArgumentError("timestamp is set; start, end and resolution should be 0").
				WithParams(timeRange.GetTimestamp(), timeRange.GetTimeStart(),
				timeRange.GetTimeEnd(), timeRange.GetResolution())
		}
	} else {
		if timeRange.GetTimeStart() == 0 && timeRange.GetTimeEnd() == 0 {
			return derrors.NewInvalidArgumentError("timestamp is not set; either start, end or both should be set").
				WithParams(timeRange.GetTimestamp(), timeRange.GetTimeStart(),
				timeRange.GetTimeEnd(), timeRange.GetResolution())
		}
	}

	return nil
}

func ValidQueryMetricsRequest(request *grpc_inventory_manager_go.QueryMetricsRequest) derrors.Error {
	// We check the asset selector so we know we have an organization ID.
	derr := ValidAssetSelector(request.GetAssets())
	if derr != nil {
		return derr
	}

	// Check the time range to either be a point in time or a range
	derr = ValidTimeRange(request.GetTimeRange())
	if derr != nil {
		return derr
	}

	// See [NP-1520]
	if len(request.GetAssets().GetAssetIds()) != 1 && request.GetAggregation() == grpc_inventory_manager_go.QueryMetricsRequest_NONE {
		return derrors.NewInvalidArgumentError("metrics for more than one asset requested without aggregation method")
	}

	return nil
}

func ValidAssetUninstalledId (request  *grpc_inventory_go.AssetUninstalledId) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("organization_id cannot be empty")
	}
	if request.AssetId == "" {
		return derrors.NewInvalidArgumentError("asset_id cannot be empty")
	}
	return nil
}
