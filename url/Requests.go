package url

import (
	"github.com/wmm1996528/requests/tls"
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
	Cookies        map[string]string
	Data           *Values
	Json           map[string]interface{}
	Body           string
	Auth           []string
	Timeout        time.Duration
	AllowRedirects bool
	Verify         bool
	ForceHTTP1     bool
	TlsProfile     tls.TlsVersion
	Method         string
	Url            string
	Proxy          string
}
