package internal

type Config struct {
	Port string
}

func (c *Config) GetPort() string {
	return c.Port
}
