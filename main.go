package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/vricap/gherver/http" // import the http package
)

const ADDR string = "127.0.0.1"
const PORT string = ":8000"

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	City string `json:"city"`
}

func main() {
	h := http.Init()
	// h.DisableLog() // to disable logging
	h.LoadStatic("/static", "./resource") // load static file

	h.Routes = []*http.Routes{
		{
			Path:   "/",
			Method: "GET",
			Handler: func(h *http.Http) {
				// set your own headers
				h.Response.Headers.SetStatusCode(200) // we didn't do any checking / validate the status code or content-type, it will simply send it as is as you wrote it
				h.Response.Headers.ContType = "text/html"

				// set response body
				h.Response.SendDocument("./resource/index.html") // send a document
				// h.Response.SendMedia("./resource/gophers.png") // or a media
			},
		},
		{
			Path:   "/foo",
			Method: "GET",
			Handler: func(h *http.Http) {
				h.Response.Headers.SetStatusCode(200)
				h.Response.Headers.ContType = "application/json"

				// send a json by struct
				person := Person{Name: "John Doe", Age: 30, City: "New York"}
				jsonData, err := json.Marshal(person)
				if err != nil {
					log.Fatalf("Error marshaling JSON: %s", err)
				}
				h.Response.Body.Byte(jsonData)

				// or directly like this
				// h.Response.Body.Byte([]byte(`{"user_id": 1, "message": "Hello world"}`))

			},
		},
	}

	// start the server
	fmt.Println("Server listening on: " + ADDR + PORT + "\n")
	h.StartServer(ADDR, PORT)
}
