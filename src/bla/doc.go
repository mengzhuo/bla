package bla

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"
)

const (
	StatusPublic = "public"
)

type Doc struct {
	Title  string
	Date   time.Time
	Tags   []string
	Public bool

	Content []byte
}

func newDoc(rdr io.Reader) (d *Doc, err error) {

	scanner := bufio.NewScanner(rdr)
	var lineCnt int
	d = &Doc{}

	for scanner.Scan() && lineCnt < 5 {
		if len(scanner.Bytes()) == 0 {
			fmt.Println("scanner end:", lineCnt)
			break
		}
		lineCnt += 1

		switch lineCnt {
		case 1:
			d.Title = scanner.Text()
		case 2:
			d.Date, err = time.Parse("2006-01-02",
				strings.TrimSpace((scanner.Text())))
			if err != nil {
				return nil, err
			}
		case 3:
			d.Tags = strings.Split(scanner.Text(), ",")
			for i, v := range d.Tags {
				d.Tags[i] = strings.TrimSpace(v)
			}
		case 4:
			d.Public = (scanner.Text() == StatusPublic)
		}
	}

	if lineCnt < 2 {
		return nil, fmt.Errorf("too less header %v", d)
	}

	d.Content, err = ioutil.ReadAll(rdr)
	if err != nil {
		return nil, err
	}
	return d, nil
}

type Page struct {
	Title   string
	Content string
}

type docsByTime []*Doc

func (s docsByTime) Len() int           { return len(s) }
func (s docsByTime) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s docsByTime) Less(i, j int) bool { return s[i].Date.After(s[j].Date) }
