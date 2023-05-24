package requests

import (
	"bytes"
	"encoding/json"
	http "github.com/bogdanfinn/fhttp"
	"github.com/bogdanfinn/fhttp/cookiejar"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/wmm1996528/requests/models"
	"github.com/wmm1996528/requests/tls"
	"github.com/wmm1996528/requests/url"
	"io"
)

func NewSession(tlsVersion tls.TlsVersion) *Session {
	client := tls.NewClient(tlsVersion)
	return &Session{Client: client}
}

type Session struct {
	Params       *url.Params
	Headers      *url.Header
	Cookies      *cookiejar.Jar
	Auth         []string
	Proxies      string
	Verify       bool
	Cert         []string
	Ja3          string
	MaxRedirects int
	request      *url.Request
	tlsVersion   int
	Client       tls_client.HttpClient
}

// 预请求处理

// http请求方式基础函数
func (s *Session) Request(method, rawurl string, request *url.Request) (*models.Response, error) {
	if request == nil {
		request = url.NewRequest()
	}
	request.Url = rawurl
	resp, err := s.Do(method, request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// get请求方式
func (s *Session) Get(rawurl string, req *url.Request) (*models.Response, error) {
	return s.Request(http.MethodGet, rawurl, req)
}

// post请求方式
func (s *Session) Post(rawurl string, req *url.Request) (*models.Response, error) {
	return s.Request(http.MethodPost, rawurl, req)
}

// options请求方式
func (s *Session) Options(rawurl string, req *url.Request) (*models.Response, error) {
	return s.Request(http.MethodOptions, rawurl, req)
}

// head请求方式
func (s *Session) Head(rawurl string, req *url.Request) (*models.Response, error) {
	return s.Request(http.MethodHead, rawurl, req)
}

// put请求方式
func (s *Session) Put(rawurl string, req *url.Request) (*models.Response, error) {
	return s.Request(http.MethodPut, rawurl, req)
}

// patch请求方式
func (s *Session) Patch(rawurl string, req *url.Request) (*models.Response, error) {
	return s.Request(http.MethodPatch, rawurl, req)
}

// delete请求方式
func (s *Session) Delete(rawurl string, req *url.Request) (*models.Response, error) {
	return s.Request(http.MethodDelete, rawurl, req)
}

// connect请求方式
func (s *Session) Connect(rawurl string, req *url.Request) (*models.Response, error) {
	return s.Request(http.MethodConnect, rawurl, req)
}

func (s *Session) Do(method string, request *url.Request) (*models.Response, error) {
	if s.Client == nil {
		// 初始化个新的
		s.Client = tls.NewClient(request.TlsProfile)

	}

	request.Method = method
	preq, err := s.PreRequest(request)
	if err != nil {
		return nil, err
	}
	do, err := s.Client.Do(preq)
	if err != nil {
		return nil, err
	}
	defer do.Body.Close()
	return s.PreResponse(request, do)
}

func (s *Session) PreRequest(request *url.Request) (*http.Request, error) {
	datas, err := s.PreData(request)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(request.Method, request.Url, bytes.NewBuffer(datas))
	if err != nil {
		return nil, err
	}

	var headers map[string]string
	if request.Headers != nil {
		headers = request.Headers.GetAll()
	} else {
		headers = url.DefaultHeaders
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return req, nil
}

// PreData 解析post参数
func (s *Session) PreData(request *url.Request) ([]byte, error) {
	if request.Json != nil {
		return json.Marshal(request.Json)
	}
	if request.Data != nil {
		return []byte(request.Data.Encode()), nil
	}
	return []byte{}, nil
}

// PreResponse 处理请求返回
func (s *Session) PreResponse(request *url.Request, do *http.Response) (*models.Response, error) {
	reader := http.DecompressBody(do)
	rb, _ := io.ReadAll(reader)
	redirectURL, err := do.Location()
	redirectUrl := ""
	if err == nil {
		// 没有错误时从跳转的location 获取url
		redirectUrl = redirectURL.String()
	}
	resp := &models.Response{
		Url:        redirectUrl,
		Headers:    url.Header(do.Header),
		Cookies:    do.Cookies(),
		Text:       string(rb),
		Content:    rb,
		Body:       reader,
		StatusCode: do.StatusCode,
		Request:    request,
	}

	return resp, nil
}
