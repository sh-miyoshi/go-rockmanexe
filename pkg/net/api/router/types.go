package routerapi

type ClientInfo struct {
	ID  string `yaml:"id"`
	Key string `yaml:"key"`
}

type RouteInfo struct {
	ID      string    `yaml:"id"`
	Clients [2]string `yaml:"clients"`
}
