package http

import (
	"fmt"
	"net"
	"strconv"
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

type res_headers struct { // use MAP...
	response_line []byte
	server        []byte
	date          []byte
	cont_type     []byte
	cont_length   []byte
}

func (h *res_headers) join() []byte {
	r := []byte{}
	r = append(r, h.response_line...)
	r = append(r, h.server...)
	r = append(r, h.date...)
	r = append(r, h.cont_type...)
	r = append(r, h.cont_length...)
	return r
}

func NewResponseHeader() *res_headers {
	y, m, d := time.Now().Date()

	// default header
	header := &res_headers{
		response_line: []byte("HTTP/1.1 200 OK" + "\r\n"),
		server:        []byte("Server: " + "Vricap" + "\r\n"),
		date:          []byte("Date: " + fmt.Sprintf("%v %s %v", y, m, d) + "\r\n"),
		cont_type:     []byte("Content-Type: " + "text/html" + "\r\n"),
		cont_length:   []byte("Content-Length: " + "69" + "\r\n"),
	}
	return header
}

func (h *res_headers) responseLine(code int) {
	StatusCodes := map[int]string{
		// 1xx: Informational
		100: "Continue",
		101: "Switching Protocols",
		102: "Processing",
		103: "Early Hints",

		// 2xx: Success
		200: "OK",
		201: "Created",
		202: "Accepted",
		203: "Non-Authoritative Information",
		204: "No Content",
		205: "Reset Content",
		206: "Partial Content",
		207: "Multi-Status",
		208: "Already Reported",
		226: "IM Used",

		// 3xx: Redirection
		300: "Multiple Choices",
		301: "Moved Permanently",
		302: "Found",
		303: "See Other",
		304: "Not Modified",
		305: "Use Proxy",
		307: "Temporary Redirect",
		308: "Permanent Redirect",

		// 4xx: Client Error
		400: "Bad Request",
		401: "Unauthorized",
		402: "Payment Required",
		403: "Forbidden",
		404: "Not Found",
		405: "Method Not Allowed",
		406: "Not Acceptable",
		407: "Proxy Authentication Required",
		408: "Request Timeout",
		409: "Conflict",
		410: "Gone",
		411: "Length Required",
		412: "Precondition Failed",
		413: "Payload Too Large",
		414: "URI Too Long",
		415: "Unsupported Media Type",
		416: "Range Not Satisfiable",
		417: "Expectation Failed",
		418: "I'm a teapot",
		421: "Misdirected Request",
		422: "Unprocessable Entity",
		423: "Locked",
		424: "Failed Dependency",
		425: "Too Early",
		426: "Upgrade Required",
		428: "Precondition Required",
		429: "Too Many Requests",
		431: "Request Header Fields Too Large",
		451: "Unavailable For Legal Reasons",

		// 5xx: Server Error
		500: "Internal Server Error",
		501: "Not Implemented",
		502: "Bad Gateway",
		503: "Service Unavailable",
		504: "Gateway Timeout",
		505: "HTTP Version Not Supported",
		506: "Variant Also Negotiates",
		507: "Insufficient Storage",
		508: "Loop Detected",
		510: "Not Extended",
		511: "Network Authentication Required",
	}
	h.response_line = []byte("HTTP/1.1 " + strconv.Itoa(code) + " " + StatusCodes[code] + "\r\n")
}

// func (h *res_headers)

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

	// send back the data
	header := NewResponseHeader()
	header.responseLine(404)
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
