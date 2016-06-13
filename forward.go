package gbclient

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
)

// creates a forward on gobus
// the forward is created on goBusPath and forwards all requests to forwardPath
func (g *GoBusClient) CreateForward(goBusPath, forwardPath string) error {
	if err := g.CreateCollection(goBusPath); err != nil {
		return err
	}
	type Forward struct {
		URL string `json:"url"`
	}
	fmt.Printf(goBusPath)
	sl := sling.New().Put(g.serverURL.String()).Path(goBusPath).Path("_forward")
	request, err := sl.BodyJSON(Forward{forwardPath}).Request()
	if err != nil {
		return err
	}
	r, err := g.client.Do(request)
	if err != nil {
		return err
	}
	if r.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Could not create forward: %s", r.StatusCode))
	}
	return nil
}

// deletes a forward on gobus
func (g *GoBusClient) DeleteForward(goBusPath string) error {
	request, err := sling.New().Delete(g.serverURL.String()).Path(goBusPath).Path("_forward").Request()
	if err != nil {
		return err
	}
	r, err := g.client.Do(request)
	if err != nil {
		return err
	}
	if r.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Could not delete forward: %s", r.StatusCode))
	}
	return nil
}
