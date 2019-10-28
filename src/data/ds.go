package data

import (
	"sync"
)

type ds struct {
}

var instanceDS *ds
var onceDS sync.Once

func GetDS() *ds {
	onceDS.Do(func() {
		instanceDS = &ds{}
	})
	return instanceDS
}
