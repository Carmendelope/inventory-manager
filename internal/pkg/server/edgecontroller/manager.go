/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package edgecontroller

import (
	"context"
	"fmt"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/grpc-network-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/grpc-vpn-server-go"
	"github.com/nalej/inventory-manager/internal/pkg/config"
	"github.com/nalej/inventory-manager/internal/pkg/entities"
	"github.com/nalej/inventory-manager/internal/pkg/server/contexts"
	"github.com/rs/zerolog/log"
	"time"
)

// TODO: we have only one proxy, change this when more proxies are added
const proxy = "proxy0-vpn.service.nalej"

type Manager struct{
	controllersClient grpc_inventory_go.ControllersClient
	authxClient grpc_authx_go.InventoryClient
	vpnClient grpc_vpn_server_go.VPNServerClient
	netMngrClient grpc_network_go.ServiceDNSClient
	// EdgeControllerAPIURL with the URL of the EIC API to accept join request.
	edgeControllerAPIURL string
    dnsUrl string
	config config.Config
}

func NewManager(authxClient grpc_authx_go.InventoryClient, controllerClient grpc_inventory_go.ControllersClient,
	vpnClient grpc_vpn_server_go.VPNServerClient, netManagerClient grpc_network_go.ServiceDNSClient, cfg config.Config) Manager{
	return Manager{
		authxClient:authxClient,
		controllersClient:controllerClient,
		vpnClient:vpnClient,
		netMngrClient: netManagerClient,
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
		Hostname:             fmt.Sprintf("vpn-server.%s:5555", m.config.ManagementClusterURL),
		Proxyname:            fmt.Sprintf("%s:5544", proxy),
	}
	return &grpc_inventory_manager_go.EICJoinResponse{
		OrganizationId:       added.OrganizationId,
		EdgeControllerId:     added.EdgeControllerId,
		Credentials:          credentials,
	}, nil
}

func (m * Manager) EICStart(info *grpc_inventory_manager_go.EICStartInfo) error {
	log.Debug().Interface("info", info).Msg("EICStart")
	// Update DNS entry
	netCtx, netCancel := context.WithTimeout(context.Background(), time.Second * 10)
	defer netCancel()
	addRequest := &grpc_network_go.AddServiceDNSEntryRequest{
		OrganizationId:       info.OrganizationId,
		Fqdn:                 fmt.Sprintf("%s-vpn", info.EdgeControllerId),
		Ip:                   info.Ip,
		Tags:                 []string{"EIC"},
	}
	log.Debug().Interface("request", addRequest).Msg("registering entry")
	_, err := m.netMngrClient.AddEntry(netCtx, addRequest)

	if err != nil{
		dErr := conversions.ToDerror(err)
		log.Error().Str("trace", dErr.DebugReport()).Msg("cannot register edge controller IP on the DNS")
		return err
	}
	return nil
}

func (m * Manager) UnlinkEIC(edgeControllerID *grpc_inventory_go.EdgeControllerId) error {
	// Remove credentials on the VPN
	vpnCtx, vpnCancel := contexts.VPNManagerContext()
	defer vpnCancel()
	eicUsername := entities.GetEdgeControllerName(edgeControllerID.OrganizationId, edgeControllerID.EdgeControllerId)
	deleteUserRequest := &grpc_vpn_server_go.DeleteVPNUserRequest{
		OrganizationId:       edgeControllerID.OrganizationId,
		Username:             eicUsername,
	}
	_, err := m.vpnClient.DeleteVPNUser(vpnCtx, deleteUserRequest)
	if err != nil{
		log.Warn().Interface("edgeControllerId", edgeControllerID).Msg("failed to delete EC from VPN")
		return err
	}
	// Remove entry from SM

	smCtx, smCancel := contexts.SMContext()
	defer smCancel()
	_, err = m.controllersClient.Remove(smCtx, edgeControllerID)
	if err != nil{
		log.Warn().Interface("edgeControllerId", edgeControllerID).Msg("failed to delete EC from SM")
		return err
	}
	// Send a message to the EIC with the unlink throught the EIC proxy.

	log.Warn().Msg("unlink EIC not fully implemented, EIC not informed")
	return nil
}

func (m * Manager) EICAlive( eic *grpc_inventory_go.EdgeControllerId)  error {

	ctx, cancel := contexts.SMContext()
	defer cancel()
	// set timestamp
	// send a message to system-model to update the timestamp
	_, err := m.controllersClient.Update(ctx, &grpc_inventory_go.UpdateEdgeControllerRequest{
		OrganizationId: eic.OrganizationId,
		EdgeControllerId: eic.EdgeControllerId,
		AddLabels:false,
		UpdateLastAlive: true,
		LastAliveTimestamp: time.Now().Unix(),
	})
	if err != nil {
		return err
	}

	return nil

}

