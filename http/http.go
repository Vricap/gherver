package http

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

const SIZE int = 1024 // size of the data the server will read. 1kb.

var ctx *Http = &Http{
	Response: &Response{
		Headers: &resHeaders{},
		Body:    &resBody{},
	},
	Request: &Request{
		Headers: &reqHeaders{},
		Body:    &reqBody{},
	},
}

var hand *Handler = &Handler{
	GET: func(p string, f func(h *Http)) {
	},
}

func Init() *Http {
	return ctx
}

type Http struct {
	Response *Response
	Request  *Request
	Handle   *Handler
}

func (h *Http) StartServer(ADDR, PORT string) {
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

type Handler struct {
	GET func(string, func(*Http))
}

// TODO: add errors field
type Response struct {
	Headers *resHeaders
	Body    *resBody
}

func (rs *Response) SendResponse(conn net.Conn) (error, any) {
	rs.Headers.ContLength = strconv.Itoa(len(rs.Body.Body))
	res := fmt.Sprintf(`%s
Server: %s
Date: %s
Content-Type: %s
Content-Length: %s

%s`, rs.Headers.ResponLine, rs.Headers.Server, rs.Headers.Date, rs.Headers.ContType, rs.Headers.ContLength, rs.Body.Body)
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
	rh.ResponLine = "HTTP / 1.1 200 OK"
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
	Body []byte
}

func (rb *resBody) NewResponseBody(b []byte) {
	rb.Body = b
}

// TODO: add errors field
type Request struct {
	Headers *reqHeaders
	Body    *reqBody
}

type reqHeaders struct {
	Method  string
	Path    string
	HttpVer string
}

// var allMethod = map[string]string{
// 	"GET":    "GET",
// 	"POST":   "POST",
// 	"DELETE": "DELETE",
// }

func (rh *reqHeaders) parseRequestHeaders(buf []byte) {
	s := strings.Split(string(buf), " ")
	rh.Method = s[0] // TODO: add supported method checking maybe???
	rh.Path = s[1]
	rh.HttpVer = s[2]
}

type reqBody struct {
	Body []byte
}

func HandleConnection(conn net.Conn, h *Http) {
	// close the connection when we're done
	defer conn.Close()

	// read incoming data
	buf := make([]byte, SIZE)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	// parse the request
	h.Request.Headers.parseRequestHeaders(buf)

	// send back the response
	err, data := h.Response.SendResponse(conn)
	if err != nil {
		fmt.Println(err)
		return
	}

	// print the response data
	fmt.Printf("Received: %s", data)
}

/*
EXAMPLE OF HTTP REQUEST FROM THE BROWSER

GET / HTTP/1.1
Host: localhost:8000
User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:136.0) Gecko/20100101 Firefox/136.0
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,* /*;q=0.8
Accept-Language: en-US,en;q=0.5
Accept-Encoding: gzip, deflate, br, zstd
Connection: keep-alive
Upgrade-Insecure-Requests: 1
Sec-Fetch-Dest: document
Sec-Fetch-Mode: navigate
Sec-Fetch-Site: none
Sec-Fetch-User: ?1
Priority: u=0, i

--some request body i any--
*/
