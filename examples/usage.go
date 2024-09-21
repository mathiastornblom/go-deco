package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"go-deco"
)

var (
	// {"operation": "read"}
	readBody = []byte(`{"operation": "read"}`)
)

var applicationVersion = "1.0.3"

type request struct {
	Operation string                 `json:"operation,omitempty"`
	Params    map[string]interface{} `json:"params,omitempty"`
}

// Struktur för att läsa av enhetens lista från Device list
type DeviceListResponse struct {
	ErrorCode int `json:"error_code"`
	Result    struct {
		DeviceList []struct {
			DeviceIP string `json:"device_ip"`
			DeviceMAC string `json:"mac"`
			Nickname string `json:"nickname"`
		} `json:"device_list"`
	} `json:"result"`
}

func main() {
	host := flag.String("host", "tplinkdeco.net", "The host address of the Deco-m4 API (default: tplinkdeco.net)")
	password := flag.String("password", "", "The password for authentication")
	version := flag.Bool("version", false, "Print the application version")
	flag.Parse()

	if *version {
		fmt.Printf("Application version: %s\n", applicationVersion)
		os.Exit(0)
	}

	if *password == "" {
		fmt.Println("Usage: --host <host> --password <password>")
		fmt.Println("Usage: --version")
		os.Exit(1)
	}

	fmt.Printf("Host: %s, Password: %s\nVersion: %s\n", *host, *password, applicationVersion)

	c := deco.New(*host)
	err := c.Authenticate(*password)
	if err != nil {
		log.Fatalf("Authentication failed: %s", err.Error())
	}

	// Hämta device list
	deviceListJSON, err := fetchDeviceList(c)
	if err != nil {
		log.Fatalf("Failed to fetch device list: %s", err.Error())
	}

	// Unmarshal JSON till DeviceListResponse
	var deviceList DeviceListResponse
	err = json.Unmarshal([]byte(deviceListJSON), &deviceList)
	if err != nil {
		log.Fatalf("Failed to unmarshal device list response: %s", err.Error())
	}

	// Kolla efter eventuella felkoder
	if deviceList.ErrorCode != 0 {
		log.Fatalf("Error in device list response: %d", deviceList.ErrorCode)
	}

		// Iterera över varje enhet i listan
	for _, device := range deviceList.Result.DeviceList {
		fmt.Printf("Found device: %s (IP: %s)\n", device.Nickname, device.DeviceIP)
	c := deco.New(device.DeviceIP)
	err := c.Authenticate(*password)
	if err != nil {
		log.Fatalf("Authentication failed: %s", err.Error())
	}
		endpoints := []struct {
			description string
			endpoint    string
			form        string
			body        []byte
		}{
			{"Device list", "/admin/device", "device_list", readBody},
			{"Performance", "/admin/network", "performance", readBody},
			{"Client list", "/admin/client", "client_list", createRequestBody("default")},
			{"WLAN", "/admin/wireless", "wlan", readBody},
			{"LAN ipv4", "/admin/network", "lan_ip", createRequestBody("default")},
			{"LAN ipv6", "/admin/network", "ipv6", createRequestBody("default")},
			{"WAN", "/admin/network", "wan_ipv4", readBody},
			{"Internet", "/admin/network", "internet", readBody},
			{"Mode", "/admin/device", "mode", readBody},
			{"Advanced", "/admin/wireless", "power", readBody},
			{"DHCP Dial", "/admin/network", "dhcp_dial", readBody},
		}

		for _, ep := range endpoints {
			printEndpointData(c, ep.description, ep.endpoint, ep.form, ep.body)
		}
	}
}

// Funktion för att hämta Device List från API
func fetchDeviceList(c *deco.Client) (string, error) {
	fmt.Println("[+] Fetching device list")
	result, err := c.Custom("/admin/device", deco.EndpointArgs{Form: "device_list"}, readBody)
	if err != nil {
		return "", err
	}

	// Returnera som JSON-sträng
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	return string(jsonData), nil
}

func createRequestBody(deviceMAC string) []byte {
	request := request{
		Operation: "read",
		Params:    map[string]interface{}{"device_mac": deviceMAC},
	}
	jsonRequest, _ := json.Marshal(request)
	return jsonRequest
}

func printEndpointData(c *deco.Client, description, endpoint, form string, body []byte) {
	fmt.Printf("[+] %s\n", description)
	jsonData, err := fetchAndPrintJSON(c, endpoint, form, body)
	if err != nil {
		log.Printf("Error printing %s: %s", description, err.Error())
		return
	}
	fmt.Println(jsonData)
	fmt.Println()
}

func fetchAndPrintJSON(c *deco.Client, endpoint, form string, body []byte) (string, error) {
	result, err := c.Custom(endpoint, deco.EndpointArgs{Form: form}, body)
	if err != nil {
		return "", err
	}

	jsonData, _ := json.MarshalIndent(result, "", "  ")
	return string(jsonData), nil
}
