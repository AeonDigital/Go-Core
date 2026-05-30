package xdb

// SetExecutorForTest allows injecting a mock database executor during unit tests.
func (r *DBGeneric[T, PT]) SetExecutorForTest(exec sqlExecutor) {
	r.executor = exec
}
