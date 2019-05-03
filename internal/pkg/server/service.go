/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package server

import (
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/inventory-manager/internal/pkg/config"
	"github.com/nalej/inventory-manager/internal/pkg/server/agent"
	"github.com/nalej/grpc-vpn-server-go"
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/inventory-manager/internal/pkg/server/edgecontroller"
	"github.com/nalej/inventory-manager/internal/pkg/server/inventory"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type Service struct {
	Configuration config.Config
}

// NewService creates a new service.
func NewService(conf config.Config) *Service {
	return &Service{
		conf,
	}
}

type Clients struct {
	cpnClient grpc_vpn_server_go.VPNServerClient
	authxClient grpc_authx_go.InventoryClient
}

// GetClients creates the required connections with the remote clients.
func (s *Service) GetClients() (*Clients, derrors.Error) {
	vpnConn, err := grpc.Dial(s.Configuration.VPNManagerAddress, grpc.WithInsecure())
	if err != nil {
		return nil, derrors.AsError(err, "cannot create connection with VPN manager")
	}
	aConn, err := grpc.Dial(s.Configuration.AuthxAddress, grpc.WithInsecure())
	if err != nil {
		return nil, derrors.AsError(err, "cannot create connection with authx")
	}
	imClient := grpc_vpn_server_go.NewVPNServerClient(vpnConn)
	aClient := grpc_authx_go.NewInventoryClient(aConn)
	return &Clients{imClient, aClient}, nil
}

// Run the service, launch the REST service handler.
func (s *Service) Run() error {
	cErr := s.Configuration.Validate()
	if cErr != nil {
		log.Fatal().Str("err", cErr.DebugReport()).Msg("invalid configuration")
	}
	s.Configuration.Print()

	_, cErr = s.GetClients()
	if cErr != nil {
		log.Fatal().Str("err", cErr.DebugReport()).Msg("Cannot create clients")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Configuration.Port))
	if err != nil {
		log.Fatal().Errs("failed to listen: %v", []error{err})
	}

	// Create handlers

	agentManager := agent.NewManager()
	agentHandler := agent.NewHandler(agentManager)

	ecManager := edgecontroller.NewManager()
	ecHandler := edgecontroller.NewHandler(ecManager)

	invManager := inventory.NewManager()
	invHandler := inventory.NewHandler(invManager)


	grpcServer := grpc.NewServer()
	grpc_inventory_manager_go.RegisterInventoryServer(grpcServer, invHandler)
	grpc_inventory_manager_go.RegisterAgentServer(grpcServer, agentHandler)
	grpc_inventory_manager_go.RegisterEICServer(grpcServer, ecHandler)


	if s.Configuration.Debug{
		log.Info().Msg("Enabling gRPC server reflection")
		// Register reflection service on gRPC server.
		reflection.Register(grpcServer)
	}
	log.Info().Int("port", s.Configuration.Port).Msg("Launching gRPC server")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Errs("failed to serve: %v", []error{err})
	}
	return nil
}