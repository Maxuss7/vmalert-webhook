package main

import (
	"log"
	"net/http"

	"github.com/Maxuss7/vmalert-webhook/types"
	"github.com/Maxuss7/vmalert-webhook/vlog"
	"github.com/Maxuss7/vmalert-webhook/webhook"
	"github.com/gin-gonic/gin"
)

// main function    Run Server on Port 8080
func main() {
	r := gin.Default()

	r.POST("/webhook", func(c *gin.Context) {
		log.Println("Endpoint /webhook hit")
		var payload types.AlertmanagerPayload
		if err := c.ShouldBindJSON(&payload); err != nil {
			log.Printf("Failed to parse payload: %v", err)
			c.Status(http.StatusBadRequest)
			return
		}

		for _, alert := range payload.Alerts {
			query := alert.Annotations["query"]
			desc := alert.Annotations["description"]
			log.Printf("[ALERT] Status: %s | Desc: %s | Query: %s", alert.Status, desc, query)

			logs, ingressURL, err := vlog.QueryVlog(query)
			if err != nil {
				log.Printf("Failed to query VictoriaLogs: %v", err)
				continue
			}

			if err := webhook.SendSlackMessage(alert, logs, ingressURL); err != nil {
				log.Printf("Failed to send Slack message: %v", err)
			}
		}
		c.Status(http.StatusOK)
	})

	log.Println("Starting server on :8080")
	r.Run(":8080")
}
