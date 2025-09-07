package data_structure

// OrderedIndex defines the interface for ordered data structures
type OrderedIndex interface {
	// Add adds an item with score and member
	// Returns 1 if new item added, 0 if item updated
	Add(score float64, member string) int

	// GetRank returns the rank (0-based index) of a member
	// Return -1 if member not found
	GetRank(member string) int
}

// IndexType represents the type of index to create
type IndexType string

const (
	IndexTypeBTree IndexType = "btree"
)

// IndexConfig holds configuration for creating indexes
type IndexConfig struct {
	Type   IndexType
	Degree int // For B+ Tree
}
