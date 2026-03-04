// SPDX-FileCopyrightText: 2025 OpenCHAMI Contributors
//
// SPDX-License-Identifier: MIT

package discover

import (
    "context"
    "fmt"
    "strings"
    "time"
    "bootstrap/internal/inventory"
    "bootstrap/internal/redfish"
)

// extractCabinetXname extracts the cabinet xname (e.g., x9000c1) from a BMC xname like x9000c1s0b0.
func extractCabinetXname(bmcX string) (string, error) {
    // Find the first occurrence of "c" after the leading x<number>.
    // We assume format: x<system>c<cabinet>... (e.g., x9000c1s0b0)
    if !strings.HasPrefix(bmcX, "x") {
        return "", fmt.Errorf("invalid xname: %s", bmcX)
    }
    idx := strings.Index(bmcX, "c")
    if idx == -1 {
        return "", fmt.Errorf("no cabinet identifier in xname: %s", bmcX)
    }
    // Find end of cabinet number digits
    start := idx + 1
    end := start
    for end < len(bmcX) && bmcX[end] >= '0' && bmcX[end] <= '9' {
        end++
    }
    if end == start {
        return "", fmt.Errorf("cabinet number missing in xname: %s", bmcX)
    }
    cabinet := bmcX[:end]
    return cabinet, nil
}

// DiscoverCabinets scans the BMC entries and returns a list of unique cabinet entries, discovering MAC and IP via Redfish.
func DiscoverCabinets(bmcs []inventory.Entry, user, pass string, insecure bool, timeout time.Duration) ([]inventory.Entry, error) {
    seen := make(map[string]bool)
    var out []inventory.Entry
    for _, b := range bmcs {
        cabX, err := extractCabinetXname(b.Xname)
        if err != nil {
            continue
        }
        if seen[cabX] {
            continue
        }
        // Determine host for Redfish manager query
        host := b.IP
        if host == "" {
            host = b.Xname
        }
        // Attempt to discover manager info
        ctx, cancel := context.WithTimeout(context.Background(), timeout)
        info, err := redfish.DiscoverManagerInfo(ctx, host, user, pass, insecure, timeout)
        cancel()
        mac := ""
        ip := b.IP // fallback to BMC IP if manager IP not found
        if err == nil {
            mac = info.MAC
            if info.IP != "" {
                ip = info.IP
            }
        } else {
            // Fall back to existing BMC MAC if present
            mac = b.MAC
        }
        out = append(out, inventory.Entry{Xname: cabX, MAC: mac, IP: ip})
        seen[cabX] = true
    }
    return out, nil
}
