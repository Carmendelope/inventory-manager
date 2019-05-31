/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package commands

import (
	"github.com/nalej/inventory-manager/internal/pkg/config"
	"github.com/nalej/inventory-manager/internal/pkg/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"time"
)

const DefaultControllerStatusThreshold = "10m"
const DefaultAssetStatusThreshold = "10m"

var cfg = config.Config{}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Launch the Inventory Manager",
	Long:  `Launch the Inventory Manager`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		log.Info().Msg("Launching API!")
		cfg.Debug = debugLevel
		server := server.NewService(cfg)
		server.Run()
	},
}

func init() {

	controllerThreshold, _ := time.ParseDuration(DefaultControllerStatusThreshold)
	assetThreshold, _ := time.ParseDuration(DefaultAssetStatusThreshold)

	rootCmd.AddCommand(runCmd)
	runCmd.Flags().IntVar(&cfg.Port, "port", 5510, "Port to receive management communications")
	runCmd.PersistentFlags().StringVar(&cfg.VPNManagerAddress, "vpnManagerAddress", "localhost:5666", "VPN Server Manager address (host:port)")
	runCmd.PersistentFlags().StringVar(&cfg.AuthxAddress, "authxAddress", "localhost:8810",
		"Authx address (host:port)")
	runCmd.PersistentFlags().StringVar(&cfg.SystemModelAddress, "systemModelAddress", "localhost:8800",
		"System Model address (host:port)")
	runCmd.PersistentFlags().StringVar(&cfg.DeviceManagerAddress, "deviceManagerAddress", "localhost:6010",
		"Device Manager address (host:port)")
	runCmd.PersistentFlags().StringVar(&cfg.ManagementClusterURL, "managementClusterURL", "nalej.cluster.local",
		"Management URL (base DNS)")
	runCmd.PersistentFlags().StringVar(&cfg.QueueAddress, "queueAddress", "localhost:6650",
		"Queue address (base DNS)")
	runCmd.PersistentFlags().StringVar(&cfg.NetworkManagerAddress, "networkManagerAddress", "localhost:8000",
		"Network Manager address (host:port)")
	runCmd.PersistentFlags().StringVar(&cfg.EdgeInventoryProxyAddress, "edgeInventoryProxyAddress", "localhost:5544",
		"Edge Inventory Proxy address (host:port)")
	runCmd.PersistentFlags().StringVar(&cfg.DnsURL, "dnsURL", "",
		"Management URL (base DNS)")
	runCmd.Flags().StringVar(&cfg.CACertPath, "caCertPath", "", "CA certificate path")
	runCmd.Flags().DurationVar(&cfg.ControllerThreshold, "controllerThreshold", controllerThreshold, "Threshold between ping to decide if a controller is offline/online")
	runCmd.Flags().DurationVar(&cfg.AssetThreshold, "assetThreshold", assetThreshold, "Threshold between ping to decide if an asset is offline/online")

}
