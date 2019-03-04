package storage

// Storage interface for storing data
type Storage interface {
	Propose(key, value string) error
	Commit(key string) error
	LookupVersions(key string) ([]string, error)
	Lookup(key string, version string) (string, error)
	Close()
}
