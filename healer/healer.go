package healer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log/syslog"
	"net/http"
)

type Healer interface {
	Heal() error
	Spawn(lb string) error
	Terminate(lb, id string) error
}

type TsuruHealer struct {
	Endpoint string
	seeker   Seeker
	token    string
}

var log *syslog.Writer

func init() {
	var err error
	log, err = syslog.New(syslog.LOG_INFO, "tsuru-healer")
	if err != nil {
		panic(err)
	}
}

// Heal iterates through down instances, terminate then
// and spawn new ones to replace the terminated.
func (h *TsuruHealer) Heal() error {
	log.Info("Starting healing process... this can take a while.")
	instances, err := h.seeker.SeekUnhealthyInstances()
	if err != nil {
		return err
	}
	for _, instance := range instances {
		if err := h.Terminate(instance.LoadBalancer, instance.InstanceId); err != nil {
			// should really stop here?
			log.Err("Got error while terminating instance: " + err.Error())
			return err
		}
		if err := h.Spawn(instance.LoadBalancer); err != nil {
			log.Err("Got error while spawining instance: " + err.Error())
		}
	}
	return nil
}

func getToken(email, password, endpoint string) (string, error) {
	url := fmt.Sprintf("%s/users/%s/tokens", endpoint, email)
	b := fmt.Sprintf(`{"password": %s}`, password)
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

// Calls tsuru add-unit endpoint
func (h *TsuruHealer) Spawn(lb string) error {
	url := fmt.Sprintf("%s/apps/%s/units", h.Endpoint, lb)
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
func (h *TsuruHealer) Terminate(lb, id string) error {
	url := fmt.Sprintf("%s/apps/%s/units", h.Endpoint, lb)
	body := bytes.NewBufferString("1")
	resp, err := request("DELETE", url, h.token, body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error terminating unit: %s", resp.Status)
	}
	return nil
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

func NewTsuruHealer(email, password, endpoint string) *TsuruHealer {
	token, err := getToken(email, password, endpoint)
	if err != nil {
		panic(err.Error())
	}
	return &TsuruHealer{
		seeker:   NewAWSSeeker(),
		Endpoint: endpoint,
		token:    token,
	}
}
