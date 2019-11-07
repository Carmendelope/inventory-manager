/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
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
const ProxyContextTimeout = 30 * time.Second

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
