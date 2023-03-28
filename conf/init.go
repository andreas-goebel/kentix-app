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

	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/db"
)

// InitConfiguration creates a default configuration to demonstrate how the eliona app should be configured. This configuration
// points to a not existing endpoint and has to be changed.
func InitConfiguration(connection db.Connection) error {
	_, err := InsertConfig(context.Background(), apiserver.Configuration{
		Address:         "http://localhost:3031",
		ApiKey:          "ikcsjhzrflwz5",
		Enable:          common.Ptr(false),
		RefreshInterval: 30,
		ProjectIDs:      &[]string{"1", "2", "3"},
	})
	_, err = InsertConfig(context.Background(), apiserver.Configuration{
		Address:         "http://localhost:3032",
		ApiKey:          "ikcsjhzrflwz5",
		Enable:          common.Ptr(false),
		RefreshInterval: 30,
		ProjectIDs:      &[]string{"1", "3"},
	})
	_, err = InsertConfig(context.Background(), apiserver.Configuration{
		Address:         "http://localhost:3033",
		ApiKey:          "ikcsjhzrflwz5",
		Enable:          common.Ptr(false),
		RefreshInterval: 30,
		ProjectIDs:      &[]string{"1"},
	})
	return err
}
