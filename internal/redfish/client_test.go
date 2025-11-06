package redfish

import "testing"

func TestIsBootable_UefiPXE(t *testing.T) {
	nic := rfEthernetInterface{UefiDevicePath: "VenHw(PXE)"}
	if !isBootable(nic) {
		t.Fatal("expected bootable due to UEFI PXE")
	}
}

func TestIsBootable_DHCPOrigin(t *testing.T) {
	nic := rfEthernetInterface{IPv4Addresses: []struct {
		Address string "json:\"Address\""
		Origin  string "json:\"AddressOrigin\""
	}{{Address: "10.0.0.2", Origin: "DHCP"}}}
	if !isBootable(nic) {
		t.Fatal("expected bootable due to DHCP origin")
	}
}

func TestIsBootable_MACEnabled(t *testing.T) {
	nic := rfEthernetInterface{MACAddress: "AA:BB:CC:DD:EE:FF"}
	if !isBootable(nic) {
		t.Fatal("expected bootable with MAC and default enabled")
	}
}

func TestIsBootable_MACDisabled(t *testing.T) {
	enabled := false
	nic := rfEthernetInterface{MACAddress: "AA:BB:CC:DD:EE:FF", InterfaceEnabled: &enabled}
	if isBootable(nic) {
		t.Fatal("expected not bootable when interface disabled")
	}
}

func TestIsBootable_False(t *testing.T) {
	if isBootable(rfEthernetInterface{}) {
		t.Fatal("expected not bootable for empty NIC")
	}
}
