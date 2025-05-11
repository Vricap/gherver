package http

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

const SIZE int = 1024 // size of the data the server will read. 1kb.

type Response struct {
	Headers *resHeaders
	Body    *resBody
}

func SendResponse(conn net.Conn, h *resHeaders, b *resBody) (error, any) {
	h.contLength = strconv.Itoa(len(b.body))
	res := fmt.Sprintf(`%s
Server: %s
Date: %s
Content-Type: %s
Content-Length: %s

%s`, h.responLine, h.server, h.date, h.contType, h.contLength, b.body)
	_, err := conn.Write([]byte(res))
	if err != nil {
		return err, nil
	}
	return nil, res
}

type resHeaders struct { // use MAP...
	responLine string
	server     string
	date       string
	contType   string
	contLength string
}

func NewResponseHeader() *resHeaders {
	y, m, d := time.Now().Date()

	// default header
	header := &resHeaders{
		responLine: "HTTP/1.1 200 OK",
		server:     "Vricap",
		date:       fmt.Sprintf("%v %s %v", y, m, d),
		contType:   "text/html",
		contLength: "69",
	}
	return header
}

func (h *resHeaders) setStatusCode(code int) {
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
	h.responLine = "HTTP/1.1 " + strconv.Itoa(code) + " " + StatusCodes[code]
}

type resBody struct {
	body []byte
}

func NewResponseBody(b []byte) *resBody {
	body := &resBody{body: b}
	return body
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

	// make a new header
	header := NewResponseHeader()
	// set your own headers
	header.setStatusCode(200)
	header.contType = "text/html"
	// set response body
	// body := NewResponseBody(buf)
	body := NewResponseBody([]byte(`
<html>
	<head>
		<title>Test</title>
	</head>
	<body>
		<h1>Hello World!</h1>
		<p>Foo Bar</p>
	</body>
</html>`))

	// send back the response
	err, data := SendResponse(conn, header, body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// print the incoming data
	fmt.Printf("Received: %s", data)
}
