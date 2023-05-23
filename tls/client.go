package tls

import (
	tls_client "github.com/bogdanfinn/tls-client"
	"log"
)

func NewClient(tlsProfile int) tls_client.HttpClient {
	var tlsConfig tls_client.ClientProfile
	switch tlsProfile {
	case Chrome112:
		tlsConfig = tls_client.Chrome_112
	case Chrome111:
		tlsConfig = tls_client.Chrome_111
	case Chrome110:
		tlsConfig = tls_client.Chrome_110
	case Chrome109:
		tlsConfig = tls_client.Chrome_109
	case Chrome108:
		tlsConfig = tls_client.Chrome_108
	case Chrome107:
		tlsConfig = tls_client.Chrome_107
	case Chrome106:
		tlsConfig = tls_client.Chrome_106
	case Chrome105:
		tlsConfig = tls_client.Chrome_105
	case Ios16:
		tlsConfig = tls_client.Safari_IOS_16_0
	case Ios15:
		tlsConfig = tls_client.Safari_IOS_15_5
	default:
		tlsConfig = tls_client.Chrome_112
	}
	jar := tls_client.NewCookieJar()
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(15),
		tls_client.WithClientProfile(tlsConfig),
		tls_client.WithCookieJar(jar), // create cookieJar instance and pass it as argument
	}
	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		log.Println(err)
	}
	return client
}
