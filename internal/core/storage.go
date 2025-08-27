package core

import "github.com/spaghetti-lover/multithread-redis/internal/data_structure"

var dictStore *data_structure.Dict

func init() {
	dictStore = data_structure.CreateDict()
}
