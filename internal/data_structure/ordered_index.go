package data_structure

type OrderedIndex interface {
	Add(score float64, member string) int
	GetScore(member string) (float64, bool)
	GetRank(member string) int
}
