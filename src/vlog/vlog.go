package vlog

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/bonzonkim/vmalert-webhook/types"
	"github.com/bonzonkim/vmalert-webhook/util"
)

// QueryVlog function    Query VictoriaLogs, return logs slice, ingressURL
func QueryVlog(query string) ([]string, string, error) {
	endpoint := util.VlogEndpoint
	fullURL := util.SetFullQueryUrl(endpoint, query)

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, fullURL, nil)
	if err != nil {
		log.Printf("Failed to Create Request %v", err)
		return nil, fullURL, err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to Query VictoriaLogs %v", err)
		return nil, fullURL, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fullURL, err
	}

	if len(body) == 0 {
		log.Println("No data return from VictoriaLogs")
		return []string{}, fullURL, nil
	}

	var logs []string
	decoder := json.NewDecoder(bytes.NewReader(body))
	for decoder.More() {
		var result types.QueryResult
		if err := decoder.Decode(&result); err != nil {
			log.Printf("Failed to decode result %v", err)
			return nil, fullURL, err
		}
		logs = append(logs, result.Msg)
	}

	ingressFullURL := util.ConvertEndpointToIngress(fullURL)

	return logs, ingressFullURL, nil
}
