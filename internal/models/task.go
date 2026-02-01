package models

type Task struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Desc   string `json:"description"`
	Status string `json:"status"`
}
