package models

type Action int

const (
	Create Action = iota
	Update
	Delete
)
