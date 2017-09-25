package source

import (
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

func NewUrlSourcer() Sourcer {
	return &url{
		client: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

type url struct {
	client *http.Client
}

func (u *url) Get(path string) ([]byte, error) {
	resp, err := u.client.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	return ioutil.ReadAll(resp.Body)
}
