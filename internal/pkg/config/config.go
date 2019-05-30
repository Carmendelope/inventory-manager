/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package config

import (
	"github.com/nalej/derrors"
	"github.com/nalej/device-api/version"
	"github.com/rs/zerolog/log"
)

type Config struct {
	// Debug level is active.
	Debug bool
	// Port where the gRPC API service will listen requests.
	Port int
	//DsnURL with the host to configure dns
	DnsURL string
	// VPNManagerAddress with the host:port to connect to the VPN manager.
	VPNManagerAddress string
	// AuthxAddress with the host:port to connect to the Authx manager.
	AuthxAddress string
	// ManagementClusterURL with the host where the management cluster resides
	ManagementClusterURL string
	// SystemModelAddress with the host:port to connect to the System Model manager
	SystemModelAddress string
	// DeviceManagerAddress with the host:port to connect to the Device Manager.
	DeviceManagerAddress string
	// QueueAddress with the bus address
	QueueAddress string
	// NetworkManagerAddress with the address of the network manager
	NetworkManagerAddress string
	// EdgeInventoryProxyAddress with the address of the edge inventory proxy.
	EdgeInventoryProxyAddress string

}


func (conf *Config) Validate() derrors.Error {

	if conf.Port <= 0 {
		return derrors.NewInvalidArgumentError("port must be valid")
	}
	if conf.DnsURL == "" {
		return derrors.NewInvalidArgumentError("DnsURL must be set")
	}

	if conf.VPNManagerAddress == "" {
		return derrors.NewInvalidArgumentError("vpnManager must be set")
	}

	if conf.AuthxAddress == "" {
		return derrors.NewInvalidArgumentError("authxAddress must be set")
	}

	if conf.ManagementClusterURL == "" {
		return derrors.NewInvalidArgumentError("managementClusterURL must be set")
	}

	if conf.SystemModelAddress == "" {
		return derrors.NewInvalidArgumentError("SystemModelAddress must be set")
	}

	if conf.DeviceManagerAddress == "" {
		return derrors.NewInvalidArgumentError("deviceManagerAddress must be set")
	}
	if conf.QueueAddress == ""{
		return derrors.NewInvalidArgumentError("queueAddress must not be empty")
	}
	if conf.NetworkManagerAddress == ""{
		return derrors.NewInvalidArgumentError("networkManagerAddress cannot be empty")
	}
	if conf.EdgeInventoryProxyAddress == ""{
		return derrors.NewInvalidArgumentError("edgeInventoryProxy cannot be empty")
	}

	return nil
}

func (conf *Config) Print() {
	log.Info().Str("app", version.AppVersion).Str("commit", version.Commit).Msg("Version")
	log.Info().Int("port", conf.Port).Msg("gRPC port")
	log.Info().Str("URL", conf.DnsURL).Msg("DNS URL")
	log.Info().Str("URL", conf.VPNManagerAddress).Msg("VPN Manager component")
	log.Info().Str("URL", conf.AuthxAddress).Msg("Authx component")
	log.Info().Str("URL", conf.SystemModelAddress).Msg("System Model component")
	log.Info().Str("URL", conf.DeviceManagerAddress).Msg("Device Manager component")
	log.Info().Str("URL", conf.ManagementClusterURL).Msg("Management cluster")
	log.Info().Str("URL", conf.QueueAddress).Msg("Queue")
	log.Info().Str("URL", conf.NetworkManagerAddress).Msg("Network Manager")
	log.Info().Str("URL", conf.EdgeInventoryProxyAddress).Msg("Edge Inventory Proxy")
}
