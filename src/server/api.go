package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/yslou/frigo-server/lib/model"
	"github.com/yslou/frigo-server/lib/config"
	"crypto/tls"
	"os"
	"github.com/syncthing/syncthing/lib/tlsutil"
)

var (
	ErrNotSupportedExtension = fmt.Errorf("Not supported extenstion")
	validExt = []string {".html", ".js", ".jpg", ".png", ".svg"}
)

type apiService struct {
	model       *model.Model
	cfg         *config.Config
	stop        chan interface{}
	listener    net.Listener
	listenerMut sync.Mutex
}

func NewApiService(m *model.Model, cfg *config.Config) (*apiService, error) {
	svc := &apiService {
		model:          m,
		cfg:            cfg,
		stop:           make(chan interface{}),
		listenerMut:    sync.Mutex{},
	}

	var err error
	svc.listener, err = svc.getListener(cfg.Api)

	return svc, err
}

func (s *apiService) getListener(cfg config.ApiConfig) (net.Listener, error) {
	httpsRSABits := 2048
	tlsDefaultCommonName := "frigo"
	cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		l.Infoln("Loading HTTPS certificate:", err)
		l.Infoln("Creating new HTTPS certificate")

		var name string
		name, err = os.Hostname()
		if err != nil {
			name = tlsDefaultCommonName
		}

		cert, err = tlsutil.NewCertificate(cfg.CertFile, cfg.KeyFile, name, httpsRSABits)
	}
	if err != nil {
		return nil, err
	}
	tlsCfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS10, // No SSLv3
		CipherSuites: []uint16{
			// No RC4
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
			tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
		},
	}

	rawListener, err := net.Listen("tcp", cfg.RawAddr)
	if err != nil {
		return nil, err
	}

	tlsListener := &tlsutil.DowngradingListener{rawListener, tlsCfg}
	return tlsListener, err
}

func httpSendJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(v)
}

// check if the path has valid extension
func isValidExt(path string) bool {
	for _, ext := range validExt {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}
	return false
}

func loadPage(path string) (body []byte, err error) {
	if strings.HasSuffix(path, "/") {
		path = filepath.Join(path, "index.html")
	}
	if !isValidExt(path) {
		return body, ErrNotSupportedExtension
	}
	f := filepath.Join("www", path)
	return ioutil.ReadFile(f)
}

func staticPageHandler(w http.ResponseWriter, r *http.Request) {
	b, err := loadPage(r.URL.Path)
	if err == nil {
		w.Write(b)
	} else {
		http.Error(w, err.Error(), 404)
	}
}

func (s *apiService) ping(w http.ResponseWriter, r *http.Request) {
	httpSendJSON(w, map[string]string{"ping": "pong"})
}

func redirectToHTTPSMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.TLS == nil {
			// Redirect HTTP requests to HTTPS
			r.URL.Host = r.Host
			r.URL.Scheme = "https"
			http.Redirect(w, r, r.URL.String(), http.StatusTemporaryRedirect)
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

func noCacheMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=0, no-cache, no-store")
		w.Header().Set("Expires", time.Now().UTC().Format(http.TimeFormat))
		w.Header().Set("Pragma", "no-cache")
		h.ServeHTTP(w, r)
	})
}

func (s *apiService) Serve() {
	s.listenerMut.Lock()
	listener := s.listener
	s.listenerMut.Unlock()

	mux := http.NewServeMux()
	mux.HandleFunc("/ping:", s.ping)
	mux.HandleFunc("/", staticPageHandler)

	handler := noCacheMiddleware(mux)
	handler = redirectToHTTPSMiddleware(handler);

	srv := http.Server {
		Handler: handler,
	}

	apil.Infoln("API listening on ", listener.Addr())
	err := srv.Serve(listener)

	select {
	case e := <- s.stop:
		l.Infof("API: shutdown?? ", e)
	case <- time.After(time.Second):
		l.Warnln("API:", err)
	}
}

func (s *apiService) Stop() {
	s.listenerMut.Lock()
	listener := s.listener
	s.listenerMut.Unlock()

	close(s.stop)

	if listener != nil {
		listener.Close()
	}
}
