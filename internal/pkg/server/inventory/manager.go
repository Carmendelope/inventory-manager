/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package inventory

import (
	"context"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-device-go"
	"github.com/nalej/grpc-device-manager-go"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/inventory-manager/internal/pkg/config"
	"github.com/nalej/inventory-manager/internal/pkg/entities"
	"github.com/nalej/inventory-manager/internal/pkg/server/contexts"
	"github.com/rs/zerolog/log"
	"strings"
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
	assets, err := m.listAssets(organizationID)
	if err != nil {
		return nil, err
	}

	controllers, err := m.listControllers(organizationID)
	if err != nil {
		return nil, err
	}

	return &grpc_inventory_manager_go.InventoryList{
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

// getSummaryFromAssetInfo returns totalNumCPUs, totalStorage and totalRAM
func (m * Manager) getSummaryFromAssetInfo(assetInfo *grpc_inventory_go.AssetInfo) (int32, int64, int64){
	var totalNumCPUs int32
	var totalStorage int64
	var totalRAM int64

	if  assetInfo == nil{
		return 0, 0, 0
	}

	if assetInfo.Hardware != nil {
		for _, cpu := range assetInfo.Hardware.Cpus{
			if cpu != nil {
				totalNumCPUs = totalNumCPUs + cpu.NumCores
			}
		}
		totalRAM = totalRAM + assetInfo.Hardware.InstalledRam
	}

	if assetInfo.Storage != nil {
		for _, storage := range assetInfo.Storage {
			if storage != nil {
				totalStorage = totalStorage + storage.TotalCapacity
			}
		}
	}

	return totalNumCPUs, totalStorage, totalRAM
}

// getSummaryFromAsset returns totalNumCPUs, totalStorage and totalRAM
func (m * Manager) getSummaryFromAsset(asset *grpc_inventory_manager_go.Asset) (int32, int64, int64){
	var totalNumCPUs int32
	var totalStorage int64
	var totalRAM int64

	if asset == nil {
		return 0, 0, 0
	}

	if asset.Hardware != nil {
		for _, cpu := range asset.Hardware.Cpus {
			if cpu != nil{
				totalNumCPUs = totalNumCPUs + cpu.NumCores
			}
		}
		totalRAM = totalRAM + asset.Hardware.InstalledRam
	}

	if asset.Storage != nil {
		for _, storage := range asset.Storage {
			if storage != nil {
				totalStorage = totalStorage + storage.TotalCapacity
			}
		}
	}

	return totalNumCPUs, totalStorage, totalRAM
}

func (m *Manager) Summary (organizationID *grpc_organization_go.OrganizationId) (*grpc_inventory_manager_go.InventorySummary, error) {
	inventoryList, err := m.List(organizationID)
	if err != nil {
		return nil, err
	}

	var totalNumCPUs int32 = 0
	var totalStorage int64 = 0
	var totalRAM int64 = 0

	for _, device := range inventoryList.Devices {
		cpus, storage, ram := m.getSummaryFromAssetInfo(device.AssetInfo)
		totalNumCPUs = totalNumCPUs + cpus
		totalStorage = totalStorage + storage
		totalRAM = totalRAM + ram
	}

	for _, asset := range inventoryList.Assets {
		// TODO Refactor asset to include an AssetInfo
		cpus, storage, ram := m.getSummaryFromAsset(asset)
		totalNumCPUs = totalNumCPUs + cpus
		totalStorage = totalStorage + storage
		totalRAM = totalRAM + ram
	}

	for _, ec := range inventoryList.Controllers {
		cpus, storage, ram := m.getSummaryFromAssetInfo(ec.AssetInfo)
		totalNumCPUs = totalNumCPUs + cpus
		totalStorage = totalStorage + storage
		totalRAM = totalRAM + ram
	}

	totalNumCPUs64 := int64 (totalNumCPUs)

	return &grpc_inventory_manager_go.InventorySummary{
		OrganizationId:       organizationID.OrganizationId,
		TotalNumCpu:          totalNumCPUs64,
		// Divided by 1024 to convert MB to GB
		TotalStorage:         totalStorage/int64(1024),
		TotalRam:             totalRAM/int64(1024),
	}, nil
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
		Location:           asset.Location,
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
		AssetInfo:          ec.AssetInfo,
	}
}

func (m * Manager) UpdateAsset (updateAssetRequest *grpc_inventory_go.UpdateAssetRequest) (*grpc_inventory_go.Asset, error) {
	ctx, cancel := contexts.SMContext()
	defer cancel()

	updated , err := m.assetsClient.Update(ctx, &grpc_inventory_go.UpdateAssetRequest{
		OrganizationId:       updateAssetRequest.OrganizationId,
		AssetId:              updateAssetRequest.AssetId,
		AddLabels:            updateAssetRequest.AddLabels,
		RemoveLabels:         updateAssetRequest.RemoveLabels,
		Labels:               updateAssetRequest.Labels,
		UpdateLastOpSummary:  updateAssetRequest.UpdateLastOpSummary,
		LastOpSummary:        updateAssetRequest.LastOpSummary,
		UpdateLastAlive:      updateAssetRequest.UpdateLastAlive,
		LastAliveTimestamp:   updateAssetRequest.LastAliveTimestamp,
		UpdateIp:             updateAssetRequest.UpdateIp,
		EicNetIp:             updateAssetRequest.EicNetIp,
		UpdateLocation:       updateAssetRequest.UpdateLocation,
		Location:             updateAssetRequest.Location,
	},
	)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (m * Manager) UpdateDeviceLocation (updateDeviceRequest *grpc_inventory_manager_go.UpdateDeviceLocationRequest) (*grpc_inventory_manager_go.Device, error) {
	ctx, cancel := contexts.SMContext()
	defer cancel()

	deviceInfo := strings.Split(updateDeviceRequest.AssetDeviceId,"#")
	if len(deviceInfo) != 2{
		return nil, derrors.NewInvalidArgumentError("invalid asset_device_id")
	}
	deviceGroupId := deviceInfo[0]
	deviceId := deviceInfo[1]

	updateRequest := &grpc_device_manager_go.UpdateDeviceLocationRequest{
		OrganizationId:       updateDeviceRequest.OrganizationId,
		DeviceGroupId:        deviceGroupId,
		DeviceId:             deviceId,
		UpdateLocation:       true,
		Location:             updateDeviceRequest.Location,
	}

	updated, err := m.deviceManagerClient.UpdateDeviceLocation(ctx, updateRequest)

	if err != nil {
		return nil, err
	}

	return entities.NewDeviceFromGRPC(updated), nil
}

// decomposeDeviceAssetID convert a grpc_inventory_manager_go.DeviceId into a grpc_device_go.DeviceId
// asset_device_id is cmmpose by device_group_id#device_id
func (m *Manager) decomposeDeviceAssetID (deviceID *grpc_inventory_manager_go.DeviceId) (*grpc_device_go.DeviceId, derrors.Error) {

	if deviceID == nil {
		return nil, derrors.NewInvalidArgumentError("deviceId cannot be empty")
	}
	deviceAssetId := deviceID.AssetDeviceId
	if deviceAssetId == "" {
		return nil, derrors.NewInvalidArgumentError("device_asset_id cannot be empty")
	}

	fields := strings.Split(deviceAssetId, "#")
	if len(fields) != 2  {
		return nil, derrors.NewInvalidArgumentError("device_asset_id is wrong")

	}
	return &grpc_device_go.DeviceId{
		OrganizationId:	deviceID.OrganizationId,
		DeviceGroupId:	fields[0],
		DeviceId: 		fields[1],
	}, nil

}

func (m *Manager)GetDeviceInfo( request *grpc_inventory_manager_go.DeviceId) (*grpc_inventory_manager_go.Device, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	deviceID, decErr  := m.decomposeDeviceAssetID(request)
	if decErr != nil {
		return nil, conversions.ToGRPCError(decErr)
	}

	log.Debug().Interface("device", deviceID).Msg("Get device info")

	device, err := m.deviceManagerClient.GetDevice(ctx, deviceID)
	if err != nil {
		log.Debug().Str("error", conversions.ToDerror(err).DebugReport()).Msg("error")
		return nil, err
	}
	log.Debug().Interface("device", device).Msg("device")
	return entities.NewDeviceFromGRPC(device), nil
}

