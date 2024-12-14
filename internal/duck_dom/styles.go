package duckdom

// oh boy CSS is comming

type Styles struct{
	Width int
	Height int
	Paddding int
	Maaargin int
	Background string
	Border Border
}

type BorderStyle int
const (
	Solid BorderStyle = iota
	Dashed

)
type Border struct{
	Width int
	Style BorderStyle
	Color string
}
