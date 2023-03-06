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
	"io"
	"kentix/apiserver"
	"log"

	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type AccessManagerInfoResponse struct {
	Data AccessManager `json:"data"`
}

type AccessManager struct {
	Name        string      `json:"name"`
	IPAddress   string      `json:"ip_address"`
	MacAddress  string      `json:"mac_address"`
	Type        int         `json:"type"`
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

func GetAccessManager(conf apiserver.Configuration) (*AccessManager, error) {
	client := &http.Client{}

	url, err := url.JoinPath(conf.Address, "api/info")
	if err != nil {
		return nil, fmt.Errorf("appending endpoint to URL: %v", err)
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %v", err)
	}

	log.Printf("%+v", string(body))

	var accessManagerInfoResponse AccessManagerInfoResponse
	err = json.Unmarshal(body, &accessManagerInfoResponse)
	if err != nil {
		return nil, fmt.Errorf("parsing response: %v", err)
	}

	return &accessManagerInfoResponse.Data, nil
}
