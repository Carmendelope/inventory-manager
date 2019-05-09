/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package entities

import (
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/grpc-organization-go"
)

func ValidEICJoinToken(request *grpc_inventory_manager_go.EICJoinToken) derrors.Error{
	if request.OrganizationId == ""{
		return derrors.NewInvalidArgumentError("organization_id must not be empty")
	}
	if request.Token == ""{
		return derrors.NewInvalidArgumentError("token must not be empty")
	}
	return nil
}

func ValidOrganizationID(orgID *grpc_organization_go.OrganizationId) derrors.Error{
	if orgID.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("organization_id must not be empty")
	}
	return nil
}

func ValidEICJoinRequest(request *grpc_inventory_manager_go.EICJoinRequest) derrors.Error{
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("organization_id must not be empty")
	}
	if request.Name == "" {
		return derrors.NewInvalidArgumentError("name must not be empty")
	}
	return nil
}

func GetEdgeControllerName(organizationID string, edgeControllerID string) string {
	return fmt.Sprintf("%s-%s@eic", organizationID, edgeControllerID)
}