package files

import (
	"sync"
)

var (
	fileLocksMu sync.Mutex
	fileLocks   = make(map[string]*sync.Mutex)
)

func getFileLock(path string) *sync.Mutex {
	fileLocksMu.Lock()
	defer fileLocksMu.Unlock()
	lock, ok := fileLocks[path]
	if !ok {
		lock = &sync.Mutex{}
		fileLocks[path] = lock
	}
	return lock
}
