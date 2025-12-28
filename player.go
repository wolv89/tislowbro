package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"time"
)

type Player struct {
	address, password string
	client            *http.Client
}

const (
	STATUS_PATH = "/requests/status.xml"
)

func NewPlayer(addr, pass string) Player {
	return Player{
		address:  addr,
		password: pass,
		client: &http.Client{
			Timeout: 2 * time.Second,
		},
	}
}

// Only handles GET requests
func (p Player) apicall(url string) error {

	if len(url) == 0 {
		return fmt.Errorf("require url to call")
	}

	req, _ := http.NewRequest("GET", p.address+url, nil)
	req.SetBasicAuth("", p.password)

	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("call failed with status: %d", resp.StatusCode)
	}

	return nil

}

func (p Player) IsRunning() error {
	return p.apicall(STATUS_PATH)
}

func (p Player) Reset() error {

	params := url.Values{}
	params.Add("command", "pl_empty")

	url := fmt.Sprintf("%s?%s", STATUS_PATH, params.Encode())

	return p.apicall(url)

}

func (p Player) Queue(video Video) error {

	rawPath, _ := exec.Command("wslpath", "-w", video.path).Output()
	winPath := strings.TrimSpace(string(rawPath))

	uri := "file:///" + winPath

	params := url.Values{}
	params.Add("command", "in_enqueue")
	params.Add("input", uri)

	url := fmt.Sprintf("%s?%s", STATUS_PATH, params.Encode())
	url = strings.ReplaceAll(url, "+", "%20")

	return p.apicall(url)

}
