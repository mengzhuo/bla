package bla

import "github.com/russross/blackfriday"

func parseContent(p []byte) (output []byte) {
	blackfriday.MarkdownCommon(p)
}
