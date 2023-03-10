//  This file is part of the eliona project.
//  Copyright © 2022 LEICOM iTEC AG. All Rights Reserved.
//  ______ _ _
// |  ____| (_)
// | |__  | |_  ___  _ __   __ _
// |  __| | | |/ _ \| '_ \ / _` |
// | |____| | | (_) | | | | (_| |
// |______|_|_|\___/|_| |_|\__,_|
//
//  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
//  BUT NOT LIMITED  TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
//  NON INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
//  DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package eliona

import (
	"context"
	"fmt"
	"kentix/apiserver"
	"kentix/conf"
	"kentix/kentix"
	"time"

	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v2"
	"github.com/eliona-smart-building-assistant/go-eliona/asset"
	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/log"
)

func UpsertDeviceInfo(config apiserver.Configuration, device kentix.DeviceInfo) error {
	for _, projectId := range conf.ProjIds(config) {
		err := upsertDeviceInfo(config, projectId, device)
		if err != nil {
			return err
		}
	}
	return nil
}

type deviceInfoPayload struct {
	IPAddress       string `json:"ip_address"`
	MACAddress      string `json:"mac_address"`
	FirmwareVersion string `json:"firmware_version"`
}

func upsertDeviceInfo(config apiserver.Configuration, projectId string, device kentix.DeviceInfo) error {
	log.Debug("Eliona", "Upsert data for device: config %d and device '%s'", config.Id, device.Serial)
	assetId, err := conf.GetAssetId(context.Background(), config, projectId, device.Serial)
	if err != nil {
		return err
	}
	if assetId == nil {
		return fmt.Errorf("unable to find asset ID")
	}
	return upsertData(
		api.SUBTYPE_INFO,
		*assetId,
		deviceInfoPayload{
			IPAddress:       device.IPAddress,
			MACAddress:      device.MacAddress,
			FirmwareVersion: device.Version.Firmware,
		},
	)
}

func UpsertDoorlockData(config apiserver.Configuration, doorlock kentix.DoorLock) error {
	for _, projectId := range conf.ProjIds(config) {
		err := upsertDoorlockData(config, projectId, doorlock)
		if err != nil {
			return err
		}
	}
	return nil
}

type doorlockDataPayload struct {
	SerialNumber string `json:"serial_number"`
	Name         string `json:"name"`
	DoorContact  int    `json:"door_contact"`
}

func upsertDoorlockData(config apiserver.Configuration, projectId string, doorlock kentix.DoorLock) error {
	log.Debug("Eliona", "Upsert data for doorlock: config %d and doorlock '%s'", config.Id, doorlock.Serial)
	assetId, err := conf.GetAssetId(context.Background(), config, projectId, doorlock.Serial)
	if err != nil {
		return err
	}
	if assetId == nil {
		return fmt.Errorf("unable to find asset ID")
	}
	return upsertData(
		api.SUBTYPE_INFO,
		*assetId,
		doorlockDataPayload{
			SerialNumber: doorlock.Serial,
			Name:         doorlock.Name,
			DoorContact:  doorlock.DoorContact,
		},
	)
}

func UpsertMultiSensorData(config apiserver.Configuration, sensor kentix.SensorData) error {
	for _, projectId := range conf.ProjIds(config) {
		err := upsertMultiSensorData(config, projectId, sensor)
		if err != nil {
			return err
		}
	}
	return nil
}

type sensorDataPayload struct {
	Temperature    string `json:"temperature"`
	Humidity       string `json:"humidity"`
	DewPoint       string `json:"dew_point"`
	AirPressure    string `json:"air_pressure"`
	AirQuality     string `json:"air_quality"`
	CO2            string `json:"co2"`
	CO             string `json:"co"`
	Heat           string `json:"heat"`
	ThermalImaging string `json:"ti"`
	Motion         string `json:"motion"`
	Vibration      string `json:"vibration"`
	PeopleCount    string `json:"people_count"`
}

func upsertMultiSensorData(config apiserver.Configuration, projectId string, sensor kentix.SensorData) error {
	log.Debug("Eliona", "Upserting data for MultiSensor: config %d and MultiSensor '%s'", config.Id, sensor.Name)
	assetId, err := conf.GetAssetId(context.Background(), config, projectId, "todo") // TODO: get asset ID of the device
	if err != nil {
		return err
	}
	if assetId == nil {
		return fmt.Errorf("unable to find asset ID")
	}
	return upsertData(
		api.SUBTYPE_INFO,
		*assetId,
		sensorDataPayload{
			Temperature:    sensor.Temperature.Value,
			Humidity:       sensor.Humidity.Value,
			DewPoint:       sensor.Dewpoint.Value,
			AirPressure:    sensor.AirPressure.Value,
			AirQuality:     sensor.AirQuality.Value,
			CO2:            sensor.CO2.Value,
			CO:             sensor.CO.Value,
			Heat:           sensor.Heat.Value,
			ThermalImaging: sensor.TI.Value,
			Motion:         sensor.Motion.Value,
			Vibration:      sensor.Vibration.Value,
			PeopleCount:    sensor.PeopleCount.Value,
		},
	)
}

//

func upsertData(subtype api.DataSubtype, assetId int32, payload any) error {
	var statusData api.Data
	statusData.Subtype = subtype
	now := time.Now()
	statusData.Timestamp = *api.NewNullableTime(&now)
	statusData.AssetId = assetId
	statusData.Data = common.StructToMap(payload)
	if err := asset.UpsertDataIfAssetExists[any](statusData); err != nil {
		return fmt.Errorf("upserting data: %v", err)
	}
	return nil
}
