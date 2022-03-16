package forrest

import "errors"

var (
	PoolHasBeenClosedError     = errors.New("pool has been closed")
	TooManyGoroutineBlockError = errors.New("too many goroutine block")
)
