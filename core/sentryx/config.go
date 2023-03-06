package sentryx

type Config struct {
	Dsn string `yaml:"Dsn"`
}

func (c Config) Available() bool {
	return len(c.Dsn) > 0
}
