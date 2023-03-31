package api

import (
	"fmt"
	"io"
	"net/http"
)

var (
	client  = &http.Client{}
	headers = http.Header{
		"Content-Type": {"application/json"},
		// "Content-Type": {"text/plain"},
	}
)

func Fetch(method, endpoint string, body io.Reader) {
	request, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		fmt.Println(err)
		return
	}
	request.Header = headers

	response, err := client.Do(request)
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
