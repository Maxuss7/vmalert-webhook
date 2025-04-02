package types

// Payload from Alertmanager
type AlertmanagerPayload struct {
	Status string  `json:"status"`
	Alerts []Alert `json:"alerts"`
}

type Alert struct {
	Status      string            `json:"status"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	StartsAt    string            `json:"startsAt"`
}

// QueryResult from VictoriaLogs
type QueryResult struct {
	Time       string `json:"_time"`
	StreamId   string `json:"_stream_id"`
	Stream     string `json:"_stream"`
	Msg        string `json:"_msg"`
	Component  string `json:"component"`
	Host       string `json:"host"`
	HostIp     string `json:"host_ip"`
	SourceType string `json:"source_type"`
	File       string `json:"file"`
}
