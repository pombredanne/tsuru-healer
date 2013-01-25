package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log/syslog"
	"net/http"
	"sync"
)

var (
	log     *syslog.Writer
	mut     sync.Mutex
	healers = make(map[string]healer)
)

type healer interface {
	heal() error
}

func register(name string, h healer) {
	mut.Lock()
	defer mut.Unlock()
	log.Info(fmt.Sprintf("registering %s healer", name))
	healers[name] = h
}

func getHealers() map[string]healer {
	mut.Lock()
	defer mut.Unlock()
	return healers
}

type tsuruHealer struct {
	url string
}

func (h *tsuruHealer) heal() error {
	log.Info(fmt.Sprintf("healing tsuru healer with endpoint %s...", h.url))
	_, err := request("GET", h.url, "", nil)
	return err
}

type instanceHealer struct {
	endpoint string
	seeker   seeker
	token    string
}

func newInstanceHealer(email, password, endpoint string) *instanceHealer {
	token, err := getToken(email, password, endpoint)
	if err != nil {
		panic(err)
	}
	return &instanceHealer{
		seeker:   newAWSSeeker(),
		endpoint: endpoint,
		token:    token,
	}
}

// Heal iterates through down instances, terminate then
// and spawn new ones to replace the terminated.
func (h *instanceHealer) heal() error {
	log.Info("Starting healing process... this can take a while.")
	instances, err := h.seeker.seekUnhealthyInstances()
	if err != nil {
		return err
	}
	for _, instance := range instances {
		if err := h.terminate(instance.loadBalancer, instance.instanceId); err != nil {
			// should really stop here?
			log.Err("Got error while terminating instance: " + err.Error())
			return err
		}
		if err := h.spawn(instance.loadBalancer); err != nil {
			log.Err("Got error while spawining instance: " + err.Error())
		}
	}
	return nil
}

// Calls tsuru add-unit endpoint
func (h *instanceHealer) spawn(lb string) error {
	url := fmt.Sprintf("%s/apps/%s/units", h.endpoint, lb)
	body := bytes.NewBufferString("1")
	resp, err := request("PUT", url, h.token, body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error spawning unit: %s", resp.Status)
	}
	return nil
}

// Calls tsuru remove-unit endpoint
func (h *instanceHealer) terminate(lb, id string) error {
	url := fmt.Sprintf("%s/apps/%s/unit", h.endpoint, lb)
	body := bytes.NewBufferString(id)
	resp, err := request("DELETE", url, h.token, body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error terminating unit: %s", resp.Status)
	}
	return nil
}

// healersFromResource returns healers registered in tsuru.
func healersFromResource(endpoint string) (map[string]tsuruHealer, error) {
	url := fmt.Sprintf("%s/healers", endpoint)
	response, err := request("GET", url, "", nil)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	h := map[string]tsuruHealer{}
	data := map[string]string{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	for name, url := range data {
		h[name] = tsuruHealer{url: fmt.Sprintf("%s%s", endpoint, url)}
	}
	return h, nil
}

func getToken(email, password, endpoint string) (string, error) {
	url := fmt.Sprintf("%s/users/%s/tokens", endpoint, email)
	b := fmt.Sprintf(`{"password": "%s"}`, password)
	body := bytes.NewBufferString(b)
	resp, err := http.Post(url, "application/json", body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Error obtaining token: %s", resp.Status)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var token map[string]string
	json.Unmarshal(respBody, &token)
	if _, ok := token["token"]; !ok {
		return "", errors.New("Unknown response format.")
	}
	return token["token"], nil
}

func request(method, url, token string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", token)
	resp, err := (&http.Client{}).Do(request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
