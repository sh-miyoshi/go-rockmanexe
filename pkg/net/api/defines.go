package api

type AuthRequest struct {
	ClientID  string `json:"client_id"`
	ClientKey string `json:"client_key"`
}

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
