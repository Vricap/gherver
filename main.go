package main

import (
	"fmt"

	"github.com/vricap/gherver/http"
)

const ADDR string = "127.0.0.1"
const PORT string = ":8000"

func main() {
	h := http.Init()

	// make / set a default new response header
	h.Response.Headers.NewResponseHeader()

	// set your own headers
	h.Response.Headers.SetStatusCode(200)
	h.Response.Headers.ContType = "text/html"

	// set response body
	h.Response.Body.NewResponseBody([]byte(`
	<html>
		<head>
			<title>Test</title>
		</head>
		<body>
			<h1>Hello World!</h1>
			<p>Foo Bar</p>
		</body>
	</html>`))

	// start the server
	fmt.Println("Server listening on: " + ADDR + PORT)
	h.StartServer(ADDR, PORT)
}
