package initbmcs

import (
	"fmt"
	"strings"

	"bootstrap/internal/inventory"
)

func getBmcID(n int) int { return (n + 1) / 2 }
func getSlot(n int) int  { return ((n - 1) / 4) % 8 }
func getBlade(n int) int { return ((n - 1) / 2) % 2 }

func getNCXname(chassis string, n int) string {
	return fmt.Sprintf("%ss%db%d", chassis, getSlot(n), getBlade(n))
}

func getNCMAC(macStart string, n int) string {
	return fmt.Sprintf("%s:%d%d:%d0", macStart, 3, getSlot(n), getBlade(n))
}

func ParseChassisSpec(spec string) map[string]string {
	out := map[string]string{}
	if strings.TrimSpace(spec) == "" {
		return out
	}
	parts := strings.Split(spec, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		kv := strings.SplitN(p, "=", 2)
		if len(kv) != 2 {
			continue
		}
		k := strings.TrimSpace(kv[0])
		v := strings.TrimSpace(kv[1])
		if k != "" && v != "" {
			out[k] = v
		}
	}
	return out
}

// Generate creates the BMC entries for an initial inventory.
func Generate(chassis map[string]string, nodesPerChassis, nodesPerBMC, startNID int, bmcSubnetBase string) []inventory.Entry {
	var bmcs []inventory.Entry
	nid := startNID
	for c, macPref := range chassis {
		for i := nid; i < nid+nodesPerChassis; i += nodesPerBMC {
			x := getNCXname(c, i)
			ip := fmt.Sprintf("%s.%d", bmcSubnetBase, getBmcID(i))
			mac := strings.ToLower(getNCMAC(macPref, i))
			bmcs = append(bmcs, inventory.Entry{Xname: x, MAC: mac, IP: ip})
		}
		nid = nid + nodesPerChassis
	}
	return bmcs
}
