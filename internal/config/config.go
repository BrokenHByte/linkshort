package config

type ServerConfig struct {
	ServerAddr string
	PrefixLink string
}

func NewServerConfig() *ServerConfig {
	var c ServerConfig
	c.ServerAddr = ":8080"
	c.PrefixLink = "http://localhost:8080"
	return &c
}
