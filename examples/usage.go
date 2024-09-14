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
	readBody []byte = []byte{123, 34, 111, 112, 101, 114, 97, 116, 105, 111, 110, 34, 58, 34, 114, 101, 97, 100, 34, 125}
)

type request struct {
    Operation string                 `json:"operation,omitempty"`
    Params    map[string]interface{} `json:"params,omitempty"`
}

func main() {
	host := flag.String("host", "tplinkdeco.net", "The host address of the Deco-m4 API (default: tplinkdeco.net)")
	password := flag.String("password", "", "The password for authentication")
	flag.Parse()

	if *password == "" {
		fmt.Println("Usage: --host <host> --password <password>")
		os.Exit(1)
	}

	fmt.Printf("Host: %s, Password: %s\n", *host, *password) // Debugging line

	c := deco.New(*host)
	err := c.Authenticate(*password)
	if err != nil {
		log.Fatal(err.Error())
	}

if _, err := printPerformance(c); err != nil {
		log.Println("Error printing performance:", err)
	}
	if _, err := printDeviceList(c); err != nil {
		log.Println("Error printing device list:", err)
	}
	if _, err := printClientList(c); err != nil {
		log.Println("Error printing client list:", err)
	}
	if _, err := printWLAN(c); err != nil {
		log.Println("Error printing WLAN:", err)
	}
	if _, err := printLAN(c); err != nil {
		log.Println("Error printing LAN:", err)
	}
	if _, err := printWAN(c); err != nil {
		log.Println("Error printing WAN:", err)
	}
	if _, err := printInternet(c); err != nil {
		log.Println("Error printing internet:", err)
	}
	if _, err := printModel(c); err != nil {
		log.Println("Error printing model:", err)
	}
	if _, err := printEnviroment(c); err != nil {
		log.Println("Error printing environment:", err)
	}
	if _, err := printStatus(c); err != nil {
		log.Println("Error printing status:", err)
	}
	if _, err := printFirmware(c); err != nil {
		log.Println("Error printing firmware:", err)
	}
	if _, err := printAdvanced(c); err != nil {
		log.Println("Error printing advanced settings:", err)
	}
}

func printPerformance(c *deco.Client) (string, error){
	fmt.Println("[+] Performance")
	result, err := c.Custom("/admin/network", deco.EndpointArgs{Form: "performance"}, readBody)
	if err != nil {
		return "", err
	}
	// Print response as json
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsonData))
	fmt.Println()
	return string(jsonData), nil
}

func printDeviceList(c *deco.Client) (string, error){
	fmt.Println("[+] Device list")
	result, err := c.Custom("/admin/device", deco.EndpointArgs{Form: "device_list"}, readBody)
	if err != nil {
		return "", err
	}
	// Print response as json
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsonData))
	fmt.Println()
	return string(jsonData), nil
}

	func printClientList(c *deco.Client) (string, error){
	fmt.Println("[+] Client list")
		request := request{
		Operation: "read",
		Params:    map[string]interface{}{"device_mac": "default"},
	}
	jsonRequest, _ := json.Marshal(request)
	result, err := c.Custom("/admin/client", deco.EndpointArgs{Form: "client_list"}, jsonRequest)
	if err != nil {			
		return "", err
	}
	// Print response as json
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsonData))
	fmt.Println()
	return string(jsonData), nil
}

func printWLAN(c *deco.Client) (string, error){
	fmt.Println("[+] WLAN")
	result, err := c.Custom("/admin/wireless", deco.EndpointArgs{Form: "wlan"}, readBody)
	if err != nil {			
		return "", err
	}
	// Print response as json
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsonData))
	fmt.Println()
	return string(jsonData), nil
}

func printLAN(c *deco.Client) (string, error){
	fmt.Println("[+] LAN")
	result, err := c.Custom("/admin/network", deco.EndpointArgs{Form: "lan_ip"}, readBody)
	if err != nil {	
		return "", err
	}
	// Print response as json
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsonData))
	fmt.Println()
	return string(jsonData), nil
}

func printWAN(c *deco.Client) (string, error){
	fmt.Println("[+] WAN")
	result, err := c.Custom("/admin/network", deco.EndpointArgs{Form: "wan_ipv4"}, readBody)
	if err != nil {	
		return "", err
	}
	// Print response as json
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsonData))
	fmt.Println()
	return string(jsonData), nil
}

func printInternet(c *deco.Client) (string, error){
	fmt.Println("[+] Internet")
	result, err := c.Custom("/admin/network", deco.EndpointArgs{Form: "internet"}, readBody)
	if err != nil {	
		return "", err
	}
	// Print response as json
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsonData))
	fmt.Println()
	return string(jsonData), nil
}

func printModel(c *deco.Client) (string, error){
	fmt.Println("[+] Model")
	result, err := c.Custom("/admin/device", deco.EndpointArgs{Form: "model"}, readBody)
	if err != nil {
		return "", err
	}
	// Print response as json
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsonData))
	fmt.Println()
	return string(jsonData), nil
}

func printEnviroment(c *deco.Client) (string, error){
	fmt.Println("[+] Enviroment")
	result, err := c.Custom("/admin/system", deco.EndpointArgs{Form: "envar"}, readBody)
	if err != nil {
		return "", err
	}
	// Print response as json
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsonData))
	fmt.Println()
	return string(jsonData), nil
}

func printStatus(c *deco.Client) (string, error){
	fmt.Println("[+] Status")
	result, err := c.Custom("/admin/status", deco.EndpointArgs{Form: "all"}, readBody)
	if err != nil {	
		return "", err
	}
	// Print response as json
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsonData))
	fmt.Println()
	return string(jsonData), nil
}

func printFirmware(c *deco.Client) (string, error) {
	fmt.Println("[+] Firmware")
	result, err := c.Custom("/admin/firmware", deco.EndpointArgs{Form: "upgrade"}, readBody)
	if err != nil {	
		return "", err
	}
	// Print response as json
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsonData))
	fmt.Println()
	return string(jsonData), nil
}

func printAdvanced(c *deco.Client) (string, error){
	fmt.Println("[+] Advanced")
	result, err := c.Custom("/admin/wireless", deco.EndpointArgs{Form: "power"}, readBody)
	if err != nil {
		return "", err
	}
	// Print response as json
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsonData))
	fmt.Println()
	return string(jsonData), nil
}