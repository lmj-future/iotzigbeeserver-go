package http

import (
	"context"
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"syscall"
	"time"

	"github.com/creasty/defaults"
	"github.com/docker/go-connections/tlsconfig"
	"github.com/gorilla/mux"
	"github.com/h3c/iotzigbeeserver-go/globalconstant/globallogger"
)

// Params http params
type Params = map[string]string

// Headers http headers
type Headers = http.Header

// Server http server
type Server struct {
	cfg    ServerInfo
	svr    *http.Server
	uri    *url.URL
	addr   string
	auth   func(u, p string) bool
	router *mux.Router
}

// ParseURL parses a url string
func ParseURL(addr string) (*url.URL, error) {
	if strings.HasPrefix(addr, "unix://") {
		parts := strings.SplitN(addr, "://", 2)
		return &url.URL{
			Scheme: parts[0],
			Host:   parts[1],
		}, nil
	}
	return url.Parse(addr)
}

// NewTLSServerConfig loads tls config for server
func NewTLSServerConfig(c Certificate) (*tls.Config, error) {
	if c.Cert == "" && c.Key == "" {
		return nil, nil
	}
	return tlsconfig.Server(tlsconfig.Options{CAFile: c.CA, KeyFile: c.Key, CertFile: c.Cert, ClientAuth: tls.VerifyClientCertIfGiven})
}

// NewServer creates a new http server
func NewServer(c ServerInfo, a func(u, p string) bool) (*Server, error) {
	defaults.Set(&c)

	uri, err := ParseURL(c.Address)
	if err != nil {
		return nil, err
	}
	tls, err := NewTLSServerConfig(c.Certificate)
	if err != nil {
		return nil, err
	}
	router := mux.NewRouter()
	return &Server{
		cfg:    c,
		auth:   a,
		uri:    uri,
		router: router,
		svr: &http.Server{
			WriteTimeout: c.Timeout,
			ReadTimeout:  c.Timeout,
			TLSConfig:    tls,
			Handler:      router,
		},
	}, nil
}

// Handle handle requests
func (s *Server) Handle(handle func(Params, []byte) ([]byte, error), method, path string, params ...string) {
	s.router.HandleFunc(path, func(res http.ResponseWriter, req *http.Request) {
		globallogger.Log.Infof("[%s] %s", req.Method, req.URL.String())
		// if s.auth != nil {
		// 	if !s.auth(req.Header.Get(headerKeyUsername), req.Header.Get(headerKeyPassword)) {
		// 		http.Error(res, errAccountUnauthorized.Error(), 401)
		// 		globallogger.Log.Errorf("[%s] %s %s", req.Method, req.URL.String(), errAccountUnauthorized.Error())
		// 		return
		// 	}
		// }
		var err error
		var reqBody []byte
		if req.Body != nil {
			defer req.Body.Close()
			reqBody, err = ioutil.ReadAll(req.Body)
			if err != nil {
				http.Error(res, err.Error(), 400)
				globallogger.Log.Errorf("[%s] %s %s", req.Method, req.URL.String(), err.Error())
				return
			}
		}
		resBody, err := handle(mux.Vars(req), reqBody)
		if err != nil {
			http.Error(res, err.Error(), 400)
			globallogger.Log.Errorf("[%s] %s %s", req.Method, req.URL.String(), err.Error())
			return
		}
		if resBody != nil {
			res.Write(resBody)
		}
	}).Methods(method).Queries(params...)
}

// Start starts server
func (s *Server) Start() error {
	if s.uri.Scheme == "unix" {
		if err := syscall.Unlink(s.uri.Host); err != nil {
			globallogger.Log.Errorf(err.Error())
		}
	}

	l, err := net.Listen(s.uri.Scheme, s.uri.Host)
	if err != nil {
		return err
	}
	if s.uri.Scheme == "tcp" {
		l = tcpKeepAliveListener{l.(*net.TCPListener)}
	}
	s.addr = l.Addr().String()
	go s.svr.Serve(l)
	return nil
}

// Close closese server
func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.svr.IdleTimeout)
	defer cancel()
	return s.svr.Shutdown(ctx)
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}
