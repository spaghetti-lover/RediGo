package data_structure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZSet_GetRank(t *testing.T) {
	bptree := NewBPlusTree(3)

	ss := NewSortedSet(bptree)
	ss.Add(20.0, "k2")
	ss.Add(40.0, "k4")
	ss.Add(10.0, "k1")
	ss.Add(30.0, "k3")
	ss.Add(50.0, "k5")
	ss.Add(60.0, "k6")
	ss.Add(80.0, "k8")
	ss.Add(70.0, "k7")

	rank := ss.GetRank("k1")
	score, _ := ss.GetScore("k1")
	assert.EqualValues(t, 0, rank)
	assert.EqualValues(t, 10.0, score)

	rank = ss.GetRank("k2")
	score, _ = ss.GetScore("k2")
	assert.EqualValues(t, 1, rank)
	assert.EqualValues(t, 20.0, score)

	rank = ss.GetRank("k3")
	score, _ = ss.GetScore("k3")
	assert.EqualValues(t, 2, rank)
	assert.EqualValues(t, 30.0, score)

	rank = ss.GetRank("k4")
	score, _ = ss.GetScore("k4")
	assert.EqualValues(t, 3, rank)
	assert.EqualValues(t, 40.0, score)

	rank = ss.GetRank("k5")
	score, _ = ss.GetScore("k5")
	assert.EqualValues(t, 4, rank)
	assert.EqualValues(t, 50.0, score)

	rank = ss.GetRank("k6")
	score, _ = ss.GetScore("k6")
	assert.EqualValues(t, 5, rank)
	assert.EqualValues(t, 60.0, score)

	rank = ss.GetRank("k7")
	score, _ = ss.GetScore("k7")
	assert.EqualValues(t, 6, rank)
	assert.EqualValues(t, 70.0, score)

	rank = ss.GetRank("k8")
	score, _ = ss.GetScore("k8")
	assert.EqualValues(t, 7, rank)
	assert.EqualValues(t, 80.0, score)
}
