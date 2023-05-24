package requests

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/binary"
	"encoding/json"
	"github.com/andybalholm/brotli"
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/wmm1996528/requests/models"
	"github.com/wmm1996528/requests/tls"
	"github.com/wmm1996528/requests/url"
	"io"
	url2 "net/url"
	"strings"
)

func NewSession(tlsVersion tls.TlsVersion) *Session {
	client := tls.NewClient(tlsVersion)
	return &Session{Client: client}
}

type Session struct {
	Headers        map[string]string
	Cookies        map[string]string
	Auth           []string
	Proxy          string
	Verify         bool
	tlsVersion     int
	Client         tls_client.HttpClient
	AllowRedirects bool
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
	// * 处理cookie
	if request.Cookies != nil {
		var cks []*http.Cookie
		for k, v := range request.Cookies {
			cks = append(cks, &http.Cookie{
				Name:  k,
				Value: v,
			})
		}
		uri, _ := url2.Parse(request.Url)
		s.Client.SetCookies(uri, cks)
	}
	// * 处理代理
	if request.Proxy != "" {
		s.Client.SetProxy(request.Proxy)
	} else {
		if s.Proxy != "" {
			s.Client.SetProxy(s.Proxy)
		}
	}
	// * 是否自动跳转
	if request.AllowRedirects {
		s.Client.SetFollowRedirect(request.AllowRedirects)
	} else {
		if s.AllowRedirects {
			s.Client.SetFollowRedirect(s.AllowRedirects)
		} else {
			s.Client.SetFollowRedirect(s.AllowRedirects)
		}
	}

	// * 组合header
	var headers map[string]string
	if request.Headers != nil {
		headers = request.Headers.GetAll()
	} else {
		if len(s.Headers) != 0 {
			headers = s.Headers
		} else {
			headers = url.DefaultHeaders
		}
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
	if request.Body != "" {
		return []byte(request.Body), nil
	}
	return []byte{}, nil
}

// PreResponse 处理请求返回
func (s *Session) PreResponse(request *url.Request, do *http.Response) (*models.Response, error) {
	rb, err := io.ReadAll(do.Body)
	if err != nil {
		return nil, err
	}
	encoding := do.Header.Get("Content-Encoding")
	DecompressBody(&rb, encoding)

	redirectURL, err := do.Location()
	redirectUrl := ""
	if err == nil {
		// 没有错误时从跳转的location 获取url
		redirectUrl = redirectURL.String()
	} else {
		redirectUrl = request.Url
	}
	resp := &models.Response{
		Url:        redirectUrl,
		Headers:    url.Header(do.Header),
		Cookies:    s.getAllCookie(request),
		Text:       string(rb),
		Content:    rb,
		StatusCode: do.StatusCode,
		Request:    request,
	}
	s.Cookies = resp.Cookies

	return resp, nil
}

func (s *Session) getAllCookie(request *url.Request) map[string]string {
	res := make(map[string]string)
	uri, _ := url2.Parse(request.Url)
	cks := s.Client.GetCookies(uri)
	for _, ck := range cks {
		res[ck.Name] = ck.Value
	}

	return res
}

// 解码Body数据
func DecompressBody(content *[]byte, encoding string) {
	if encoding != "" {
		if strings.ToLower(encoding) == "gzip" {
			decodeGZip(content)
		} else if strings.ToLower(encoding) == "deflate" {
			decodeDeflate(content)
		} else if strings.ToLower(encoding) == "br" {
			decodeBrotli(content)
		}
	}
}

// 解码GZip编码
func decodeGZip(content *[]byte) error {
	if content == nil {
		return nil
	}
	b := new(bytes.Buffer)
	binary.Write(b, binary.LittleEndian, content)
	r, err := gzip.NewReader(b)
	if err != nil {
		return err
	}
	defer r.Close()
	*content, err = io.ReadAll(r)
	if err != nil {
		return err
	}
	return nil
}

// 解码deflate编码
func decodeDeflate(content *[]byte) error {
	var err error
	if content == nil {
		return err
	}
	r := flate.NewReader(bytes.NewReader(*content))
	defer r.Close()
	*content, err = io.ReadAll(r)
	if err != nil {
		return err
	}
	return nil
}

// 解码br编码
func decodeBrotli(content *[]byte) error {
	var err error
	if content == nil {
		return err
	}
	r := brotli.NewReader(bytes.NewReader(*content))
	*content, err = io.ReadAll(r)
	if err != nil {
		return err
	}
	return nil
}
