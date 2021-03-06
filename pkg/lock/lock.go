package lock

import (
	"github.com/aldor007/mort/pkg/response"
)

// Lock is responding for collapsing request for same object
type Lock interface {
	// Lock try  get a lock for given key
	Lock(key string) (observer LockResult, acquired bool)
	// Release remove lock for given key
	Release(key string)
	// NotifyAndRelease remove lock for given key and notify all clients waiting for result
	NotifyAndRelease(key string, res *response.Response)
}

// LockResult contain struct
type LockResult struct {
	ResponseChan chan *response.Response // channel on which you get response
	Cancel       chan bool               // channel for notify about cancel of waiting
}

type lockData struct {
	Key         string
	notifyQueue []LockResult
}

// AddWatcher add next request waiting for lock to expire or return result
func (l *lockData) AddWatcher() LockResult {
	d := LockResult{}
	d.ResponseChan = make(chan *response.Response, 1)
	d.Cancel = make(chan bool, 1)
	l.notifyQueue = append(l.notifyQueue, d)
	return d
}

// NewNopLock create lock that do nothing
func NewNopLock() *NopLock {
	return &NopLock{}
}

// NopLock will never  collapse any request
type NopLock struct {
}

// Lock always return that lock was acquired
func (l *NopLock) Lock(_ string) (LockResult, bool) {
	return LockResult{}, true
}

// Release do nothing
func (l *NopLock) Release(_ string) {

}

// NotifyAndRelease do nothing
func (l *NopLock) NotifyAndRelease(_ string, _ *response.Response) {

}
