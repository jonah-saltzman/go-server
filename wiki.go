package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	""

	"github.com/julienschmidt/httprouter"
)

func viewHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	title := params.ByName("title")
	page, _ := loadPage(title)
	m := page.toMap()
	j := createJson(m)
	if j == nil {
		fmt.Println("error")
		return
	}
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Write(j)
}

func createJson(m map[string]string) ([]byte) {
	j, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return j
}

func newPageFromJson(p *page.Page, data io.ReadCloser) (*page.Page, error) {
	err := json.NewDecoder(data).Decode(&p)
	if err != nil {
		log.Fatal(err)
	}
	return p, nil
}

func createHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var p page.Page
	_, err := newPageFromJson(&p, r.Body)
	if err != nil {
		responder(w, createMessage("Error", err.Error()), http.StatusBadRequest)
		return
	}
	if p.save() != nil {
		responder(w, createMessage("Error", "Error saving new page"), http.StatusInternalServerError)
		return
	}
	m := p.toMap()
	responder(w, m, http.StatusCreated)
}

func createMessage(title string, body string) map[string]string {
	m := make(map[string]string)
	m[title] = body
	return m
}

func responder(w http.ResponseWriter, m map[string]string, s int) {
	w.WriteHeader(s)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if m != nil {
		r := createJson(m)
		w.Write(r)
		return
	}
	w.Write([]byte("Good job"))
	return
}

func loadPage(title string) (*page.Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &page.Page{Title: title, Body: string(body)}, nil
}

func main() {
	router := httprouter.New()
	router.GET("/view/:title", viewHandler)
	router.POST("/create", createHandler)
	log.Fatal(http.ListenAndServe(":8080", router))
}