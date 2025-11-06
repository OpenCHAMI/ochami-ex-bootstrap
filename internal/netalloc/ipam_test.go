package netalloc

import (
	"strings"
	"testing"
)

func TestAllocatorGatewayReserved(t *testing.T) {
	a, err := NewAllocator("10.0.0.0/24")
	if err != nil {
		t.Fatalf("NewAllocator: %v", err)
	}
	ip1, err := a.Next()
	if err != nil {
		t.Fatalf("Next: %v", err)
	}
	if !strings.HasPrefix(ip1, "10.0.0.") {
		t.Fatalf("unexpected ip1: %s", ip1)
	}
	if ip1 != "10.0.0.2" {
		t.Fatalf("expected first allocation to skip .1 and be .2, got %s", ip1)
	}
}

func TestAllocatorReserveAndNext(t *testing.T) {
	a, err := NewAllocator("10.0.1.0/29") // hosts .1-.6, gw .1 reserved
	if err != nil {
		t.Fatalf("NewAllocator: %v", err)
	}
	// First Next should be .2
	if ip, _ := a.Next(); ip != "10.0.1.2" {
		t.Fatalf("got %s want 10.0.1.2", ip)
	}
	// Reserve .4, ensure sequence does not return .4
	a.Reserve("10.0.1.4")
	// Next should be .3
	if ip, _ := a.Next(); ip != "10.0.1.3" {
		t.Fatalf("got %s want 10.0.1.3", ip)
	}
	// Next should skip .4 (reserved) and return .5
	if ip, _ := a.Next(); ip != "10.0.1.5" {
		t.Fatalf("got %s want 10.0.1.5", ip)
	}
}
