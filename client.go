package deco

import (
	"crypto/md5"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"go-deco/utils"
)

const (
	// userName is the default (hardcoded) username
	userName = "admin"
)

var baseURL = url.URL{
	Scheme: "http",
	Path:   "/cgi-bin/luci/",
}

// Client is a client for sending requests to the Deco-m4 API
type Client struct {
	c        *http.Client
	aes      *utils.AESKey
	rsa      *rsa.PublicKey
	hash     string
	stok     string
	sequence uint
}

// ClientListResp is the structure of the client_list endpoint
type ClientListResp struct {
	ErrorCode int `json:"error_code"`
	Result    struct {
		ClientList []struct {
			AccessHost     string `json:"access_host"`
			ClientMesh     bool   `json:"client_mesh"`
			ClientType     string `json:"client_type"`
			ConnectionType string `json:"band5"`
			DownSpeed      uint   `json:"down_speed"`
			EnablePriority bool   `json:"enable_priority"`
			Interface      string `json:"interface"`
			IP             string `json:"ip"`
			MAC            string `json:"mac"`
			Name           string `json:"name"`
			Online         bool   `json:"online"`
			OwnerID        string `json:"owner_id"`
			RemainTime     int    `json:"remain_time"`
			SpaceID        string `json:"space_id"`
			UpSpeed        uint   `json:"up_speed"`
			WireType       string `json:"wire_type"`
		} `json:"client_list"`
	} `json:"result"`
}

// DeviceListResp is the structure of the device_list endpoint
type DeviceListResp struct {
	ErrorCode int `json:"error_code"`
	Result    struct {
		DeviceList []struct {
			DeviceIP          string   `json:"device_ip"`
			DeviceID          string   `json:"device_id,omitempty"`
			DeviceType        string   `json:"device_type"`
			NandFlash         bool     `json:"nand_flash"`
			OwnerTransfer     bool     `json:"owner_transfer,omitempty"`
			Previous          string   `json:"previous"`
			BSSID5G           string   `json:"bssid_5g"`
			BSSID2G           string   `json:"bssid_2g"`
			BSSIDSta5G        string   `json:"bssid_sta_5g"`
			BSSIDSta2G        string   `json:"bssid_sta_2g"`
			ParentDeviceID    string   `json:"parent_device_id,omitempty"`
			SoftwareVer       string   `json:"software_ver"`
			Role              string   `json:"role"`
			ProductLevel      int      `json:"product_level"`
			HardwareVer       string   `json:"hardware_ver"`
			InetStatus        string   `json:"inet_status"`
			SupportPLC        bool     `json:"support_plc"`
			MAC               string   `json:"mac"`
			SetGatewaySupport bool     `json:"set_gateway_support"`
			InetErrorMsg      string   `json:"inet_error_msg"`
			ConnectionType    []string `json:"connection_type,omitempty"`
			CustomNickname    string   `json:"custom_nickname,omitempty"`
			Nickname          string   `json:"nickname"`
			GroupStatus       string   `json:"group_status"`
			OemID             string   `json:"oem_id"`
			SignalLevel       struct {
				Band24 string `json:"band2_4"`
				Band5  string `json:"band5"`
			} `json:"signal_level"`
			DeviceModel       string `json:"device_model"`
			OversizedFirmware bool   `json:"oversized_firmware"`
			SpeedGetSupport   bool   `json:"speed_get_support,omitempty"`
			HwID              string `json:"hw_id"`
		} `json:"device_list"`
	} `json:"result"`
}

// PerfResp is the structure of the performance endpoint
type PerfResp struct {
	ErrorCode int `json:"error_code"`
	Result    struct {
		CPU float32 `json:"cpu_usage"`
		MEM float32 `json:"mem_usage"`
	} `json:"result"`
}

// New creates a new Go client for the Deco-m4 API
func New(target string) *Client {
	jar, _ := cookiejar.New(nil)
	c := &http.Client{Timeout: 10 * time.Second, Jar: jar}

	baseURL.Host = target

	return &Client{
		c: c,
	}
}

// Authenticate will generate the keys needed for the communication with the router.
func (c *Client) Authenticate(password string) error {
	c.aes = utils.GenerateAESKey()
	c.hash = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s%s", userName, password))))

	passwordKey, err := c.getPasswordKey()
	if err != nil {
		return err
	}

	sessionKey, sequence, err := c.getSessionKey()
	if err != nil {
		return err
	}
	c.rsa = sessionKey
	c.sequence = sequence

	encryptedPassword, err := utils.EncryptRsa(password, passwordKey)
	if err != nil {
		return err
	}

	loginReq := loginRequest{
		Operation: "login",
		Params: loginParams{
			Password: string(encryptedPassword),
		},
	}

	loginJSON, err := json.Marshal(loginReq)
	if err != nil {
		return err
	}
	args := EndpointArgs{
		Form: "login",
	}
	var result loginResponse
	err = c.doEncryptedPost(";stok=/login", args, loginJSON, true, &result)
	if err != nil {
		return err
	}
	if result.ErrCode != -5002 && result.Result.Stok != "" {
		c.stok = result.Result.Stok
		return nil
	}
	return fmt.Errorf("authentication failed: %w", err)
}

// Performance returns the current cpu and mem usage.
func (c *Client) Performance() (*PerfResp, error) {
	var result PerfResp
	err := c.doEncryptedPost(fmt.Sprintf(";stok=%s/admin/network", c.stok), EndpointArgs{Form: "performance"}, readBody, false, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DeviceList returns the list of connected deco APs
func (c *Client) DeviceList() (*DeviceListResp, error) {
	var result DeviceListResp
	err := c.doEncryptedPost(fmt.Sprintf(";stok=%s/admin/device", c.stok), EndpointArgs{Form: "device_list"}, readBody, false, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ClientList returns the list of connected devices
func (c *Client) ClientList() (*ClientListResp, error) {
	var result ClientListResp
	request := request{
		Operation: "read",
		Params:    map[string]interface{}{"device_mac": "default"},
	}
	jsonRequest, _ := json.Marshal(request)
	err := c.doEncryptedPost(fmt.Sprintf(";stok=%s/admin/client", c.stok), EndpointArgs{Form: "client_list"}, jsonRequest, false, &result)
	if err != nil {
		return nil, err
	}
	for index := range result.Result.ClientList {
		name, err := base64.StdEncoding.DecodeString(result.Result.ClientList[index].Name)
		if err == nil {
			result.Result.ClientList[index].Name = string(name)
		}
	}
	return &result, nil
}

// Reboot the deco nodes by mac address
func (c *Client) Reboot(macAddrs ...string) (map[string]interface{}, error) {
	var result map[string]interface{}
	var macList []map[string]string

	for _, mac := range macAddrs {
		macList = append(macList, map[string]string{
			"mac": strings.ToUpper(mac),
		})
	}

	request := request{
		Operation: "reboot",
		Params: map[string]interface{}{
			"mac_list": macList,
		},
	}

	jsonRequest, _ := json.Marshal(request)
	err := c.doEncryptedPost(fmt.Sprintf(";stok=%s/admin/device", c.stok), EndpointArgs{Form: "system"}, jsonRequest, false, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Custom lets you make a custom request
func (c *Client) Custom(path string, params EndpointArgs, body []byte) (interface{}, error) {
	var result interface{}
	err := c.doEncryptedPost(fmt.Sprintf(";stok=%s%s", c.stok, path), params, body, false, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
