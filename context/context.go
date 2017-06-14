package context

import (
	"crypto/rsa"
	"sync"
)

var (
	globalContextLock sync.Mutex
	globalContext     Context
)

type Context struct {
	PrivateKey *rsa.PrivateKey
}

func GlobalContext() Context {
	globalContextLock.Lock()
	defer globalContextLock.Unlock()
	return globalContext
}

func SetGlobalContext(ctx Context) {
	globalContextLock.Lock()
	defer globalContextLock.Unlock()
	globalContext = ctx
}
