package bla

import "time"

type Doc struct {
	Title   string
	Date    time.Time
	Tags    []string
	Content string
}

type Page struct {
	Title   string
	Content string
}

type docsByTime []*Doc

func (s docsByTime) Len() int           { return len(s) }
func (s docsByTime) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s docsByTime) Less(i, j int) bool { return s[i].Date.After(s[j].Date) }
