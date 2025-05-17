package main

import (
	"fmt"

	"github.com/vricap/gherver/http"
)

const ADDR string = "127.0.0.1"
const PORT string = ":8000"

func main() {
	h := http.Init()

	h.Routes = []*http.Routes{
		{
			Path:   "/",
			Method: "GET",
			Handler: func(h *http.Http) {
				// set your own headers
				h.Response.Headers.SetStatusCode(200) // we didn't do any checking / validate the status code or content-type, it will simply send it as is as you wrote it
				h.Response.Headers.ContType = "text/html"

				// set response body
				h.Response.SendHtml("./resource/index.html")
			},
		},
		{
			Path:   "/foo",
			Method: "GET",
			Handler: func(h *http.Http) {
				// set your own headers
				h.Response.Headers.SetStatusCode(200)
				h.Response.Headers.ContType = "text/html"

				// set response body
				h.Response.Body.NewResponseBody(`
<html>
	<head>
		<title>Foo</title>
	</head>
	<body>
		<h1>Foo!</h1>
		<p>Foo Bar</p>
	</body>
</html>`)
			},
		},
	}

	// start the server
	fmt.Println("Server listening on: " + ADDR + PORT + "\n")
	h.StartServer(ADDR, PORT)
}
