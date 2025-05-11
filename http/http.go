package http

import (
	"fmt"
	"net"
	"time"
)

const SIZE int = 1024 // size of the data the server will read. 1kb.

type Response struct {
	Headers *res_headers
	Body    *res_body
}

func (r *Response) join() []byte {
	rs := []byte{}
	rs = append(rs, r.Headers.join()...)
	rs = append(rs, []byte("\r\n")...)
	rs = append(rs, r.Body.join()...)
	return rs
}

type res_headers struct {
	response_line []byte
	server        []byte
	date          []byte
	cont_type     []byte
	cont_length   []byte
}

func (h *res_headers) join() []byte { // use MAP...
	r := []byte{}
	r = append(r, h.response_line...)
	r = append(r, h.server...)
	r = append(r, h.date...)
	r = append(r, h.cont_type...)
	r = append(r, h.cont_length...)
	return r
}

type res_body struct {
	body []byte
}

func (rs *res_body) join() []byte {
	return rs.body
}

func HandleConnection(conn net.Conn) {
	// close the connection when we're done
	defer conn.Close()

	// read incoming data
	buf := make([]byte, SIZE)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	y, m, d := time.Now().Date()
	// send back the data
	header := &res_headers{
		response_line: []byte("HTTP/1.1 200 OK" + "\r\n"),
		server:        []byte("Server: " + "Monzo" + "\r\n"),
		date:          []byte("Date: " + fmt.Sprintf("%v %s %v", y, m, d) + "\r\n"),
		cont_type:     []byte("Content-Type: " + "text/html" + "\r\n"),
		cont_length:   []byte("Content-Length: " + "69" + "\r\n"),
	}
	body := &res_body{
		// body: []byte("Helo World!" + "\r\n"),
		// body: []byte("<html><body><h1>Helo world!</h1></body></html>"),
		body: []byte(buf), // send back the data from the request, along with the http response header
	}
	res := &Response{
		Headers: header,
		Body:    body,
	}
	_, err = conn.Write([]byte(res.join()))
	if err != nil {
		fmt.Println(err)
		return
	}

	// print the incoming data
	fmt.Printf("Received: %s", res.join())
}
