package utils

import "sync"

type StringContainer struct {
	data map[string]bool
	mu   sync.RWMutex
}

// NewStringContainer 创建一个新的字符串容器
func NewStringContainer() *StringContainer {
	return &StringContainer{
		data: make(map[string]bool),
	}
}

// Add 方法用于向容器添加一个字符串
func (c *StringContainer) Add(str string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[str] = true
}

// Contains 方法用于检查容器是否包含指定的字符串
func (c *StringContainer) Contains(str string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, exists := c.data[str]
	return exists
}
