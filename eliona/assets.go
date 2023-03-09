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

	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v2"
	"github.com/eliona-smart-building-assistant/go-eliona/asset"
	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/log"
)

// CreateAssetsIfNecessary create all assets for specification if not already exists
func CreateAssetsIfNecessary(config apiserver.Configuration, spec kentix.DeviceInfo) error {
	for _, projectId := range conf.ProjIds(config) {
		if err := createAssetIfNecessary(config, projectId, nil, spec); err != nil {
			return fmt.Errorf("creating assets for device %s: %v", spec.Serial, err)
		}
	}
	return nil
}

// createAssetIfNecessary create asset for specification if not already exists
func createAssetIfNecessary(config apiserver.Configuration, projectId string, parentAssetId *int32, spec kentix.DeviceInfo) error {
	// Get known asset id from configuration
	assetID, err := conf.GetAssetId(context.Background(), config, projectId, spec.Serial)
	if err != nil {
		return fmt.Errorf("finding asset ID: %v", err)
	}
	if assetID != nil {
		return nil
	}

	// If no asset exists for project and configuration, create a new one
	name := name(spec)
	description := description(spec)

	newId, err := asset.UpsertAsset(api.Asset{
		ProjectId:               projectId,
		GlobalAssetIdentifier:   spec.Serial,
		Name:                    *api.NewNullableString(common.Ptr(name)),
		AssetType:               spec.AssetType,
		Description:             *api.NewNullableString(common.Ptr(description)),
		ParentLocationalAssetId: *api.NewNullableInt32(parentAssetId),
	})
	if err != nil {
		return fmt.Errorf("upserting asset into Eliona: %v", err)
	}
	if newId == nil {
		return fmt.Errorf("cannot create asset: %s", name)
	}

	// Remember the asset id for further usage
	if err := conf.InsertAsset(context.Background(), config, projectId, spec.Serial, *newId); err != nil {
		return fmt.Errorf("inserting asset to config db: %v", err)
	}

	log.Debug("eliona", "Created new asset for project %s and device %s.", projectId, spec.Serial)

	return nil
}

func name(specification kentix.DeviceInfo) string {
	return fmt.Sprintf("%s (%s)", specification.Name, specification.IPAddress)
}

func description(specification kentix.DeviceInfo) string {
	return fmt.Sprintf("%s (%s)", specification.Name, specification.Serial)
}
