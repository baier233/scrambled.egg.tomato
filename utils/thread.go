package utils

import (
	"sync"
	"sync/atomic"
)

//从https://github.com/Vedadiyan/gothreads/blob/master/threads.go 抄袭的

/*
#include <stdio.h>
#include <stdint.h>

#ifdef _WIN32
#include <windows.h>

typedef HANDLE handle;

extern void callback(int id);

static inline handle Create(int id) {
    HANDLE hThread;
    DWORD dwThreadId;
    DWORD WINAPI thread_func(LPVOID lpParam)
    {
        callback((int)(uintptr_t)lpParam);
        return 0;
    }
    hThread = CreateThread(NULL, 0, thread_func, (LPVOID)(uintptr_t)id, 0, &dwThreadId);
    if (hThread == NULL) {
        printf("Failed to create thread!\n");
    }
    return hThread;
}

static inline void Terminate(handle handle) {
    TerminateThread(handle, 0);
    CloseHandle(handle);
}

static inline void Close(handle handle) {
    CloseHandle(handle);
}

#else
#include <pthread.h>
#include <unistd.h>

typedef pthread_t handle;

extern void callback(int id);

static inline handle Create(int id) {
    pthread_t thread;
    pthread_create(&thread, NULL, (void* (*)(void*))callback, (void*)(uintptr_t)id);
    return thread;
}

static inline void Terminate(handle thread) {
    pthread_cancel(thread);
}

static inline void Close(handle thread) {
    // nothing to do here
    // pthread_t handles are not explicitly closed
}

#endif
*/
import "C"

//export callback
func callback(id C.int) {
	fn, ok := _threadPool.Load(id)
	if ok {
		fn.(func())()
		return
	}
	panic("thread id not found")
}

var _threadPool sync.Map
var _id atomic.Int32

type Thread struct {
	id     C.int
	done   chan bool
	result chan any
	handle C.handle
}

func New(fn func() any) *Thread {
	thread := Thread{
		id:     C.int(_id.Add(1)),
		done:   make(chan bool),
		result: make(chan any),
	}
	_threadPool.Store(thread.id, func() {
		result := fn()
		C.Close(thread.handle)
		thread.done <- true
		thread.result <- result
	})
	return &thread
}

func (t *Thread) Start() {
	t.handle = C.Create(C.int(t.id))
}

func (t *Thread) Stop() {
	C.Terminate(t.handle)
	_threadPool.Delete(t.id)
}

func (t *Thread) Await() <-chan any {
	return t.result
}
