/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package entities

import (
	"fmt"
	"github.com/nalej/derrors"
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

func ValidUpdateAssetRequest (request *grpc_inventory_go.UpdateAssetRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("organization_id cannot be empty")
	}
	if request.AssetId == "" {
		return derrors.NewInvalidArgumentError("asset_id cannot be empty")
	}
	return nil
}