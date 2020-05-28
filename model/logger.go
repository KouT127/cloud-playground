package model

import (
	"cloud.google.com/go/logging"
	"encoding/json"
	"log"
)

type Entry struct {
	Message   string           `json:"message"`
	Severity  logging.Severity `json:"severity,omitempty"`
	Trace     string           `json:"logging.googleapis.com/trace,omitempty"`
	Component string           `json:"component,omitempty"`
}

func (e Entry) String() string {
	out, err := json.Marshal(e)
	if err != nil {
		log.Printf("json.Marshal: %v", err)
	}
	return string(out)
}
