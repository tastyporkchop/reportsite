package usgsatom

import (
	"golang.org/x/tools/blog/atom"
)

type USGSFeed struct {
	atom.Feed
	Entry []*USGSEntry `xml:"entry"`
}

type USGSEntry struct {
	atom.Entry
	Point    string     `xml:"point"`
	Elev     string     `xml:"elev"`
	Category []Category `xml:"category"`
}

type Category struct {
	Label string `xml:"label,attr"`
	Term  string `xml:"term,attr"`
}
