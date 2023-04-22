package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/net/config"
)

type AuthResponse struct {
	SessionID string `json:"session_id"`
}

type SessionInfo struct {
	ID            string `json:"id"`
	OwnerUserID   string `json:"owner_user_id"`
	OwnerClientID string `json:"owner_client_id"`
	GuestUserID   string `json:"guest_user_id"`
	GuestClientID string `json:"guest_client_id"`
}

func VersionCheck(version string) error {
	c := config.Get()
	if c.AcceptableVersion == "" {
		return nil
	}

	if c.AcceptableVersion != version {
		return fmt.Errorf("router acceptable version is %s, but got %s", c.AcceptableVersion, version)
	}

	return nil
}

func ClientAuth(clientID string, clientKey string) (*AuthResponse, error) {
	c := config.Get()
	url := fmt.Sprintf("%s/api/v1/client/auth", c.APIAddr)

	req := map[string]interface{}{
		"client_id":  clientID,
		"client_key": clientKey,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	handler := NewHandler(c.APIAddr, "", true, 30)
	httpRes, err := handler.request("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode == http.StatusOK {
		var res AuthResponse
		if err := json.NewDecoder(httpRes.Body).Decode(&res); err != nil {
			return nil, err
		}

		return &res, nil
	}

	return nil, fmt.Errorf("request failed")
}

func GetSessionInfo(sessionID string) (*SessionInfo, error) {
	c := config.Get()
	url := fmt.Sprintf("%s/api/v1/session/%s", c.APIAddr, sessionID)

	handler := NewHandler(c.APIAddr, "", true, 30)
	httpRes, err := handler.request("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode == http.StatusOK {
		var res SessionInfo
		if err := json.NewDecoder(httpRes.Body).Decode(&res); err != nil {
			return nil, err
		}

		return &res, nil
	}

	return nil, fmt.Errorf("request failed")
}
