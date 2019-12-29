package main

import "encoding/json"

type event struct {
	Type    string
	Topic   string
	Payload *payload
}

func (e *event) UnmarshalJSON(data []byte) error {
	m := make(map[string]interface{})
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	if v, ok := m["type"]; ok {
		e.Type = v.(string)
	}

	if v, ok := m["topic"]; ok {
		e.Topic = v.(string)
	}

	if v, ok := m["payload"]; ok {
		// rm := json.RawMessage(v.([]byte))
		e.Payload = new(payload)
		if err := json.Unmarshal([]byte(v.(string)), e.Payload); err != nil {
			return err
		}
	}

	return nil
}

type payload struct {
	Type     string `json:"type"`
	Value    string `json:"value"`
	OldType  string `json:"oldType"`
	OldValue string `json:"oldValue"`
}
