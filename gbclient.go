package gbclient

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"path"
	"time"
)

// GoBus represents a gobus instance and is used to generate queries
type GoBusClient struct {
	serverURL     *url.URL
	ownURL        *url.URL
	client        *http.Client
	hookReceivers []hookReceiver
	nextHookId    int
}

// returns a new GoBusClient
func NewGoBusClient(host, port string) *GoBusClient {
	serverURL, err := url.Parse("http://" + host + ":" + port)
	if err != nil {
		log.Fatalf("Could not parse url: %s", err)
	}
	client := &http.Client{}
	return &GoBusClient{
		serverURL:  serverURL,
		client:     client,
		nextHookId: 0,
	}
}

// creates an url from the given path and the baseUrl
func (g *GoBusClient) getServerURL(relPath string) string {
	newURL := *g.serverURL
	newURL.Path = path.Join(newURL.Path, relPath)

	log.Printf(newURL.Path)
	return newURL.String()
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's taken from net/http/server.go. It's used so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

// starts the server to receive callbacks from registered hooks
// if host is empty, listens on localhost
// if port ist empty, automatically assigns a port
func (g *GoBusClient) StartServer(host, port string) error {
	addrString := host
	addrString += ":"
	if port == "" {
		addrString += "0"
	} else {
		addrString += port
	}
	listener, err := net.Listen("tcp", addrString)
	if err != nil {
		return err
	}
	newOwnURL, err := url.Parse("http://" + listener.Addr().String() + "/")
	if err != nil {
		return err
	}
	g.ownURL = newOwnURL
	mux := http.NewServeMux()
	mux.HandleFunc("/", g.getHookHandler())
	go func() {
		server := &http.Server{Addr: addrString, Handler: mux}
		err = server.Serve(tcpKeepAliveListener{listener.(*net.TCPListener)})
	}()
	if err != nil {
		return err
	}
	return nil
}
