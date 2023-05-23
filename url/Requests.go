package url

import (
	"github.com/bogdanfinn/fhttp/cookiejar"
	"time"
)

func NewRequest() *Request {
	return &Request{
		AllowRedirects: true,
		Verify:         true,
	}
}

type Request struct {
	Params         *Params
	Headers        *Header
	Cookies        *cookiejar.Jar
	Data           *Values
	Json           map[string]interface{}
	Body           string
	Auth           []string
	Timeout        time.Duration
	AllowRedirects bool
	Proxies        string
	Verify         bool
	ForceHTTP1     bool
	TlsProfile     int
	Method         string
	Url            string
}
