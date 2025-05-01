package util

import (
	"log"
	"net/url"
	"os"
	"strings"
)

var (
	VlogEndpoint  = os.Getenv("VICTORIALOGS_ENDPOINT")
	SlackEndpoint = os.Getenv("SLACK_ENDPOINT")
	VlogIngress   = os.Getenv("VLOG_INGRESS")
)

// SetFullQueryUrl function   make full queriable url to extract Log detail from VictoriaLogs
// @endpoint: VictoriaLogs Endpoint
// @query: Extracted from Alertmanager 'query' annotation.
// return: full URL string
func SetFullQueryUrl(endpoint, query string) string {
	params := url.Values{}
	params.Set("query", query)

	fullURL := endpoint + "?" + params.Encode()

	return fullURL
}

// ConvertEndpointToIngress function   make full queriable url with Ingress url so users can see the actual Log in VictoriaLogs vmui
// @endpoint: VictoriaLogs Endpoint
// return: Ingress url.
func ConvertEndpointToIngress(endpoint string) string {
	partial := strings.Split(endpoint, "?")

	if VlogIngress == "" {
		log.Println("No Ingress is set")
		return ""
	}

	return VlogIngress + "?" + partial[1]
}
