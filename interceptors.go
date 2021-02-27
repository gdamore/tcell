package tcell

// interceptors implements the DrawInterceptAdder interface. It provides
// somewhere to store the callbacks.
type interceptors struct {
	beforeFunc DrawInterceptFunc
	afterFunc  DrawInterceptFunc
}

func (icepts interceptors) before(s Screen, sync bool) {
	if icepts.beforeFunc != nil {
		icepts.beforeFunc(s, sync)
	}
}

func (icepts interceptors) after(s Screen, sync bool) {
	if icepts.afterFunc != nil {
		icepts.afterFunc(s, sync)
	}
}

// AddDrawIntercept wraps the existing draw intercept function with the given
// one.
func (icepts *interceptors) AddDrawIntercept(fn DrawInterceptFunc) {
	icepts.beforeFunc = wrapDrawInterceptFunc(icepts.beforeFunc, fn)
}

// AddDrawInterceptAfter adds the draw interceptor after everything is drawn.
func (icepts *interceptors) AddDrawInterceptAfter(fn DrawInterceptFunc) {
	icepts.afterFunc = wrapDrawInterceptFunc(icepts.afterFunc, fn)
}

func wrapDrawInterceptFunc(oldFn, newFn DrawInterceptFunc) DrawInterceptFunc {
	// directly use the new function if we don't have an old one.
	if oldFn == nil {
		return newFn
	}

	return func(s Screen, sync bool) {
		newFn(s, sync)
		oldFn(s, sync)
	}
}

type noopMutex struct{}

func (noopMutex) Lock()   {}
func (noopMutex) Unlock() {}
