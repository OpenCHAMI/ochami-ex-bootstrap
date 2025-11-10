package redfish

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

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

func TestClientURLs(t *testing.T) {
	host := "example.com"
	user := "admin"
	pass := "password"
	insecure := true
	tests := []struct {
		name     string
		call     func(c *client) error
		wantPath string
	}{
		{
			name: "GET Systems",
			call: func(c *client) error {
				_, err := c.firstSystemPath(context.Background())
				return err
			},
			wantPath: "/redfish/v1/Systems",
		},
		{
			name: "GET EthernetInterfaces",
			call: func(c *client) error {
				_, err := c.listEthernetInterfaces(context.Background())
				return err
			},
			wantPath: "/redfish/v1/EthernetInterfaces",
		},
		{
			name: "POST SimpleUpdate",
			call: func(c *client) error {
				return c.post(context.Background(), "/UpdateService/Actions/SimpleUpdate", map[string]string{})
			},
			wantPath: "/redfish/v1/UpdateService/Actions/SimpleUpdate",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotPaths []string
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPaths = append(gotPaths, r.URL.Path)
				w.Header().Set("Content-Type", "application/json")
				// Return mock Redfish responses
				switch r.URL.Path {
				case "/redfish/v1/Systems":
					w.Write([]byte(`{"Members":[{"@odata.id":"/redfish/v1/Systems/1"}]}`))
				case "/redfish/v1/EthernetInterfaces":
					w.Write([]byte(`{"Members":[{"@odata.id":"/redfish/v1/EthernetInterfaces/1"}]}`))
				case "/redfish/v1/EthernetInterfaces/1":
					w.Write([]byte(`{"Id":"1","MACAddress":"aa:bb:cc:dd:ee:ff"}`))
				default:
					w.Write([]byte(`{}`))
				}
			}))
			defer ts.Close()

			c := newClient(host, user, pass, insecure, 0)
			c.base = ts.URL + "/redfish/v1"

			if err := tt.call(c); err != nil {
				t.Fatalf("call failed: %v", err)
			}
			// Check that the first request was to the expected path
			if len(gotPaths) == 0 {
				t.Fatal("no requests made")
			}
			if gotPaths[0] != tt.wantPath {
				t.Errorf("got path %q, want %q", gotPaths[0], tt.wantPath)
			}
		})
	}
}

func TestResolvePath(t *testing.T) {
	c := &client{base: "https://example.com/redfish/v1"}
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "Relative path",
			path: "/Systems",
			want: "https://example.com/redfish/v1/Systems",
		},
		{
			name: "Absolute URL",
			path: "http://other.com/Systems",
			want: "http://other.com/Systems",
		},
		{
			name: "Already resolved path",
			path: "https://example.com/redfish/v1/Systems",
			want: "https://example.com/redfish/v1/Systems",
		},
		{
			name: "Path without leading slash",
			path: "Systems",
			want: "https://example.com/redfish/v1/Systems",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := c.resolvePath(tt.path)
			if got != tt.want {
				t.Errorf("resolvePath(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}
