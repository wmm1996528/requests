package url

import (
	"errors"
	"net/textproto"
	"strings"
)

const HeaderOrderKey = "Header-Order:"

// PHeaderOrderKey is a magic Key for setting http2 pseudo header order.
// If the header is nil it will use regular GoLang header order.
// Valid fields are :authority, :method, :path, :scheme
const PHeaderOrderKey = "PHeader-Order:"

// * 请求header 部分
type Header Common

func (h Header) Add(key, value string) {
	textproto.MIMEHeader(h).Add(key, value)
}

// Set sets the header entries associated with Key to the
// single element value. It replaces any existing Values
// associated with Key. The Key is case insensitive; it is
// canonicalized by textproto.CanonicalMIMEHeaderKey.
// To use non-canonical keys, assign to the map directly.
func (h Header) Set(key, value string) {
	textproto.MIMEHeader(h).Set(key, value)
}

// Get gets the first value associated with the given Key. If
// there are no Values associated with the Key, Get returns "".
// It is case insensitive; textproto.CanonicalMIMEHeaderKey is
// used to canonicalize the provided Key. To use non-canonical keys,
// access the map directly.
func (h Header) Get(key string) string {
	return textproto.MIMEHeader(h).Get(key)
}
func (h Header) GetAll() map[string]string {
	res := make(map[string]string)
	for k, _ := range h {
		res[k] = h.Get(k)
	}
	return res
}

// Values returns all Values associated with the given Key.
// It is case insensitive; textproto.CanonicalMIMEHeaderKey is
// used to canonicalize the provided Key. To use non-canonical
// keys, access the map directly.
// The returned slice is not a copy.
func (h Header) Values(key string) []string {
	return textproto.MIMEHeader(h).Values(key)
}

// get is like Get, but Key must already be in CanonicalHeaderKey form.
func (h Header) get(key string) string {
	if v := h[key]; len(v) > 0 {
		return v[0]
	}
	return ""
}

// has reports whether h has the provided Key defined, even if it's
// set to 0-length slice.
func (h Header) has(key string) bool {
	_, ok := h[key]
	return ok
}

// Del deletes the Values associated with Key.
// The Key is case insensitive; it is canonicalized by
// CanonicalHeaderKey.
func (h Header) Del(key string) {
	textproto.MIMEHeader(h).Del(key)
}

// Clone returns a copy of h or nil if h is nil.
func (h Header) Clone() Header {
	if h == nil {
		return nil
	}

	// Find total number of Values.
	nv := 0
	for _, vv := range h {
		nv += len(vv)
	}
	sv := make([]string, nv) // shared backing array for headers' Values
	h2 := make(Header, len(h))
	for k, vv := range h {
		n := copy(sv, vv)
		h2[k] = sv[:n:n]
		sv = sv[n:]
	}
	return h2
}

// 初始化Headers结构体
func NewHeaders() *Header {
	headers := &Header{}
	(*headers)[PHeaderOrderKey] = (*headers)[PHeaderOrderKey]
	return headers
}

// 解析Headers字符串为结构体
func ParseHeaders(headers string) *Header {
	h := Header{}
	headerOrder := []string{}
	lines := strings.Split(headers, "\n")
	for _, header := range lines {
		header = strings.TrimSpace(header)
		if header == "" || strings.Index(header, ":") == 0 || strings.Index(header, "/") == 0 || strings.Index(header, "#") == 0 {
			continue
		}
		keyValue := strings.SplitN(header, ":", 2)
		if len(keyValue) != 2 {
			panic(errors.New("该字符串不符合http头部标准！"))
		}
		key := keyValue[0]
		value := keyValue[1]
		h.Set(key, value)
		headerOrder = append(headerOrder, strings.ToLower(key))
	}
	h[HeaderOrderKey] = headerOrder
	return &h
}
