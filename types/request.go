package types

type Method string

const (
	GET    Method = "GET"
	POST   Method = "POST"
	PUT    Method = "PUT"
	PATH   Method = "PATCH"
	DELETE Method = "DELETE"
)

type Request struct {
	Raw      string
	Url      string
	Headers  string
	Body     string
	Name     string
	FileName string
	Verb     Method
}

type RequestLog struct {
	Request  Request
	Response string
}

type Response struct {
	Status string
	Header string
	Body   string
}
