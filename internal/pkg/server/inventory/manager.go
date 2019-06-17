/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package inventory

import (
	"context"
	"github.com/nalej/grpc-device-go"
	"github.com/nalej/grpc-device-manager-go"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/inventory-manager/internal/pkg/config"
	"github.com/nalej/inventory-manager/internal/pkg/entities"
	"time"
)

const DefaultTimeout = time.Second * 10

type Manager struct {
	deviceManagerClient grpc_device_manager_go.DevicesClient
	assetsClient        grpc_inventory_go.AssetsClient
	controllersClient   grpc_inventory_go.ControllersClient
	cfg                 config.Config
}

func NewManager(deviceManagerClient grpc_device_manager_go.DevicesClient,
	assetsClient grpc_inventory_go.AssetsClient,
	controllersClient grpc_inventory_go.ControllersClient, cfg config.Config) Manager {
	return Manager{
		deviceManagerClient: deviceManagerClient,
		assetsClient:        assetsClient,
		controllersClient:   controllersClient,
		cfg:                 cfg,
	}
}

func (m *Manager) List(organizationID *grpc_organization_go.OrganizationId) (*grpc_inventory_manager_go.InventoryList, error) {

	devices, err := m.listDevices(organizationID)
	if err != nil {
		return nil, err
	}

	/*devicesIM := make([]*grpc_inventory_manager_go.Device, 0)
	for _, device := range devices {
		deviceIM := &grpc_inventory_manager_go.Device{
			OrganizationId: device.OrganizationId,
			Location: device.Location,
			Labels: device.Labels,
			RegisterSince: device.RegisterSince,
			DeviceId: device.DeviceId,
			DeviceGroupId: device.DeviceGroupId,
			AssetDeviceId: device.AssetDeviceId,

		}
		devicesIM = append(devicesIM, deviceIM)
	}
	*/
	assets, err := m.listAssets(organizationID)
	if err != nil {
		return nil, err
	}

	controllers, err := m.listControllers(organizationID)
	if err != nil {
		return nil, err
	}

	return &grpc_inventory_manager_go.InventoryList{
		//Devices:     devicesIM,
		Devices: devices,
		Assets:      assets,
		Controllers: controllers,
	}, nil
}

func (m *Manager) listDevices(organizationID *grpc_organization_go.OrganizationId) ([]*grpc_inventory_manager_go.Device, error) {
	ctxDM, cancelDM := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancelDM()
	groups, err := m.deviceManagerClient.ListDeviceGroups(ctxDM, organizationID)
	if err != nil {
		return nil, err
	}
	result := make([]*grpc_inventory_manager_go.Device, 0)
	for _, deviceGroup := range groups.Groups {
		ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
		defer cancel()
		deviceGroupId := &grpc_device_go.DeviceGroupId{
			OrganizationId: deviceGroup.OrganizationId,
			DeviceGroupId:  deviceGroup.DeviceGroupId,
		}
		devices, err := m.deviceManagerClient.ListDevices(ctx, deviceGroupId)
		if err != nil {
			return nil, err
		}
		for _, dev := range devices.Devices {
			result = append(result, entities.NewDeviceFromGRPC(dev))
		}

	}
	return result, nil
}

func (m *Manager) listAssets(organizationID *grpc_organization_go.OrganizationId) ([]*grpc_inventory_manager_go.Asset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	assets, err := m.assetsClient.List(ctx, organizationID)
	if err != nil {
		return nil, err
	}
	return m.toAssetFromList(assets.Assets), nil
}

func (m *Manager) listControllers(organizationID *grpc_organization_go.OrganizationId) ([]*grpc_inventory_manager_go.EdgeController, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	controllers, err := m.controllersClient.List(ctx, organizationID)
	if err != nil {
		return nil, err
	}
	result := make([]*grpc_inventory_manager_go.EdgeController, 0)
	for _, ec := range controllers.Controllers {
		toAdd := m.toController(ec)
		result = append(result, toAdd)
	}
	return result, nil
}

func (m *Manager) GetControllerExtendedInfo(edgeControllerID *grpc_inventory_go.EdgeControllerId) (*grpc_inventory_manager_go.EdgeController, []*grpc_inventory_manager_go.Asset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	controller, err := m.controllersClient.Get(ctx, edgeControllerID)
	if err != nil {
		return nil, nil, err
	}
	assetCtx, assetCancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer assetCancel()
	assets, err := m.assetsClient.ListControllerAssets(assetCtx, edgeControllerID)
	if err != nil {
		return nil, nil, err
	}
	return m.toController(controller), m.toAssetFromList(assets.Assets), nil
}

func (m *Manager) GetAssetInfo(assetID *grpc_inventory_go.AssetId) (*grpc_inventory_manager_go.Asset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	asset, err := m.assetsClient.Get(ctx, assetID)
	if err != nil {
		return nil, err
	}
	return m.toAsset(asset), nil
}

func (m *Manager) Summary(organizationID *grpc_organization_go.OrganizationId) (*grpc_inventory_manager_go.InventorySummary, error) {
	panic("implement me")
}

func (m *Manager) toAssetFromList(assets []*grpc_inventory_go.Asset) []*grpc_inventory_manager_go.Asset {
	result := make([]*grpc_inventory_manager_go.Asset, 0)
	for _, asset := range assets {
		toAdd := m.toAsset(asset)
		result = append(result, toAdd)
	}
	return result
}

func (m *Manager) toAsset(asset *grpc_inventory_go.Asset) *grpc_inventory_manager_go.Asset {
	status := grpc_inventory_manager_go.ConnectedStatus_OFFLINE
	if asset.LastAliveTimestamp != 0 {
		timeCalculated := time.Unix(asset.LastAliveTimestamp, 0).Add(m.cfg.AssetThreshold).Unix()
		if timeCalculated > time.Now().Unix() {
			status = grpc_inventory_manager_go.ConnectedStatus_ONLINE
		}
	}

	return &grpc_inventory_manager_go.Asset{
		OrganizationId:     asset.OrganizationId,
		EdgeControllerId:   asset.EdgeControllerId,
		AssetId:            asset.AssetId,
		AgentId:            asset.AgentId,
		Show:               asset.Show,
		Created:            asset.Created,
		Labels:             asset.Labels,
		Os:                 asset.Os,
		Hardware:           asset.Hardware,
		Storage:            asset.Storage,
		EicNetIp:           asset.EicNetIp,
		LastOpSummary:      asset.LastOpResult,
		LastAliveTimestamp: asset.LastAliveTimestamp,
		Status:             status,
	}
}

func (m *Manager) toController(ec *grpc_inventory_go.EdgeController) *grpc_inventory_manager_go.EdgeController {
	status := grpc_inventory_manager_go.ConnectedStatus_OFFLINE
	if ec.LastAliveTimestamp != 0 {
		timeCalculated := time.Unix(ec.LastAliveTimestamp, 0).Add(m.cfg.ControllerThreshold).Unix()
		if timeCalculated > time.Now().Unix() {
			status = grpc_inventory_manager_go.ConnectedStatus_ONLINE
		}
	}
	return &grpc_inventory_manager_go.EdgeController{
		OrganizationId:     ec.OrganizationId,
		EdgeControllerId:   ec.EdgeControllerId,
		Show:               ec.Show,
		Created:            ec.Created,
		Name:               ec.Name,
		Labels:             ec.Labels,
		LastAliveTimestamp: ec.LastAliveTimestamp,
		Status:             status,
		Location:           ec.Location,
	}
}
