package main

import "time"

type LogData struct {
	Level   string    `json:"level"`
	Message string    `json:"message"`
	From    string    `json:"from"`
	Time    time.Time `json:"logged_at"`
}
