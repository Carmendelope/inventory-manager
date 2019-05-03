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

	return nil
}

func (conf *Config) Print() {
	log.Info().Str("app", version.AppVersion).Str("commit", version.Commit).Msg("Version")
	log.Info().Int("port", conf.Port).Msg("gRPC port")
	log.Info().Str("URL", conf.VPNManagerAddress).Msg("VPN Manager")
	log.Info().Str("URL", conf.AuthxAddress).Msg("Authx")
}
