package dto

import "encoding/json"

type BaseMessage struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}

type DriverLocationMessage struct {
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	CarPackage string  `json:"carPackage"`
}
