package tls

// * 要使用的tls版本
type TlsVersion int

const DefaultTls TlsVersion = 1
const (
	Chrome112 = DefaultTls + iota
	Chrome111
	Chrome110
	Chrome109
	Chrome108
	Chrome107
	Chrome106
	Chrome105
	Ios16
	Ios15
)
