package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/wmm1996528/requests/url"
)

// Response结构体
type Response struct {
	Url        string
	Headers    url.Header
	Cookies    map[string]string
	Text       string
	Content    []byte
	StatusCode int
	History    []*Response
	Request    *url.Request
}

// 使用自带库JSON解析
func (res *Response) Json() (map[string]interface{}, error) {
	js := make(map[string]interface{})
	err := json.Unmarshal(res.Content, &js)
	return js, err
}

// 使用go-simplejson解析
func (res *Response) SimpleJson() (*simplejson.Json, error) {
	return simplejson.NewJson(res.Content)
}

// 状态码是否错误
func (res *Response) RaiseForStatus() error {
	var err error
	if res.StatusCode >= 400 && res.StatusCode < 500 {
		err = errors.New(fmt.Sprintf("%d Client Error", res.StatusCode))
	} else if res.StatusCode >= 500 && res.StatusCode < 600 {
		err = errors.New(fmt.Sprintf("%d Server Error", res.StatusCode))
	}
	return err
}
