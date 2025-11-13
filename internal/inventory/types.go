// SPDX-FileCopyrightText: 2025 OpenCHAMI Contributors
//
// SPDX-License-Identifier: MIT

// Package inventory defines types for inventory YAML files.
package inventory

// Entry represents a BMC or Node record in the YAML file.
type Entry struct {
	Xname string `yaml:"xname"`
	MAC   string `yaml:"mac"`
	IP    string `yaml:"ip"`
}

// FileFormat is the root YAML structure with bmcs and nodes.
type FileFormat struct {
	BMCs  []Entry `yaml:"bmcs"`
	Nodes []Entry `yaml:"nodes"`
}
