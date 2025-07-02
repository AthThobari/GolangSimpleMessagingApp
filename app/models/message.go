package models

import "time"

type MessagePayload struct {
	From     string    `json:"from"`
	Messsage string    `json:"message"`
	Date     time.Time `json:"date"`
}
