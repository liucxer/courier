package envconf

// mask value string for security value
type SecurityStringer interface {
	SecurityString() string
}

type EnvVar struct {
	KeyPath string
	Value   string
	Mask    string

	Optional      bool
	IsUpstream    bool
	IsCopy        bool
	IsExpose      bool
	IsHealthCheck bool
}

func (envVar *EnvVar) metaFromFlags(flags map[string]bool) {
	envVar.Optional = false

	if flags["opt"] {
		envVar.Optional = true
	}
	if flags["upstream"] {
		envVar.IsUpstream = true
	}

	if flags["copy"] {
		envVar.IsCopy = true
	}

	if flags["expose"] {
		envVar.IsExpose = true
	}

	if flags["healthCheck"] {
		envVar.IsHealthCheck = true
	}
}

func (envVar *EnvVar) Key(prefix string) string {
	return prefix + "__" + envVar.KeyPath
}
