package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/kr/pretty"
)

func main() {
	if err := play(); err != nil {
		log.Fatal(err)
	}
}
func play() error {
	fmt.Println("SwayDee")
	docFolder := "txt"
	files, err := ioutil.ReadDir(docFolder)
	if err != nil {
		return err
	}

	rawDocs := make(map[string]string)
	for _, f := range files {
		docPath := filepath.Join(docFolder, f.Name())
		buf, err := ioutil.ReadFile(docPath)
		if err != nil {
			return err
		}
		rawDocs[f.Name()] = string(buf)
		doc := NewDocument(f.Name(), string(buf))
		doc.Filters = append(doc.Filters, isUrl)
		doc.PPFuncs = append(doc.PPFuncs, strings.ToLower)
		doc.index()
		break
	}
	return nil
}

type Document struct {
	filename string
	content  string
	words    []string

	PPFuncs []PPFunc
	Filters []Filter
}

type PPFunc func(string) string
type Filter func(string) bool

func isUrl(s string) bool {
	if _, err := url.Parse(s); err != nil {
		return false
	}
	return true
}

func NewDocument(filename, content string) Document {
	return Document{filename: filename, content: content}
}

type Gram struct {
	N      int
	Tokens []string
}

func (d Document) WordFrequency() map[string]int {
	occurences := make(map[string]int)
	for _, w := range d.words {
		for _, pp := range d.PPFuncs {
			w = pp(w)
		}
		occurences[w] += 1
	}
	return occurences
}
func (d Document) Ngrams(n int) []Gram {
	return nil
}

func (d Document) wordTokenize() []string {
	return strings.FieldsFunc(d.content, func(r rune) bool {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return false
		}
		return true
	})
}
func (d *Document) index() {
	d.words = d.wordTokenize()
	pretty.Println(d.WordFrequency())
}
