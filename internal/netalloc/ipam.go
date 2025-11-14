// SPDX-FileCopyrightText: 2025 OpenCHAMI Contributors
//
// SPDX-License-Identifier: MIT

// Package netalloc provides IP address allocation utilities.
package netalloc

import (
	"context"
	"net"

	ipam "github.com/metal-stack/go-ipam"
)

// Allocator manages IP address allocation within a specified subnet.
type Allocator struct {
	ipm    ipam.Ipamer
	prefix *ipam.Prefix
}

// NewAllocator creates a new Allocator for the given CIDR subnet.
func NewAllocator(cidr string) (*Allocator, error) {
	ctx := context.Background()
	ipm := ipam.New(ctx)
	pr, err := ipm.NewPrefix(ctx, cidr)
	if err != nil {
		return nil, err
	}
	// Previously we reserved the first host (gateway) to avoid collisions.
	// Removing that reservation allows allocation of the .1 address when desired.
	return &Allocator{ipm: ipm, prefix: pr}, nil
}

// Reserve marks the specified IP address as reserved in the allocator.
func (a *Allocator) Reserve(ip string) {
	_, _ = a.ipm.AcquireSpecificIP(context.Background(), a.prefix.Cidr, ip)
}

// Next allocates and returns the next available IP address in the subnet.
func (a *Allocator) Next() (string, error) {
	addr, err := a.ipm.AcquireIP(context.Background(), a.prefix.Cidr)
	if err != nil {
		return "", err
	}
	return addr.IP.String(), nil
}

// Contains checks if the given IP address is within the allocator's subnet.
func (a *Allocator) Contains(ip string) bool {
	_, n, err := net.ParseCIDR(a.prefix.Cidr)
	if err != nil {
		return false
	}
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}
	return n.Contains(parsedIP)
}
