// SPDX-FileCopyrightText: 2025 OpenCHAMI Contributors
//
// SPDX-License-Identifier: MIT

package discover

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "time"

    "bootstrap/internal/inventory"
)

func TestDiscoverCabinets(t *testing.T) {
    // Mock Redfish manager Ethernet interfaces
    ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        switch r.URL.Path {
        case "/redfish/v1/Managers/BMC/EthernetInterfaces":
            // Return collection with one member
            w.Header().Set("Content-Type", "application/json")
            _, _ = w.Write([]byte(`{"Members":[{"@odata.id":"/redfish/v1/Managers/BMC/EthernetInterfaces/1"}]}`))
        case "/redfish/v1/Managers/BMC/EthernetInterfaces/1":
            w.Header().Set("Content-Type", "application/json")
            _, _ = w.Write([]byte(`{"Id":"1","MACAddress":"aa:bb:cc:dd:ee:ff","IPv4Addresses":[{"Address":"192.168.100.10","AddressOrigin":"Static"}]}`))
        default:
            w.WriteHeader(http.StatusNotFound)
        }
    }))
    defer ts.Close()

    // Host without scheme (host:port) for client construction
    host := strings.TrimPrefix(ts.URL, "https://")

    bmcs := []inventory.Entry{{Xname: "x9000c1s0b0", MAC: "", IP: host}}

    cabinets, err := DiscoverCabinets(bmcs, "user", "pass", true, 5*time.Second)
    if err != nil {
        t.Fatalf("DiscoverCabinets returned error: %v", err)
    }
    if len(cabinets) != 1 {
        t.Fatalf("expected 1 cabinet, got %d", len(cabinets))
    }
    cab := cabinets[0]
    if cab.Xname != "x9000c1" {
        t.Errorf("unexpected cabinet xname: %s", cab.Xname)
    }
    if cab.MAC != "aa:bb:cc:dd:ee:ff" {
        t.Errorf("unexpected cabinet MAC: %s", cab.MAC)
    }
    if cab.IP != "192.168.100.10" {
        t.Errorf("unexpected cabinet IP: %s", cab.IP)
    }
}
