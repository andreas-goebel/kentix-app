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

	"github.com/eliona-smart-building-assistant/go-utils/log"
)

// CreateAssetsIfNecessary create all assets for specification including sub specification if not already exists
func CreateAssetsIfNecessary(config apiserver.Configuration, spec kentix.DeviceInfo) error {
	for _, projectId := range conf.ProjIds(config) {
		_, err := createAssetIfNecessary(config, projectId, nil, spec)
		if err != nil {
			return fmt.Errorf("creating assets for device %s: %v", spec.Serial, err)
		}
	}
	return nil
}

// createAssetIfNecessary create asset for specification if not already exists
func createAssetIfNecessary(config apiserver.Configuration, projectId string, parentAssetId *int32, spec kentix.DeviceInfo) (*int32, error) {
	log.Debug("eliona", "Creating new asset for project %s and device %s.", projectId, spec.Serial)

	// Get known asset id from configuration
	existingId, err := conf.GetAssetId(context.Background(), config, projectId, spec.Serial)
	if err != nil {
		return nil, fmt.Errorf("finding asset ID: %v", err)
	}
	if existingId != nil {
		log.Debug("eliona", "already have asset ID %v.", existingId)
		return existingId, nil
	}

	log.Debug("eliona", "Will create new asset for project %s and device %s.", projectId, spec.Serial)

	return nil, nil
}

func name(specification kentix.DeviceInfo) string {
	return fmt.Sprintf("%s (%s)", specification.Name, specification.IPAddress)
}

func description(specification kentix.DeviceInfo) string {
	return fmt.Sprintf("%s (%s)", specification.Name, specification.Serial)
}
