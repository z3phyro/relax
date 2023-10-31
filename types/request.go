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
	Raw     string
	Url     string
	Headers string
	Body    string
	Name    string
	Verb    Method
}
