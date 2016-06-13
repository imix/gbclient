package gbclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/dghubble/sling"
)

// create a collection
func (g *GoBusClient) CreateCollection(path string) error {
	request, err := sling.New().Put(g.serverURL.String()).Path(path).Request()
	if err != nil {
		return err
	}
	response, err := g.client.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusCreated {
		return errors.New(fmt.Sprintf("Could not create collection: %s", response.StatusCode))
	}
	return nil
}

// add an item to a collection
func (g *GoBusClient) AddToCollection(path string, data io.Reader) error {
	request, err := sling.New().Post(g.serverURL.String()).Path(path).Body(data).Request()
	if err != nil {
		return err
	}
	response, err := g.client.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusCreated {
		return errors.New(fmt.Sprintf("Could not add to collection: %s", response.StatusCode))
	}
	return nil
}

// get the list of items in a collection
func (g *GoBusClient) GetCollection(path string) ([]string, error) {
	request, err := sling.New().Get(g.serverURL.String()).Path(path).Request()
	if err != nil {
		return nil, err
	}
	response, err := g.client.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Could not get collection: %s", response.StatusCode))
	}
	var result []string
	data, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
