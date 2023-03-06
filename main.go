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

package main

import (
	"context"
	"kentix/conf"
	"kentix/eliona"
	"time"

	"github.com/eliona-smart-building-assistant/go-eliona/app"
	"github.com/eliona-smart-building-assistant/go-eliona/asset"
	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/db"
	"github.com/eliona-smart-building-assistant/go-utils/log"
)

// The main function starts the app by starting all services necessary for this app and waits
// until all services are finished.
func main() {
	log.Info("Kentix", "Starting the app.")

	// Necessary to close used init resources, because db.Pool() is used in this app.
	defer db.ClosePool()

	// Init the app before the first run.
	app.Init(db.Pool(), app.AppName(),
		app.ExecSqlFile("conf/init.sql"),
		asset.InitAssetTypeFile("eliona/asset-type-alarm-manager.json"),
		asset.InitAssetTypeFile("eliona/asset-type-access-manager.json"),
		asset.InitAssetTypeFile("eliona/asset-type-doorlock.json"),
		asset.InitAssetTypeFile("eliona/asset-type-multi-sensor.json"),
		conf.InitConfiguration,
		eliona.InitEliona,
	)

	common.WaitFor(
		common.Loop(collectData, time.Second),
		listenApiRequests,
	)

	// At the end set all configuration inactive
	_, err := conf.SetAllConfigsInactive(context.Background())
	if err != nil {
		log.Error("Kentix", "setting all configs inactive: %v", err)
	}

	log.Info("Kentix", "Terminating the app.")
}
