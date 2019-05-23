/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package edgecontroller

import (
	"fmt"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-vpn-server-go"
	"github.com/nalej/inventory-manager/internal/pkg/config"
	"github.com/nalej/inventory-manager/internal/pkg/entities"
	"github.com/nalej/inventory-manager/internal/pkg/server/contexts"
	"github.com/rs/zerolog/log"
)

// TODO: we have only one proxy, change this when more proxies are added
const proxy = "proxy0.vpn.service.nalej"

type Manager struct{
	controllersClient grpc_inventory_go.ControllersClient
	authxClient grpc_authx_go.InventoryClient
	vpnClient grpc_vpn_server_go.VPNServerClient
	// EdgeControllerAPIURL with the URL of the EIC API to accept join request.
	edgeControllerAPIURL string
    dnsUrl string
	config config.Config
}

func NewManager(authxClient grpc_authx_go.InventoryClient, controllerClient grpc_inventory_go.ControllersClient,
	vpnClient grpc_vpn_server_go.VPNServerClient, cfg config.Config) Manager{
	return Manager{
		authxClient:authxClient,
		controllersClient:controllerClient,
		vpnClient:vpnClient,
		edgeControllerAPIURL: fmt.Sprintf("eic-api.%s", cfg.ManagementClusterURL),
		dnsUrl: cfg.DnsURL,
		config: cfg,
	}
}

func (m * Manager) CreateEICToken(orgID *grpc_organization_go.OrganizationId) (*grpc_inventory_manager_go.EICJoinToken, error) {
	ctx, cancel := contexts.AuthxContext()
	defer cancel()
	token, err := m.authxClient.CreateEICJoinToken(ctx, orgID)
	if err != nil{
		return nil, err
	}
	return &grpc_inventory_manager_go.EICJoinToken{
		OrganizationId:       token.OrganizationId,
		Token:                token.Token,
		Cacert:               token.Cacert,
		JoinUrl:              m.edgeControllerAPIURL,
		DnsUrl:               m.dnsUrl,
	}, nil
}

func (m * Manager) EICJoin(request *grpc_inventory_manager_go.EICJoinRequest) (*grpc_inventory_manager_go.EICJoinResponse, error) {

	log.Debug().Interface("request", request).Msg("EICJoin")

	// Add the EIC to system model
	ctx, cancel := contexts.InventoryContext()
	defer cancel()
	toAdd := &grpc_inventory_go.AddEdgeControllerRequest{
		OrganizationId:       request.OrganizationId,
		Name:                 request.Name,
		Labels:               request.Labels,
	}
	added, err := m.controllersClient.Add(ctx, toAdd)
	if err != nil{
		return nil, err
	}

	eicUsername := entities.GetEdgeControllerName(request.OrganizationId, added.EdgeControllerId)

	// Create a set of credentials
	vpnCtx, vpnCancel := contexts.VPNManagerContext()
	defer vpnCancel()
	eicUser := &grpc_vpn_server_go.AddVPNUserRequest{
		Username:            eicUsername,
		OrganizationId:      request.OrganizationId,
	}

	vpnCredentials, err := m.vpnClient.AddVPNUser(vpnCtx, eicUser)
	if err != nil{
		return nil, err
	}

	credentials := &grpc_inventory_manager_go.VPNCredentials{
		// TODO Is this needed?
		Cacert:               "",
		Username:             vpnCredentials.Username,
		Password:             vpnCredentials.Password,
		Hostname:             proxy,
	}
	return &grpc_inventory_manager_go.EICJoinResponse{
		OrganizationId:       added.OrganizationId,
		EdgeControllerId:     added.EdgeControllerId,
		Credentials:          credentials,
	}, nil
}

func (m * Manager) EICStart(info *grpc_inventory_manager_go.EICStartInfo) (*grpc_common_go.Success, error) {
	log.Debug().Interface("info", info).Msg("EICStart")
	//TODO implement logic
	return &grpc_common_go.Success{}, nil
}

func (m * Manager) UnlinkEIC(edgeControllerID *grpc_inventory_go.EdgeControllerId) error {
	panic("implement me")
}
