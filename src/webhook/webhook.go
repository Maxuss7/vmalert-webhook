package webhook

import (
	"fmt"
	"strings"

	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/Maxuss7/vmalert-webhook/types"
	"github.com/Maxuss7/vmalert-webhook/util"
)

// SendSlackMessage function    send Alert to Slack with Logs, vmui url
// if logs is over 70, send vmui url only.
func SendSlackMessage(alert types.Alert, logs []string, logUrl string) error {
	attachment := slack.Attachment{}

	for k, v := range alert.Labels {
		attachment.AddField(slack.Field{
			Title: k,
			Value: v,
			Short: true,
		})
	}

	desc := alert.Annotations["description"]

	if logUrl != "" {
		desc += fmt.Sprintf("\n<%s|See Logs in VMUI>", logUrl)
	}
	if len(logs) > 0 {
		desc += "\n*Recent Logs:*\n"
		max := min(len(logs), 70)
		for _, line := range logs[:max] {
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

	err := slack.Send(util.SlackEndpoint, "", payload)
	if len(err) > 0 {
		return fmt.Errorf("Slack error(s): %v", err)
	}
	return nil
}
