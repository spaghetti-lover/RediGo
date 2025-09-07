package data_structure

type SortedSet struct {
	Index       OrderedIndex
	MemberScore map[string]float64
}

// NewSortedSet creates a new SortedSet with the specified index configuration
func NewSortedSet(config IndexConfig) (*SortedSet, error) {
	index, err := NewOrderedIndex(config)
	if err != nil {
		return nil, err
	}

	return &SortedSet{
		Index:       index,
		MemberScore: make(map[string]float64),
	}, nil
}

// NewSortedSetWithBTree creates a SortedSet with B+ Tree (convenience function)
func NewSortedSetWithBTree(degree int) (*SortedSet, error) {
	return NewSortedSet(IndexConfig{
		Type:   IndexTypeBTree,
		Degree: degree,
	})
}

func (ss *SortedSet) Add(score float64, member string) int {
	// Check if member exists.
	result := ss.Index.Add(score, member)
	ss.MemberScore[member] = score
	return result
}

func (ss *SortedSet) GetScore(member string) (float64, bool) {
	score, exists := ss.MemberScore[member]
	return score, exists
}

func (ss *SortedSet) GetRank(member string) int {
	return ss.Index.GetRank(member)
}
