/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package contexts

import (
	"context"
	"time"
)

const InventoryContextTimeout = 30 * time.Second
const AuthxContextTimeout = 30 * time.Second
const VPNContextTimeout = 30 * time.Second
const SMContextTimeout = 30 * time.Second
const ProxyContextTimeout = 60 * time.Second

// AuthxContext generates a new gRPC for authx connections
func AuthxContext() (context.Context, func()) {
	return context.WithTimeout(context.Background(), AuthxContextTimeout)
}

// VPNManagerContext generates a new gRPC context for VPN connections
func VPNManagerContext() (context.Context, func()) {
	return context.WithTimeout(context.Background(), VPNContextTimeout)
}

// InventoryContext generates a new gRPC context for inventory connections
func InventoryContext() (context.Context, func()) {
	return context.WithTimeout(context.Background(), InventoryContextTimeout)
}

// SMContext generates a new gRPC context for system model connections
func SMContext() (context.Context, func()) {
	return context.WithTimeout(context.Background(), SMContextTimeout)
}

// ProxyContext generates a new gRPC context for edge inventory proxy connections
func ProxyContext() (context.Context, func()) {
	return context.WithTimeout(context.Background(), ProxyContextTimeout)
}
