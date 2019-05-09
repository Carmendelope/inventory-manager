/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package commands

import (
	"github.com/nalej/inventory-manager/internal/pkg/config"
	"github.com/nalej/inventory-manager/internal/pkg/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

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
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().IntVar(&cfg.Port, "port", 5510, "Port to receive management communications")
	runCmd.PersistentFlags().StringVar(&cfg.VPNManagerAddress, "vpnManagerAddress", "localhost:5666", "VPN Server Manager address (host:port)")
	runCmd.PersistentFlags().StringVar(&cfg.AuthxAddress, "authxAddress", "localhost:8810",
		"Authx address (host:port)")
	runCmd.PersistentFlags().StringVar(&cfg.ManagementClusterURL, "managementClusterURL", "nalej.cluster.local",
		"Management URL (base DNS)")

}