package binance

import (
	"io/ioutil"
	"net/http"
)

func Get(url string) (string, error) {
	httpClient := http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	res, getErr := httpClient.Do(req)
	if getErr != nil {
		return "", err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return "", err
	}

	return string(body), nil
}
