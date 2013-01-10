package healer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Healer interface {
	Heal() error
	Spawn(lb string) error
	Terminate(id string) error
}

type TsuruHealer struct {
	Endpoint string
}

// Heal iterates through down instances, terminate then
// and spawn new ones to replace the terminated.
func (h *TsuruHealer) Heal() error {
	// instances, err := h.seeker.SeekUnhealthyInstances()
	// if err != nil {
	//     return err
	// }
	// for _, instance := range instances {
	//     err := h.Terminate(instance.Id)
	//     if err != nil {
	//         // should really stop here?
	//         return err
	//     }
	//     err := h.Spawn(instance)
	//     if err != nil {
	//         // should really stop here?
	//         return err
	//     }
	// }
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
	resp, err := http.Post(url, "text/plain", body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error spawning unit: %s", resp.Status)
	}
	return nil
}

// Calls tsuru remove-unit endpoint
func (h *TsuruHealer) Terminate(id string) error {
	return nil
}

func NewTsuruHealer(email, password, endpoint string) *TsuruHealer {
	getToken(email, password, endpoint)
	return &TsuruHealer{Endpoint: endpoint}
}
