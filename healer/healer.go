package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log/syslog"
	"net/http"
	"os"
	"sync"
)

var (
	log     *syslog.Writer
	mut     sync.Mutex
	healers = make(map[string]*healer)
)

type healer struct {
	url string
}

func setHealers(h map[string]*healer) {
	mut.Lock()
	healers = h
	mut.Unlock()
}

func getHealers() map[string]*healer {
	mut.Lock()
	defer mut.Unlock()
	return healers
}

func (h *healer) heal() error {
	log.Info(fmt.Sprintf("healing tsuru healer with endpoint %s...", h.url))
	r, err := request("GET", h.url, nil)
	if err == nil {
		r.Body.Close()
	}
	return err
}

// healersFromResource returns healers registered in tsuru.
func healersFromResource(endpoint string) (map[string]*healer, error) {
	url := fmt.Sprintf("%s/healers", endpoint)
	response, err := request("GET", url, nil)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return nil, err
	}
	var h map[string]*healer
	data := map[string]string{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	h = make(map[string]*healer, len(data))
	for name, url := range data {
		h[name] = &healer{url: fmt.Sprintf("%s%s", endpoint, url)}
	}
	return h, nil
}

func request(method, url string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if token := os.Getenv("TSURU_TOKEN"); token != "" {
		request.Header.Add("Authorization", token)
		request.Header.Add("Token-Owner", os.Getenv("TSURU_TOKEN_OWNER"))
	}
	resp, err := (&http.Client{}).Do(request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
