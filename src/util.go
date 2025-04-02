package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/bonzonkim/vmalert-webhook/types"
)

var (
	vlogEndpoint = os.Getenv("VICTORIALOGS_ENDPOINT")
	slackEndpoint = os.Getenv("SLACK_ENDPOINT")
)


// QueryVictoriaLogs query to VictoriaLogs endpoint which is `vlog:9428/select/logsql/query`
// then return the log stream.
func QueryVictoriaLogs(query string) ([]string, error) {
	params := url.Values{}
	params.Set("query", query)

	fullURL := vlogEndpoint + "?" + params.Encode()
	log.Printf("Querying VictoriaLogs: %s", fullURL)

	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("Response Status: %s", resp.Status)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("Response Body: %s", string(body))

	if len(body) == 0 {
		log.Println("No data returned from VictoriaLogs")
		return []string{}, nil
	}

	var logs []string
	decoder := json.NewDecoder(bytes.NewReader(body))
	for decoder.More() {
		var result types.QueryResult
		if err := decoder.Decode(&result); err != nil {
			return nil, err
		}
		logs = append(logs, result.Msg)
	}
	return logs, nil
}

// SendSlackMessage make Slack attachment with logs stream then send to slack channel using variable SLACK_ENDPOINT
func SendSlackMessage(alert types.Alert, logs []string) error {
	attachment := slack.Attachment{}

	for k, v := range alert.Labels {
		attachment.AddField(slack.Field{
			Title: k,
			Value: v,
			Short: true,
		})
	}

	desc := alert.Annotations["description"]
	if len(logs) > 0 {
		desc += "\n*Recent Logs:*\n"
		for _, line := range logs {
			desc += fmt.Sprintf("• `%s`\n", line)
		}
	}
	attachment.Text = &desc

	color := "danger"
	if strings.ToLower(alert.Status) == "resolved" {
		color = "good"
	}
	attachment.Color = &color

	payload := slack.Payload{
		Text:        fmt.Sprintf("[%s] %s", strings.ToUpper(alert.Status), alert.Labels["alertname"]),
		Attachments: []slack.Attachment{attachment},
	}

	err := slack.Send(slackEndpoint, "", payload)
	if len(err) > 0 {
		return fmt.Errorf("Slack error(s): %v", err)
	}
	return nil
}
