package routerapi

type ClientInfo struct {
	ID  string `json:"id"`
	Key string `json:"key"`
}

type RouteInfo struct {
	ID      string    `json:"id"`
	Clients [2]string `json:"clients"`
}

type RouteAddRequest struct {
	Clients [2]string `json:"clients"`
}
