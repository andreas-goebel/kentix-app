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

package conf

import (
	"context"
	"kentix/apiserver"
	dbkentix "kentix/db/kentix"

	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/db"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// todo: unmock appname
const appname = "kentix"

// InsertConfig inserts or updates
func InsertConfig(ctx context.Context, config apiserver.Configuration) (apiserver.Configuration, error) {
	dbConfig := dbConfigFromApiConfig(config)
	err := dbConfig.Insert(ctx, db.Database(appname), boil.Infer())
	if err != nil {
		return apiserver.Configuration{}, err
	}
	return config, err
}

func dbConfigFromApiConfig(apiConfig apiserver.Configuration) (dbConfig dbkentix.Configuration) {
	dbConfig.ID = null.Int64FromPtr(apiConfig.Id).Int64
	dbConfig.Address = null.StringFrom(apiConfig.Address)
	dbConfig.APIKey = null.StringFrom(apiConfig.ApiKey)
	dbConfig.Enable = null.BoolFromPtr(apiConfig.Enable)
	dbConfig.RefreshInterval = apiConfig.RefreshInterval
	if apiConfig.RequestTimeout != nil {
		dbConfig.RequestTimeout = *apiConfig.RequestTimeout
	}
	dbConfig.Active = null.BoolFromPtr(apiConfig.Active)
	if apiConfig.ProjectIDs != nil {
		dbConfig.ProjectIds = *apiConfig.ProjectIDs
	}
	return dbConfig
}

func apiConfigFromDbConfig(dbConfig *dbkentix.Configuration) (apiConfig apiserver.Configuration) {
	apiConfig.Id = &dbConfig.ID
	apiConfig.Address = dbConfig.Address.String
	apiConfig.ApiKey = dbConfig.APIKey.String
	apiConfig.Enable = dbConfig.Enable.Ptr()
	apiConfig.RefreshInterval = dbConfig.RefreshInterval
	apiConfig.RequestTimeout = &dbConfig.RequestTimeout
	apiConfig.Active = dbConfig.Active.Ptr()
	apiConfig.ProjectIDs = common.Ptr[[]string](dbConfig.ProjectIds)
	return apiConfig
}

func dbSensorFromApiSensor(apiSensor apiserver.Sensor) (dbSensor dbkentix.Sensor) {
	dbSensor.AssetID = null.Int32FromPtr(apiSensor.AssetID)

	return dbSensor
}

func GetConfigs(ctx context.Context) ([]apiserver.Configuration, error) {
	dbConfigs, err := dbkentix.Configurations().All(ctx, db.Database(appname))
	if err != nil {
		return nil, err
	}
	var apiConfigs []apiserver.Configuration
	for _, dbConfig := range dbConfigs {
		apiConfigs = append(apiConfigs, apiConfigFromDbConfig(dbConfig))
	}
	return apiConfigs, nil
}

func GetAssetId(ctx context.Context, config apiserver.Configuration, projId string, deviceId string) (*int32, error) {
	dbAssets, err := dbkentix.Sensors(
		dbkentix.SensorWhere.ConfigurationID.EQ(null.Int64FromPtr(config.Id).Int64),
		dbkentix.SensorWhere.ProjectID.EQ(projId),
		dbkentix.SensorWhere.SerialNumber.EQ(deviceId),
	).All(ctx, db.Database(appname))
	if err != nil || len(dbAssets) == 0 {
		return nil, err
	}
	return common.Ptr(int32(dbAssets[0].ID)), nil
}

func InsertAsset(ctx context.Context, config apiserver.Configuration, projId string, SerialNumber string, assetId int32) error {
	var dbAsset dbkentix.Sensor
	dbAsset.ConfigurationID = null.Int64FromPtr(config.Id).Int64
	dbAsset.ProjectID = projId
	dbAsset.SerialNumber = SerialNumber
	dbAsset.AssetID = null.Int32From(assetId)
	return dbAsset.Insert(ctx, db.Database(appname), boil.Infer())
}

func SetConfigActiveState(ctx context.Context, config apiserver.Configuration, state bool) (int64, error) {
	return dbkentix.Configurations(
		dbkentix.ConfigurationWhere.ID.EQ(null.Int64FromPtr(config.Id).Int64),
	).UpdateAll(ctx, db.Database(appname), dbkentix.M{
		dbkentix.ConfigurationColumns.Active: state,
	})
}

func ProjIds(config apiserver.Configuration) []string {
	if config.ProjectIDs == nil {
		return []string{}
	}
	return *config.ProjectIDs
}

func IsConfigActive(config apiserver.Configuration) bool {
	return config.Active == nil || *config.Active
}

func IsConfigEnabled(config apiserver.Configuration) bool {
	return config.Enable == nil || *config.Enable
}

func SetAllConfigsInactive(ctx context.Context) (int64, error) {
	return dbkentix.Configurations().UpdateAll(ctx, db.Database(appname), dbkentix.M{
		dbkentix.ConfigurationColumns.Active: false,
	})
}
