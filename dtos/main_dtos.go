package dtos

import (
	"encoding/json"
)

type MainRouteResponse struct {
	Status      int                     `json:"status"`
	Message     string                  `json:"message"`
	Server      ServerDetailStatus      `json:"server"`
	Application ApplicationDetailStatus `json:"application"`
}

type ServerDetailStatus struct {
	Timestamp      string `json:"timestamp"`
	ProcessingTime string `json:"processing_time"`
	HostName       string `json:"host_name"`
}

type ApplicationDetailStatus struct {
	LogLevel       string `json:"log_level"`
	GoRoutineCount int    `json:"go_routine_count"`
}

func (r *MainRouteResponse) ToJSON() []byte {
	resp, err := json.Marshal(r)
	if err != nil {
		return []byte("{}")
	}
	return resp
}
