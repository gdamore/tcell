package tcell

// interceptors implements the DrawInterceptAdder interface. It provides
// somewhere to store the callbacks.
type interceptors struct {
	beforeFunc DrawInterceptFunc
	afterFunc  DrawInterceptFunc
}

func (icepts interceptors) before(s Screen, sync bool) bool {
	if icepts.beforeFunc != nil {
		return icepts.beforeFunc(s, sync)
	}
	return false
}

func (icepts interceptors) after(s Screen, sync bool) bool {
	if icepts.afterFunc != nil {
		return icepts.afterFunc(s, sync)
	}
	return false
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

	return func(s Screen, sync bool) bool {
		a := newFn(s, sync)
		b := oldFn(s, sync)
		return a || b
	}
}

type noopMutex struct{}

func (noopMutex) Lock()   {}
func (noopMutex) Unlock() {}
