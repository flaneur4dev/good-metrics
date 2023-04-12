package client

import (
	"fmt"
	"io"
	"net/http"
)

type MetricsClient struct {
	client  *http.Client
	scheme  string
	addr    string
	baseURL string
	headers http.Header
}

func New(ad string, secure bool) *MetricsClient {
	mc := &MetricsClient{
		client: &http.Client{},
		scheme: "http://",
		addr:   ad,
		headers: http.Header{
			"Content-Type": {"application/json"},
		},
	}

	if secure {
		mc.scheme = "https://"
	}
	mc.baseURL = mc.scheme + mc.addr

	return mc
}

func (mc *MetricsClient) Fetch(method, endpoint string, body io.Reader) {
	request, err := http.NewRequest(method, mc.baseURL+endpoint, body)
	if err != nil {
		fmt.Println(err)
		return
	}
	request.Header = mc.headers

	response, err := mc.client.Do(request)
	if err != nil {
		fmt.Println(err)
		return
	}

	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	response.Body.Close()

	fmt.Println("response:", string(resBody))
}
