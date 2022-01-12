package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	. "github.com/jonah-saltzman/go-server/page"
)

func viewHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	title := params.ByName("title")
	page, err := LoadPage(title)
	if err != nil {
		responder(writer, createMessage("Error", "Page not found"), http.StatusNotFound)
		return
	}
	m := page.ToMap()
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

func createHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var p Page
	_, err := NewPageFromJson(&p, r.Body)
	if err != nil {
		responder(w, createMessage("Error", err.Error()), http.StatusBadRequest)
		return
	}
	if p.Save() != nil {
		responder(w, createMessage("Error", "Error saving new page"), http.StatusInternalServerError)
		return
	}
	m := p.ToMap()
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

func main() {
	router := httprouter.New()
	router.GET("/view/:title", viewHandler)
	router.POST("/create", createHandler)
	log.Fatal(http.ListenAndServe(":8080", router))
}