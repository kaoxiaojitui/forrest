// Run! Forrest, Run!

package forrest

var defaultAntsPool = NewDefaultPool()

func Submit(task func()) error {
	return defaultAntsPool.Submit(task)
}

func Running() int32 {
	return defaultAntsPool.Running()
}