package gbclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

// Hook Receiver describes an endpoint for a hook event
type hookReceiver struct {
	Path string
	Func func(gbc *GoBusClient, r *HookResponse)
}

type hookRequest struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type HookResponse struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Path   string `json:"path"`
	Method string `json:"method"`
	Item   bool   `json:"item"` // is the affected resource an item or a collection
}

func (g *GoBusClient) getHookURL() *url.URL {
	newURL := *g.ownURL
	name := strconv.Itoa(g.nextHookId)
	newURL.Path = path.Join(newURL.Path, name)
	return &newURL
}

// create a hook for a certain resource
func (g *GoBusClient) HookResource(resPath string, receiverFunc func(gbc *GoBusClient, h *HookResponse)) error {
	hookURL := g.getHookURL()
	data, err := json.Marshal(hookRequest{
		Name: "a_hook",
		Url:  hookURL.String(),
	})
	if err != nil {
		return err
	}
	response, err := http.Post(g.getServerURL(path.Join(resPath, "_hooks")), "", bytes.NewReader(data))
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusCreated {
		return errors.New(fmt.Sprintf("Resource could not be hooked: %d", response.StatusCode))
	}
	newReceiver := hookReceiver{
		Path: hookURL.Path,
		Func: receiverFunc,
	}
	g.hookReceivers = append(g.hookReceivers, newReceiver)
	return nil
}

// creates a handler that handles all defined events for hooks
func (g *GoBusClient) getHookHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, hr := range g.hookReceivers {
			if r.URL.Path == path.Join(g.ownURL.Path, hr.Path) {
				var hookResp HookResponse
				data, err := ioutil.ReadAll(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				err = json.Unmarshal(data, &hookResp)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				go hr.Func(g, &hookResp)
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, "OK")
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Not Found")
	}
}
