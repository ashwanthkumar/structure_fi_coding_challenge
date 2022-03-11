package binance

import (
	"io/ioutil"
	"net/http"
)

func http_Get(url string) (string, error) {
	httpClient := http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	res, getErr := httpClient.Do(req)
	if getErr != nil {
		return "", getErr
	}

	return readRespBody(res)
}

func readRespBody(res *http.Response) (string, error) {
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return "", readErr
	}

	return string(body), nil
}
