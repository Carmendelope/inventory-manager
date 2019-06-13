/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package device

import (
	"fmt"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-edge-inventory-proxy-go"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-network-go"
	"github.com/nalej/grpc-vpn-server-go"
	"github.com/nalej/inventory-manager/internal/pkg/config"
)

type Manager struct {
	controllersClient 	grpc_inventory_go.ControllersClient
	assetsClient  		grpc_inventory_go.AssetsClient
	proxyClient 		grpc_edge_inventory_proxy_go.EdgeControllerProxyClient
	authxClient       	grpc_authx_go.InventoryClient
	certClient        	grpc_authx_go.CertificatesClient
	vpnClient         	grpc_vpn_server_go.VPNServerClient
	netMngrClient     	grpc_network_go.ServiceDNSClient
	// EdgeControllerAPIURL with the URL of the EIC API to accept join request.
	edgeControllerAPIURL string
	dnsUrl               string
	config               config.Config
}

func NewManager(
	authxClient grpc_authx_go.InventoryClient, certClient grpc_authx_go.CertificatesClient, controllerClient grpc_inventory_go.ControllersClient,
	vpnClient grpc_vpn_server_go.VPNServerClient, netManagerClient grpc_network_go.ServiceDNSClient, assetClient grpc_inventory_go.AssetsClient,
	proxyClient grpc_edge_inventory_proxy_go.EdgeControllerProxyClient, cfg config.Config) Manager {
	return Manager{
		authxClient:          authxClient,
		certClient:           certClient,
		controllersClient:    controllerClient,
		vpnClient:            vpnClient,
		netMngrClient:        netManagerClient,
		assetsClient:         assetClient,
		proxyClient:          proxyClient,
		edgeControllerAPIURL: fmt.Sprintf("eic-api.%s", cfg.ManagementClusterURL),
		dnsUrl:               cfg.DnsURL,
		config:               cfg,
	}
}