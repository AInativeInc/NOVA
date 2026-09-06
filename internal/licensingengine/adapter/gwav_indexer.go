package adapter

type IndexRecord struct {
	EntityID    string
	EntityType  string
	Descriptors []string
	Metadata    map[string]string
}

type GWAVIndexer interface {
	Index(record IndexRecord) error
	HarmonicSearch(query string, limit int) ([]IndexRecord, error)
}
