package core

import "github.com/spaghetti-lover/multithread-redis/internal/data_structure"

var dictStore *data_structure.Dict
var zsetStore map[string]*data_structure.SortedSet
var setStore map[string]*data_structure.SimpleSet

func init() {
	dictStore = data_structure.CreateDict()
	zsetStore = make(map[string]*data_structure.SortedSet)
	setStore = make(map[string]*data_structure.SimpleSet)
}
