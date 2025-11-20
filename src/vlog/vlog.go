package vlog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/bonzonkim/vmalert-webhook/types"
	"github.com/bonzonkim/vmalert-webhook/util"
)

const defaultLimit = 50

// QueryVlog queries VictoriaLogs, returns slice of messages and ingress url.
func QueryVlog(query string) ([]string, string, error) {
	endpoint := util.VlogEndpoint
	// append a limit param to avoid massive responses if backend supports it
	u, err := url.Parse(endpoint)
	if err == nil {
		q := u.Query()
		// only set limit if not present already
		if q.Get("limit") == "" {
			q.Set("limit", fmt.Sprintf("%d", defaultLimit))
		}
		// query param will be set by SetFullQueryUrl as well; merge back
		endpoint = u.String()
	}

	fullURL := util.SetFullQueryUrl(endpoint, query)

	client := http.Client{
		Timeout: 15 * time.Second,
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, fullURL, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return nil, fullURL, err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to query VictoriaLogs: %v", err)
		return nil, fullURL, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fullURL, fmt.Errorf("VictoriaLogs returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fullURL, err
	}
	if len(body) == 0 {
		log.Println("No data returned from VictoriaLogs")
		return []string{}, fullURL, nil
	}

	var logs []string
	decoder := json.NewDecoder(bytes.NewReader(body))
	for decoder.More() {
		var result types.QueryResult
		if err := decoder.Decode(&result); err != nil {
			// try to skip broken record and continue
			log.Printf("Failed to decode result (skipping): %v", err)
			// attempt to consume one token to avoid infinite loop
			_, _ = decoder.Token()
			continue
		}
		// if _msg empty, try other fields or skip
		if result.Msg == "" {
			// try fallback to Stream or File
			if result.Stream != "" {
				logs = append(logs, result.Stream)
			} else {
				// skip empty
				continue
			}
		} else {
			logs = append(logs, result.Msg)
		}
	}

	ingressFullURL := util.ConvertEndpointToIngress(fullURL)

	return logs, ingressFullURL, nil
}
