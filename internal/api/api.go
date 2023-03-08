package api

import (
	"fmt"
	"io"
	"net/http"
	// "net/http/httputil"
	"os"
)

var (
	client  = &http.Client{}
	baseURL = "http://127.0.0.1:8080"
	headers = http.Header{
		"Content-Type": {"text/plain"},
	}
)

func Fetch(method, endpoint string, body io.Reader) {
	request, err := http.NewRequest(method, baseURL+"/"+endpoint, body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	request.Header = headers

	// requestDump, err := httputil.DumpRequest(request, true)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// fmt.Println(string(requestDump))

	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	response.Body.Close()

	fmt.Println(string(resBody))
}
