package class

import "encoding/xml"

type Response[T any] struct {
	Status string `xml:"status,attr"`
	Response T `xml:"result"`
}

type Result[T any] struct {
	XMLName     xml.Name `xml:"result"`
	TotalCount  int      `xml:"total-count,attr"`
	Count       int      `xml:"count,attr"`
	Entry       T      `xml:"entry"`
}
