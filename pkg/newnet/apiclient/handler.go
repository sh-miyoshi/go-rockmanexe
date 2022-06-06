package apiclient

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

// Handler ...
type Handler struct {
	client      *http.Client
	serverAddr  string
	accessToken string
}

// NewHandler ...
func NewHandler(serverAddr string, accessToken string, insecure bool, timeout uint) *Handler {
	h := &Handler{
		serverAddr:  serverAddr,
		accessToken: accessToken,
	}

	h.client = createClient(serverAddr, insecure, time.Duration(timeout)*time.Second)

	return h
}

// trimHTTPPrefix trims "http://" and "https://"
func trimHTTPPrefix(addr string) string {
	addr = strings.TrimPrefix(addr, "http://")
	addr = strings.TrimPrefix(addr, "https://")
	return addr
}

func createClient(serverAddr string, insecure bool, timeout time.Duration) *http.Client {
	tlsConfig := tls.Config{
		ServerName: trimHTTPPrefix(serverAddr),
	}

	if insecure {
		tlsConfig.InsecureSkipVerify = true
	}

	tr := &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: &tlsConfig,
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   timeout,
	}
	return client
}

func (h *Handler) request(method string, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if h.accessToken != "" {
		req.Header.Add("Authorization", fmt.Sprintf("bearer %s", h.accessToken))
	}
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	dump, _ := httputil.DumpRequest(req, false)
	logger.Debug("server request dump: %q", dump)

	res, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}

	dump, _ = httputil.DumpResponse(res, true)
	logger.Debug("server response dump: %q", dump)
	return res, nil
}
