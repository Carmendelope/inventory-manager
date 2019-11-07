/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

import (
	"github.com/nalej/derrors"
	"github.com/nalej/device-api/version"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"time"
)

type Config struct {
	// Debug level is active.
	Debug bool
	// Port where the gRPC API service will listen requests.
	Port int
	// DsnURL with the host to configure dns
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
	// CACertPath with the path of the CA.
	CACertPath string
	// CACertRaw certificate
	CACertRaw string
	// ControllerThreshold maximum time (seconds) between ping to decide if a controller is offline or online
	ControllerThreshold time.Duration
	// AssetThreshold maximum time (seconds) between ping to decide if an asset is offline or online
	AssetThreshold time.Duration
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
	if conf.QueueAddress == "" {
		return derrors.NewInvalidArgumentError("queueAddress must not be empty")
	}
	if conf.NetworkManagerAddress == "" {
		return derrors.NewInvalidArgumentError("networkManagerAddress cannot be empty")
	}
	if conf.EdgeInventoryProxyAddress == "" {
		return derrors.NewInvalidArgumentError("edgeInventoryProxy cannot be empty")
	}
	if conf.CACertPath == "" {
		return derrors.NewInvalidArgumentError("caCertPath cannot be empty")
	}

	err := conf.loadCACert()
	if err != nil {
		return err
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
	log.Info().Str("Cert Path", conf.CACertPath).Msg("CA files")
	log.Info().Str("EdgeController", conf.ControllerThreshold.String()).Str("Asset", conf.AssetThreshold.String()).Msg("Online/Offline Threshold")
}

// LoadCert loads the CA certificate in memory.
func (conf *Config) loadCACert() derrors.Error {
	content, err := ioutil.ReadFile(conf.CACertPath)
	if err != nil {
		return derrors.AsError(err, "cannot load management CA certificate")
	}
	conf.CACertRaw = string(content)
	return nil
}
