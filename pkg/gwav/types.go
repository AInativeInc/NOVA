package gwav

type HarmonicQuery struct {
	SemanticDescriptor string            `json:"semantic_descriptor"`
	Constraints        map[string]string `json:"constraints,omitempty"`
}

type HarmonicMatch struct {
	EntityID   string             `json:"entity_id"`
	EntityType string             `json:"entity_type"`
	Score      float64            `json:"score"`
	Metadata   map[string]string  `json:"metadata,omitempty"`
	Vectors    map[string]float64 `json:"vectors,omitempty"`
}
