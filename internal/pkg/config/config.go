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
	// VPNManagerAddress with the host:port to connect to the VPN manager.
	VPNManagerAddress string
	// AuthxAddress with the host:port to connect to the Authx manager.
	AuthxAddress string
	// ManagementClusterURL with the host where the management cluster resides
	ManagementClusterURL string
	//VPNServerURL with the URL of the VPN server accepting connections.
	VPNServerURL string
	// SystemModelAddress with the host:port to connect to the System Model manager
	SystemModelAddress string
	// DeviceManagerAddress with the host:port to connect to the Device Manager.
	DeviceManagerAddress string
	//DsnURL with the host to configure dns
	DnsURL string
}


func (conf *Config) Validate() derrors.Error {

	if conf.Port <= 0 {
		return derrors.NewInvalidArgumentError("port must be valid")
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

	if conf.DnsURL == "" {
		return derrors.NewInvalidArgumentError("DnsURL must be set")
	}

	return nil
}

func (conf *Config) Print() {
	log.Info().Str("app", version.AppVersion).Str("commit", version.Commit).Msg("Version")
	log.Info().Int("port", conf.Port).Msg("gRPC port")
	log.Info().Str("URL", conf.VPNManagerAddress).Msg("VPN Manager component")
	log.Info().Str("URL", conf.AuthxAddress).Msg("Authx component")
	log.Info().Str("URL", conf.SystemModelAddress).Msg("System Model component")
	log.Info().Str("URL", conf.DeviceManagerAddress).Msg("Device Manager component")
	log.Info().Str("URL", conf.ManagementClusterURL).Msg("Management cluster")
	log.Info().Str("URL", conf.VPNServerURL).Msg("VPN Server URL")
	log.Info().Str("URL", conf.DnsURL).Msg("DNS URL")
}
