package main

import (
	"fmt"

	"github.com/vricap/gherver/http"
)

const ADDR string = "127.0.0.1"
const PORT string = ":8000"

func main() {
	h := http.Init()
	h.LoadStatic("/static", "./resource")

	h.Routes = []*http.Routes{
		{
			Path:   "/",
			Method: "GET",
			Handler: func(h *http.Http) {
				// set your own headers
				h.Response.Headers.SetStatusCode(200) // we didn't do any checking / validate the status code or content-type, it will simply send it as is as you wrote it
				h.Response.Headers.ContType = "text/html"

				// set response body
				h.Response.SendDocument("./resource/index.html")
				// h.Response.SendMedia("./resource/gophers.png")
			},
		},
		{
			Path:   "/foo",
			Method: "GET",
			Handler: func(h *http.Http) {
				h.Response.Headers.SetStatusCode(200)
				h.Response.Headers.ContType = "application/json"
				h.Response.Body.Byte([]byte(`{"user_id": 1, "message": "Hello world"}`))
			},
		},
	}

	// start the server
	fmt.Println("Server listening on: " + ADDR + PORT + "\n")
	h.StartServer(ADDR, PORT)
}
