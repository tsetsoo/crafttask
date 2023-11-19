package crafttask

import "sync/atomic"

var currentId atomic.Uint64

func getNewId() uint64 {
	return currentId.Add(1)
}
