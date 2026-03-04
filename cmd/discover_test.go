// SPDX-FileCopyrightText: 2025 OpenCHAMI Contributors
//
// SPDX-License-Identifier: MIT

package cmd

import (
    "net/http"
    "net/http/httptest"
    "os"
    "testing"
    "time"
    "bootstrap/internal/inventory"
    "gopkg.in/yaml.v3"
)

func TestDiscoverCommand(t *testing.T) {
    // Set up mock Redfish server (TLS) with required endpoints.
    ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        switch r.URL.Path {
        case "/redfish/v1/Systems":
            w.Header().Set("Content-Type", "application/json")
            _, _ = w.Write([]byte(`{"Members":[{"@odata.id":"/redfish/v1/Systems/Self"}]}`))
        case "/redfish/v1/Systems/Self/EthernetInterfaces":
            w.Header().Set("Content-Type", "application/json")
            _, _ = w.Write([]byte(`{"Members":[{"@odata.id":"/redfish/v1/Systems/Self/EthernetInterfaces/1"}]}`))
        case "/redfish/v1/Systems/Self/EthernetInterfaces/1":
            w.Header().Set("Content-Type", "application/json")
            // Return a bootable NIC with a MAC address.
            _, _ = w.Write([]byte(`{"Id":"1","MACAddress":"aa:bb:cc:dd:ee:01"}`))
        case "/redfish/v1/Managers/BMC/EthernetInterfaces":
            w.Header().Set("Content-Type", "application/json")
            _, _ = w.Write([]byte(`{"Members":[{"@odata.id":"/redfish/v1/Managers/BMC/EthernetInterfaces/1"}]}`))
        case "/redfish/v1/Managers/BMC/EthernetInterfaces/1":
            w.Header().Set("Content-Type", "application/json")
            _, _ = w.Write([]byte(`{"Id":"1","MACAddress":"aa:bb:cc:dd:ee:ff","IPv4Addresses":[{"Address":"192.168.100.10","AddressOrigin":"Static"}]}`))
        default:
            // Return empty JSON for any unexpected path.
            w.Header().Set("Content-Type", "application/json")
            _, _ = w.Write([]byte(`{}`))
        }
    }))
    defer ts.Close()

    // Extract host (without scheme) for the BMC entry.
    host := ts.URL[len("https://"):]

    // Prepare a temporary inventory file.
    tmpFile, err := os.CreateTemp("", "inventory-*.yaml")
    if err != nil {
        t.Fatalf("failed to create temp file: %v", err)
    }
    defer os.Remove(tmpFile.Name())

    // Write initial inventory with a single BMC entry.
    inv := inventory.FileFormat{BMCs: []inventory.Entry{{Xname: "x9000c1s0b0", MAC: "", IP: host}}}
    data, _ := yaml.Marshal(&inv)
    if err := os.WriteFile(tmpFile.Name(), data, 0o644); err != nil {
        t.Fatalf("write inventory failed: %v", err)
    }

    // Set required environment variables for Redfish authentication.
    os.Setenv("REDFISH_USER", "admin")
    os.Setenv("REDFISH_PASSWORD", "secret")
    defer os.Unsetenv("REDFISH_USER")
    defer os.Unsetenv("REDFISH_PASSWORD")

    // Configure command‑line flags directly (they are package‑level variables).
    discFile = tmpFile.Name()
    discBMCSubnet = "192.168.100.0/24"
    discNodeSubnet = "192.168.100.0/24"
    discNodeStartIP = ""
    discInsecure = true
    discTimeout = 5 * time.Second
    discSSHPubKey = ""
    discDryRun = false

    // Execute the discover command.
    if err := discoverCmd.RunE(discoverCmd, []string{}); err != nil {
        t.Fatalf("discover command failed: %v", err)
    }

    // Load the resulting inventory.
    out, err := os.ReadFile(tmpFile.Name())
    if err != nil {
        t.Fatalf("reading output file: %v", err)
    }
    var result inventory.FileFormat
    if err := yaml.Unmarshal(out, &result); err != nil {
        t.Fatalf("unmarshal result: %v", err)
    }

    // Verify we have exactly one node and one cabinet.
    if len(result.Nodes) != 1 {
        t.Fatalf("expected 1 node, got %d", len(result.Nodes))
    }
    if len(result.Cabinets) != 1 {
        t.Fatalf("expected 1 cabinet, got %d", len(result.Cabinets))
    }

    // Node checks – the MAC should match the mock system NIC.
    node := result.Nodes[0]
    if node.MAC != "aa:bb:cc:dd:ee:01" {
        t.Errorf("node MAC mismatch: got %s, want aa:bb:cc:dd:ee:01", node.MAC)
    }
    // Cabinet checks – MAC and IP should come from manager info.
    cab := result.Cabinets[0]
    if cab.Xname != "x9000c1" {
        t.Errorf("cabinet xname mismatch: got %s, want x9000c1", cab.Xname)
    }
    if cab.MAC != "aa:bb:cc:dd:ee:ff" {
        t.Errorf("cabinet MAC mismatch: got %s, want aa:bb:cc:dd:ee:ff", cab.MAC)
    }
    if cab.IP != "192.168.100.10" {
        t.Errorf("cabinet IP mismatch: got %s, want 192.168.100.10", cab.IP)
    }
}
