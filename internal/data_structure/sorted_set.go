package data_structure

type SortedSet struct {
	Index       OrderedIndex
	MemberScore map[string]float64
}

func NewSortedSet(index OrderedIndex) *SortedSet {
	return &SortedSet{
		Index:       index,
		MemberScore: make(map[string]float64),
	}
}

func (ss *SortedSet) Add(score float64, member string) int {
	return ss.Index.Add(score, member)
}

func (ss *SortedSet) GetScore(member string) (float64, bool) {
	return ss.Index.GetScore(member)
}

func (ss *SortedSet) GetRank(member string) int {
	return ss.Index.GetRank(member)
}
