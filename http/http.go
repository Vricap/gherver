package http

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

const SIZE int = 1024 // size of the data the server will read. 1kb.

func Init() *http {
	return &http{
		Response: &Response{
			Headers: &resHeaders{},
			Body:    &resBody{},
		},
	}
}

type http struct {
	Response *Response
}

func (h *http) StartServer(ADDR, PORT string) {
	// listen for incoming connections on port 8000
	ln, err := net.Listen("tcp", ADDR+PORT)
	if err != nil {
		fmt.Println(err)
		return
	}

	// accept incoming connections and handle them
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		// handle the connections in a new goroutine
		go HandleConnection(conn, h)
	}
}

type Response struct {
	Headers *resHeaders
	Body    *resBody
}

func (rs *Response) SendResponse(conn net.Conn) (error, any) {
	rs.Headers.ContLength = strconv.Itoa(len(rs.Body.body))
	res := fmt.Sprintf(`%s
Server: %s
Date: %s
Content-Type: %s
Content-Length: %s

%s`, rs.Headers.ResponLine, rs.Headers.Server, rs.Headers.Date, rs.Headers.ContType, rs.Headers.ContLength, rs.Body.body)
	_, err := conn.Write([]byte(res))
	if err != nil {
		return err, nil
	}
	return nil, res
}

type resHeaders struct { // use MAP...
	ResponLine string
	Server     string
	Date       string
	ContType   string
	ContLength string
}

func (rh *resHeaders) NewResponseHeader() {
	y, m, d := time.Now().Date()

	// default header
	rh.ResponLine = "HTTP/1.1 200 OK"
	rh.Server = "Vricap"
	rh.Date = fmt.Sprintf("%v %s %v", y, m, d)
	rh.ContType = "text/html"
	rh.ContLength = "69"
}

func (h *resHeaders) SetStatusCode(code int) {
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
	h.ResponLine = "HTTP/1.1 " + strconv.Itoa(code) + " " + StatusCodes[code]
}

type resBody struct {
	body []byte
}

func (rb *resBody) NewResponseBody(b []byte) {
	rb.body = b
}

func HandleConnection(conn net.Conn, h *http) {
	// close the connection when we're done
	defer conn.Close()

	// read incoming data
	buf := make([]byte, SIZE)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	// send back the response
	err, data := h.Response.SendResponse(conn)
	if err != nil {
		fmt.Println(err)
		return
	}

	// print the response data
	fmt.Printf("Received: %s", data)
}
