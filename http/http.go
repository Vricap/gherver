package http

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const SIZE int = 1024 // size of the data the server will read. 1kb.

func Init() *Http {
	return &Http{
		Response: &Response{
			Headers: &resHeaders{},
			Body:    &resBody{},
		},
		Request: &Request{
			Headers: &reqHeaders{},
			Body:    &reqBody{},
		},
		Routes:  []*Routes{},
		Static:  []*Static{},
		Logging: true,
	}
}

type Http struct {
	Response *Response
	Request  *Request
	Routes   []*Routes
	Static   []*Static
	Logging  bool
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

func (h *Http) LoadStatic(prefix, dir string) {
	fs, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, n := range fs {
		path := dir + "/" + n.Name()
		ext := getExt(path)
		v, ok := docuContTypes[ext]
		if !ok {
			v, ok = mediaContypes[ext]
			if !ok {
				fmt.Printf("Unsupported file type: '%s'\n", ext)
				return
			}
		}
		content, err := os.ReadFile(path)
		if err != nil {
			fmt.Println(err)
			return
		}
		h.Static = append(h.Static, &Static{
			oriPath:    path,
			prefixPath: prefix + "/" + n.Name(),
			contType:   v,
			content:    content,
		})
	}
}

func (h *Http) DisableLog() {
	h.Logging = false
}

type Static struct {
	oriPath    string
	prefixPath string
	contType   string
	content    []byte
}

type Routes struct {
	Path    string
	Method  string
	Handler func(*Http)
}

type Response struct {
	Headers *resHeaders
	Body    *resBody
}

func (rs *Response) SendHtml(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	rs.Body.Body = data
}

func getExt(path string) string {
	l := len(path) - 1
	var ext []byte = []byte{}
	for i := l; i != 0; i-- {
		if path[i] == '.' {
			break
		}
		ext = append([]byte(string(path[i])), ext...)
	}
	return string(ext)
}

var docuContTypes map[string]string = map[string]string{
	// Documents
	"pdf": "application/pdf", "txt": "text/plain", "html": "text/html", "htm": "text/html", "css": "text/css", "js": "application/javascript", "json": "application/json", "xml": "application/xml", "csv": "text/csv",

	// Archives
	"zip": "application/zip", "tar": "application/x-tar", "gz": "application/gzip", "rar": "application/vnd.rar", "7z": "application/x-7z-compressed",
}

var mediaContypes map[string]string = map[string]string{
	// Images
	"png": "image/png", "jpg": "image/jpeg", "jpeg": "image/jpeg", "gif": "image/gif", "bmp": "image/bmp", "webp": "image/webp", "svg": "image/svg+xml", "ico": "image/x-icon",

	// Videos
	"mp4": "video/mp4", "webm": "video/webm", "ogg": "video/ogg", "mov": "video/quicktime", "avi": "video/x-msvideo", "mkv": "video/x-matroska",

	// Audio
	"mp3": "audio/mpeg", "wav": "audio/wav", "flac": "audio/flac", "aac": "audio/aac", "m4a": "audio/mp4",
}

func (rs *Response) SendDocument(path string) {
	ext := getExt(path)
	t, ok := docuContTypes[ext]
	if !ok {
		fmt.Printf("Unsupported file type: '%s'\n", ext)
		return
	}
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	rs.Headers.ContType = t
	rs.Body.Body = data
}

func (rs *Response) SendMedia(path string) {
	ext := getExt(path)
	t, ok := mediaContypes[ext]
	if !ok {
		fmt.Printf("Unsupported file type: '%s'\n", ext)
		return
	}
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	rs.Headers.ContType = t
	rs.Body.Body = data
}

func (rs *Response) constructResponse() []byte {
	// set some default value if user didn't set them
	if len(rs.Headers.ResponLine) == 0 && len(rs.Body.Body) == 0 {
		rs.Body.Byte([]byte(`
<html>
	<head>
		<title>Error</title>
	</head>
	<body>
		<h1>500 INTERNAL SERVER ERROR!</h1>
		<p>Server didn't set the status code and the body!</p>
	</body>
</html>`))
	}
	rs.Headers.ContLength = strconv.Itoa(len(rs.Body.Body))
	if len(rs.Headers.ResponLine) == 0 {
		rs.Headers.ResponLine = "HTTP/1.1 500 Internal Server Error"
	}
	if len(rs.Headers.Server) == 0 {
		rs.Headers.Server = "GHERVER"
	}
	if len(rs.Headers.Date) == 0 {
		y, m, d := time.Now().Date()
		rs.Headers.Date = fmt.Sprintf("%v %s %v", y, m, d)
	}
	if len(rs.Headers.ContType) == 0 {
		rs.Headers.ContType = "text/html"
	}
	if len(rs.Headers.AccessControlAllowOrigin) == 0 {
		rs.Headers.ContType = "*"
	}

	res := fmt.Sprintf(`%s
Server: %s
Date: %s
Content-Type: %s
Content-Length: %s
Access-Control-Allow-Origin: %s

%s`, rs.Headers.ResponLine, rs.Headers.Server, rs.Headers.Date, rs.Headers.ContType, rs.Headers.ContLength, rs.Headers.AccessControlAllowOrigin, rs.Body.Body)

	return []byte(res)
}

type resHeaders struct { // use MAP...
	ResponLine               string
	Server                   string
	Date                     string
	ContType                 string
	ContLength               string
	AccessControlAllowOrigin string
}

func (h *resHeaders) SetStatusCode(code int) {
	StatusCodes := map[int]string{
		// i try to suppress the total loc :)
		// 1xx: Informational
		100: "Continue", 101: "Switching Protocols", 102: "Processing", 103: "Early Hints",

		// 2xx: Success
		200: "OK", 201: "Created", 202: "Accepted", 203: "Non-Authoritative Information", 204: "No Content", 205: "Reset Content", 206: "Partial Content", 207: "Multi-Status", 208: "Already Reported", 226: "IM Used",

		// 3xx: Redirection
		300: "Multiple Choices", 301: "Moved Permanently", 302: "Found", 303: "See Other", 304: "Not Modified", 305: "Use Proxy", 307: "Temporary Redirect", 308: "Permanent Redirect",

		// 4xx: Client Error
		400: "Bad Request", 401: "Unauthorized", 402: "Payment Required", 403: "Forbidden", 404: "Not Found", 405: "Method Not Allowed", 406: "Not Acceptable", 407: "Proxy Authentication Required", 408: "Request Timeout", 409: "Conflict", 410: "Gone", 411: "Length Required", 412: "Precondition Failed", 413: "Payload Too Large", 414: "URI Too Long", 415: "Unsupported Media Type", 416: "Range Not Satisfiable", 417: "Expectation Failed", 418: "I'm a teapot", 421: "Misdirected Request", 422: "Unprocessable Entity", 423: "Locked", 424: "Failed Dependency", 425: "Too Early", 426: "Upgrade Required", 428: "Precondition Required", 429: "Too Many Requests", 431: "Request Header Fields Too Large", 451: "Unavailable For Legal Reasons",

		// 5xx: Server Error
		500: "Internal Server Error", 501: "Not Implemented", 502: "Bad Gateway", 503: "Service Unavailable", 504: "Gateway Timeout", 505: "HTTP Version Not Supported", 506: "Variant Also Negotiates", 507: "Insufficient Storage", 508: "Loop Detected", 510: "Not Extended", 511: "Network Authentication Required",
	}
	h.ResponLine = "HTTP/1.1 " + strconv.Itoa(code) + " " + StatusCodes[code]
}

type resBody struct {
	Body []byte
}

func (rb *resBody) Byte(b []byte) {
	rb.Body = b
}

type Request struct {
	Headers *reqHeaders
	Body    *reqBody
}

var allMethod = map[string]string{
	"GET":    "GET",
	"POST":   "POST",
	"DELETE": "DELETE",
}

func (r *Request) parseRequest(buf []byte) *err {
	payload := strings.Split(string(buf), "\n")
	s := strings.Split(string(payload[0]), " ")
	m, ok := allMethod[s[0]]
	if !ok {
		return &err{
			code:     501,
			contType: "text/plain",
			body:     fmt.Sprintf("Method 501 is not supported!"),
		}
	}
	r.Headers.Method = m
	r.Headers.Path = s[1]
	r.Headers.HttpVer = s[2]
	r.Headers.Payloads = map[string]string{}

	fmt.Println(payload[len(payload)-2])
	for i, v := range payload {
		if i == 0 {
			continue
		}
		if v == "\r" {
			r.Body.Body = []byte(payload[i+1])
			break
		}
		s := strings.SplitN(v, ":", 2)
		r.Headers.Payloads[s[0]] = strings.TrimSpace(s[1])
	}
	return nil
}

type reqHeaders struct {
	Method   string
	Path     string
	HttpVer  string
	Payloads map[string]string
}

type reqBody struct {
	Body []byte
}

type err struct {
	code     int
	contType string
	body     string
}

func (e *err) contructErrResponse(h *Http) {
	h.Response.Headers.SetStatusCode(e.code)
	h.Response.Headers.ContType = e.contType
	h.Response.Body.Byte([]byte(e.body))
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

	if h.Logging {
		fmt.Println("============================================================")
		fmt.Println("Here is the request from browser:")
		fmt.Println(string(buf))
		fmt.Println("============================================================\r\n\r\n")
	}

	// parse the request
	e := h.Request.parseRequest(buf)
	if e != nil {
		e.contructErrResponse(h)
		sendResponse(h, conn)
	}

	// handle static
	if len(h.Static) != 0 {
		for _, s := range h.Static {
			route := &Routes{
				Path:   s.prefixPath,
				Method: "GET",
				Handler: func(h *Http) {
					h.Response.Headers.SetStatusCode(200)
					h.Response.Headers.ContType = s.contType
					h.Response.Body.Byte(s.content)
				},
			}
			h.Routes = append(h.Routes, route)
		}
	}

	// handle request
	for i, v := range h.Routes {
		if v.Path == h.Request.Headers.Path && v.Method == h.Request.Headers.Method {
			v.Handler(h)
			break
		}
		// if didn't find any route with such method, return 404
		if i == len(h.Routes)-1 {
			h.Response.Headers.SetStatusCode(404)
			h.Response.Headers.ContType = "text/html"
			h.Response.Body.Byte([]byte(`
<html>
	<head>
		<title>Error</title>
	</head>
	<body>
		<h1>404 NOT FOUND!</h1>
		<p>Request routes or http method didn't exist!</p>
	</body>
</html>`))
		}
	}

	sendResponse(h, conn)
}

func sendResponse(h *Http, conn net.Conn) {
	// contruct the response
	data := h.Response.constructResponse()

	// send back the response
	_, err := conn.Write(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	if h.Logging {
		fmt.Println("============================================================")
		fmt.Println("Here is the response from server:")
		fmt.Println(string(data))
		fmt.Println("============================================================\r\n\r\n")
	}
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
