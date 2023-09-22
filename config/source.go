package config

// KV is config key value.
type KV struct {
	Key    string
	Format string
	Value  []byte
}

// Source is config source.
type Source interface {
	Load() ([]*KV, error)
	Watch() (Watcher, error)
}

// Watcher watches a source for changes.
type Watcher interface {
	Next() ([]*KV, error)
	Stop() error
}
