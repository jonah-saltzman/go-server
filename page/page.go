package page

import "os"

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