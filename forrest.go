// Run! Forrest, Run!

package forrest

var defaultPool = NewDefaultPool()

func Submit(task func()) error {
	return defaultPool.Submit(task)
}

func Running() int32 {
	return defaultPool.Running()
}
