package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/api"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/config"
)

func VersionCheck(version string) error {
	c := config.Get()
	if c.AcceptableVersion == "" {
		return nil
	}

	if c.AcceptableVersion != version {
		return errors.Newf("router acceptable version is %s, but got %s", c.AcceptableVersion, version)
	}

	return nil
}

func ClientAuth(clientID string, clientKey string) (*api.AuthResponse, error) {
	url := fmt.Sprintf("%s/api/v1/client/auth", config.APIAddr())

	req := map[string]interface{}{
		"client_id":  clientID,
		"client_key": clientKey,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	handler := NewHandler(config.APIAddr(), "", true, 30)
	httpRes, err := handler.request("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode == http.StatusOK {
		var res api.AuthResponse
		if err := json.NewDecoder(httpRes.Body).Decode(&res); err != nil {
			return nil, err
		}

		return &res, nil
	}

	return nil, errors.New("request failed")
}

func GetSessionInfo(sessionID string) (*api.SessionInfo, error) {
	url := fmt.Sprintf("%s/api/v1/session/%s", config.APIAddr(), sessionID)

	handler := NewHandler(config.APIAddr(), "", true, 30)
	httpRes, err := handler.request("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode == http.StatusOK {
		var res api.SessionInfo
		if err := json.NewDecoder(httpRes.Body).Decode(&res); err != nil {
			return nil, err
		}

		return &res, nil
	}

	return nil, errors.New("request failed")
}
