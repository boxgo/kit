package lock

import "time"

type (
	// Lock 竞争锁
	Lock interface {
		Lock(string, time.Duration) (bool, error)
		IsLocked(string) (bool, error)
		UnLock(string) error
	}
)
