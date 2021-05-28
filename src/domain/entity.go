package domain

type ResponseFormat int

const (
	JSON ResponseFormat = iota
	HTML ResponseFormat = iota
)
