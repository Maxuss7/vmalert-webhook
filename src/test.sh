#!/bin/bash

curl -X POST -H "Content-Type: application/json" -d '{"status":"firing",\
  "alerts":[{"status":"firing",\
  "labels":{"alertname":"Test-Alert",\
  "severity":"critical","\
  datasource":"victoriaLogs"},\
  "annotations":{"description":"Test alerting count: 5","query":"_time:3m * \"error\""},"startsAt":"2025-03-24T08:00:00Z"}]}' localhost:8080/webhook
