// SPDX-FileCopyrightText: 2025 OpenCHAMI Contributors
//
// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"
	"os"

	"bootstrap/internal/diag"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ochami_bootstrap",
	Short: "Bootstrap inventory generation and NIC discovery via Redfish",
	PersistentPreRun: func(cmd *cobra.Command, args []string) { //nolint:revive
		// propagate debug flag to internal diagnostics
		diag.Debug = debugFlag
	},
}

var debugFlag bool

// Execute is the entry point for the CLI.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debugFlag, "debug", false, "enable verbose debug logging")
}
