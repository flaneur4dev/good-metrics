package api

import (
	"fmt"
	"io"
	"net/http"
	// "net/http/httputil"

	"github.com/flaneur4dev/good-metrics/internal/lib/utils"
)

var (
	client     = &http.Client{}
	baseURL, _ = utils.EnvVar("ADDRESS", "http://127.0.0.1:8080").(string)
	headers    = http.Header{
		"Content-Type": {"application/json"},
		// "Content-Type": {"text/plain"},
	}
)

func Fetch(method, endpoint string, body io.Reader) {
	request, err := http.NewRequest(method, baseURL+"/"+endpoint, body)
	if err != nil {
		fmt.Println(err)
		return
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
