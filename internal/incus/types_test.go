package incus

import (
	"strings"
	"testing"

	"github.com/lxc/incus/v7/shared/api"
)

func TestInstanceRowSearchTextIncludesUserVisibleFields(t *testing.T) {
	row := InstanceRow{
		Name:      "ubuntu",
		Type:      "container",
		Status:    "RUNNING",
		IPv4:      "10.0.3.15",
		IPv6:      "fd42::1",
		Snapshots: 2,
		Profiles:  []string{"default"},
		Location:  "server-a",
	}

	text := row.SearchText()
	for _, want := range []string{"ubuntu", "container", "running", "10.0.3.15", "fd42::1", "2", "default", "server-a"} {
		if !strings.Contains(text, want) {
			t.Fatalf("expected search text to contain %q, got %q", want, text)
		}
	}
}

func TestInstanceIPsUsesGlobalNonLoopbackAddresses(t *testing.T) {
	state := &api.InstanceState{Network: map[string]api.InstanceStateNetwork{
		"lo": {
			Type:      "loopback",
			Addresses: []api.InstanceStateNetworkAddress{{Family: "inet", Address: "127.0.0.1", Scope: "local"}},
		},
		"eth0": {
			Type: "broadcast",
			Addresses: []api.InstanceStateNetworkAddress{
				{Family: "inet", Address: "10.0.3.15", Scope: "global"},
				{Family: "inet", Address: "169.254.1.1", Scope: "link"},
				{Family: "inet6", Address: "fd42::1", Scope: "global"},
				{Family: "inet6", Address: "fe80::1", Scope: "link"},
			},
		},
		"eth1": {
			Type:      "broadcast",
			Addresses: []api.InstanceStateNetworkAddress{{Family: "inet", Address: "10.0.4.20", Scope: "global"}},
		},
	}}

	ipv4, ipv6 := instanceIPs(state)
	if ipv4 != "10.0.3.15,10.0.4.20" {
		t.Fatalf("unexpected IPv4 addresses: %q", ipv4)
	}
	if ipv6 != "fd42::1" {
		t.Fatalf("unexpected IPv6 addresses: %q", ipv6)
	}
}

func TestInstanceIPsReturnsDashWhenMissing(t *testing.T) {
	ipv4, ipv6 := instanceIPs(&api.InstanceState{})
	if ipv4 != "-" || ipv6 != "-" {
		t.Fatalf("expected missing IPs to be dashes, got %q %q", ipv4, ipv6)
	}
}

func TestInstanceToRowDefaults(t *testing.T) {
	row := instanceToRow(api.Instance{Name: "ubuntu", Status: "Running", StatusCode: api.Running, Type: "container"})
	if row.IPv4 != "-" || row.IPv6 != "-" {
		t.Fatalf("expected default IP dashes, got %q %q", row.IPv4, row.IPv6)
	}
	if row.Location != "-" {
		t.Fatalf("expected default location dash, got %q", row.Location)
	}
	if row.Image != "-" {
		t.Fatalf("expected default image dash, got %q", row.Image)
	}
}
