package models

type ResourceStateType string
type EventType string

const (
	ResourceStateNormal  ResourceStateType = "Normal"
	ResourceStateDeleted ResourceStateType = "Deleted"
	EventCreate = "CREATE"
	EventUpdate = "UPDATE"
	EventDelete = "DELETE"
)
