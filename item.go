package gbclient

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/dghubble/sling"
)

// gets an item
func (g *GoBusClient) GetItem(path string) (io.Reader, error) {
	request, err := sling.New().Get(g.serverURL.String()).Path(path).Request()
	if err != nil {
		return nil, err
	}
	response, err := g.client.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Could not get item: %s", response.StatusCode))
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	response.Body.Close()
	return bytes.NewReader(data), nil
}

// create an item
func (g *GoBusClient) WriteItem(path string, data io.Reader) error {
	request, err := sling.New().Put(g.serverURL.String()).Path(path).Body(data).Request()
	if err != nil {
		return err
	}
	response, err := g.client.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusCreated || response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Could not create item: %s", response.StatusCode))
	}
	return nil
}
