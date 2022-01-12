package page

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Page struct {
	Title string
	Body string
}

func (p *Page) EncodeBody() []byte {
	return []byte(p.Body)
}

func (p *Page) ToMap() map[string]string {
	m := make(map[string]string)
	m["title"] = p.Title
	m["body"] = string(p.Body)
	return m
}

func (p *Page) Save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.EncodeBody(), 0600)
}

func NewPageFromJson(p *Page, data io.ReadCloser) (*Page, error) {
	err := json.NewDecoder(data).Decode(&p)
	if err != nil {
		log.Fatal(err)
	}
	return p, nil
}

func LoadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: string(body)}, nil
}