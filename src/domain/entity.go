package domain

type ResponseFormat int

const (
	JSON ResponseFormat = iota
	HTML ResponseFormat = iota
)

type WebTemplate struct {
	Path         string
	ListFilename string
	ViewFilename string
	EditFilename string
	NewFilename  string
}
