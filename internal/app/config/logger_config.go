package config

type LoggerConfig struct {
	Pretty       bool   // human-readable logs (for dev)
	Level        string // debug/info/warn/error
	BodyLogMax   int    // max length body for logs
	SampleErrors bool   // log body on errors
}

func LoadLoggerConfig(env Environment) *LoggerConfig {
	cfg := &LoggerConfig{
		BodyLogMax:   1024,
		SampleErrors: true,
	}

	switch env {
	case Dev:
		cfg.Pretty = true
		cfg.Level = "debug"
	case Stage:
		cfg.Pretty = false
		cfg.Level = "info"
	case Prod:
		cfg.Pretty = false
		cfg.Level = "info"
		cfg.SampleErrors = false
	}

	return cfg
}
