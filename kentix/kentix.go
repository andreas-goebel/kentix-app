//  This file is part of the eliona project.
//  Copyright © 2023 LEICOM iTEC AG. All Rights Reserved.
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

package kentix

import (
	"fmt"
	"kentix/apiserver"
	"net/url"
	"time"

	"github.com/eliona-smart-building-assistant/go-utils/http"
)

const (
	AccessPointAssetType  = "Kentix Accessmanager"
	AlarmManagerAssetType = "Kentix AlarmManager"
	MultiSensorAssetType  = "Kentix MultiSensor"
	DoorlockAssetType     = "Kentix Doorlock"
)

type infoResponse struct {
	Data DeviceInfo `json:"data"`
}

type DeviceInfo struct {
	Name        string      `json:"name"`
	IPAddress   string      `json:"ip_address"`
	MacAddress  string      `json:"mac_address"`
	Type        int         `json:"type"`
	AssetType   string      `json:"-"`
	Serial      string      `json:"serial"`
	Version     VersionInfo `json:"version"`
	OSRevision  int         `json:"os_revision"`
	BootedAt    int         `json:"booted_at"`
	LastBackup  interface{} `json:"last_backup"`
	MasterSlave MasterSlave `json:"masterslave"`
}

type VersionInfo struct {
	Firmware string `json:"firmware"`
	Atmel    string `json:"atmel"`
	FSM      string `json:"fsm"`
	GSM      string `json:"gsm"`
}

type MasterSlave struct {
	IsSlave  bool   `json:"is_slave"`
	MasterIP string `json:"master_ip"`
}

func GetDeviceInfo(conf apiserver.Configuration) (*DeviceInfo, error) {
	url, err := url.JoinPath(conf.Address, "api/info")
	if err != nil {
		return nil, fmt.Errorf("appending endpoint to URL: %v", err)
	}
	r, err := http.NewRequestWithApiKey(url, "Authorization", "Basic "+conf.ApiKey)
	if err != nil {
		return nil, fmt.Errorf("creating request to %s: %v", url, err)
	}
	infoResponse, err := http.Read[infoResponse](r, time.Duration(*conf.RequestTimeout)*time.Second, true)
	if err != nil {
		return nil, fmt.Errorf("reading response from %s: %v", url, err)
	}
	infoResponse.Data.AssetType, err = inferAssetType(infoResponse.Data.Type)
	if err != nil {
		return nil, fmt.Errorf("inferring asset type from %s: %v", url, err)
	}
	return &infoResponse.Data, nil
}

func inferAssetType(typ int) (string, error) {
	// TODO: Verify that this is the correct property to determine device type.
	switch typ {
	case 1:
		return AccessPointAssetType, nil
	case 8:
		return AlarmManagerAssetType, nil
	case 110:
		return MultiSensorAssetType, nil
	}
	return "", fmt.Errorf("unknown device type: %v", typ)
}

type accessPointResponse struct {
	Data  []DoorLock     `json:"data"`
	Links PaginationLink `json:"links"`
}

type DoorLock struct {
	ID                    int    `json:"id"`
	Name                  string `json:"name"`
	Serial                string `json:"serial"`
	Active                bool   `json:"active"`
	Default               bool   `json:"default"`
	Address               string `json:"address"`
	Type                  string `json:"type"`
	DoubleAuth            int    `json:"double_auth"`
	CoupleTime            int    `json:"couple_time"`
	SwitchAlarmZone       string `json:"switch_alarmzone"`
	TimeControlled        bool   `json:"time_controlled"`
	TimeControlledBooking bool   `json:"time_controlled_booking"`
	TimeProfileID         int    `json:"timeprofile_id"`
	SwitchOutput          int    `json:"switch_output"`
	SwitchExternal        int    `json:"switch_external"`
	CameraID              int    `json:"camera_id"`
	DoorContact           int    `json:"door_contact"`
	AlarmDelay            int    `json:"alarm_delay"`
}

type PaginationLink struct {
	Next string `json:"next"`
}

func GetAccessPointReadings(conf apiserver.Configuration) ([]DoorLock, error) {
	url, err := url.JoinPath(conf.Address, "api/devices/doorlocks")
	if err != nil {
		return nil, fmt.Errorf("appending endpoint to URL: %v", err)
	}
	return fetchDoorlocks(url, conf)
}

func fetchDoorlocks(url string, conf apiserver.Configuration) ([]DoorLock, error) {
	r, err := http.NewRequestWithApiKey(url, "Authorization", "Basic "+conf.ApiKey)
	if err != nil {
		return nil, fmt.Errorf("creating request to %s: %v", url, err)
	}
	accessPointResponse, err := http.Read[accessPointResponse](r, time.Duration(*conf.RequestTimeout)*time.Second, true)
	if err != nil {
		return nil, fmt.Errorf("reading response from %s: %v", url, err)
	}
	doorlocks := accessPointResponse.Data
	if accessPointResponse.Links.Next != "" {
		dl, err := fetchDoorlocks(accessPointResponse.Links.Next, conf)
		if err != nil {
			return nil, err
		}
		doorlocks = append(doorlocks, dl...)
	}
	return doorlocks, nil
}

type DigitalInput struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Input struct {
		Value      string `json:"value"`
		HasAlarm   bool   `json:"has_alarm"`
		Unit       string `json:"unit"`
		UnitPrefix string `json:"unit_prefix"`
	} `json:"input"`
}

type SensorState struct {
	HasAlarm bool `json:"has_alarm"`
}

type SensorValue struct {
	Value    string `json:"value"`
	HasAlarm bool   `json:"has_alarm"`
}

type SensorData struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	State       SensorState `json:"state"`
	Temperature SensorValue `json:"temperature"`
	Humidity    SensorValue `json:"humidity"`
	Dewpoint    SensorValue `json:"dewpoint"`
	AirPressure SensorValue `json:"air_pressure"`
	AirQuality  SensorValue `json:"air_quality"`
	CO2         SensorValue `json:"co2"`
	Heat        SensorValue `json:"heat"`
	CO          SensorValue `json:"co"`
	TI          SensorValue `json:"ti"`
	Motion      SensorValue `json:"motion"`
	Vibration   SensorValue `json:"vibration"`
	PeopleCount SensorValue `json:"people_count"`
}

type sensorResponse struct {
	Data SensorData `json:"data"`
}

func GetMultiSensorReadings(conf apiserver.Configuration) (*SensorData, error) {
	url, err := url.JoinPath(conf.Address, "api/devices/multisensor/values")
	if err != nil {
		return nil, fmt.Errorf("appending endpoint to URL: %v", err)
	}
	r, err := http.NewRequestWithApiKey(url, "Authorization", "Basic "+conf.ApiKey)
	if err != nil {
		return nil, fmt.Errorf("creating request to %s: %v", url, err)
	}
	sensorResponse, err := http.Read[sensorResponse](r, time.Duration(*conf.RequestTimeout)*time.Second, true)
	if err != nil {
		return nil, fmt.Errorf("reading response from %s: %v", url, err)
	}
	return &sensorResponse.Data, nil
}
