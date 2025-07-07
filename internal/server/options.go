package server

type ServerOption func(*serverConfig)

type serverConfig struct {
	addr       string
	maxClients int
	logAddr    string
	name       string
	verbose    bool
}

// базовый конструктор конфига
func newConfig(opts ...ServerOption) *serverConfig {
	cfg := &serverConfig{
		addr:       ":8001",
		maxClients: 100,
		logAddr:    "localhost:9000",
		name:       "Server",
		verbose:    false,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// Опции для изменения параметров:

func WithAddr(a string) ServerOption {
	return func(c *serverConfig) { c.addr = a }
}

func WithMaxClients(n int) ServerOption {
	return func(c *serverConfig) { c.maxClients = n }
}

func WithLogger(addr, name string, verbose bool) ServerOption {
	return func(c *serverConfig) {
		c.logAddr = addr
		c.name = name
		c.verbose = verbose
	}
}
