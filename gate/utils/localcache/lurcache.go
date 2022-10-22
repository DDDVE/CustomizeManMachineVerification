package localcache

import (
	"sync"
	"time"
)

//LUR算法结合map针对"过期时间一致"和"节点数量有限"有两种实现方案
//这里实现"过期时间一致"对应的LUR算法,"节点数量有限"实现方式不一，暂不写了
//"过期时间一致"也可以结合切片数组实现，更加节省内存

type LURcache struct {
	button      bool //默认get时自动续期，设置为false则关闭自动续期
	startNode   *node
	endNode     *node
	m           map[string]node
	unifyExpire time.Duration //统一的过期时间
	lock        sync.RWMutex
}

// 双向链表中的节点
type node struct {
	preNode *node
	sufNode *node
	expire  time.Time
	key     string
	value   interface{}
}

func NewLRUcache(expire time.Duration) *LURcache {
	startNode := &node{}
	endNode := &node{}
	startNode.sufNode = endNode
	endNode.preNode = startNode
	return &LURcache{
		startNode:   startNode,
		endNode:     endNode,
		button:      true,
		unifyExpire: expire,
	}
}

//设置get时是否自动续期
func (l *LURcache) SetButton(b bool) {
	l.button = b
}

func (l *LURcache) Add(key string, val interface{}) {
	newNode := node{
		expire: time.Now().Add(l.unifyExpire),
		key:    key,
		value:  val,
	}
	l.lock.Lock()
	l.startNode.sufNode.preNode = &newNode
	newNode.preNode = l.startNode
	newNode.sufNode = l.startNode.sufNode
	l.startNode.sufNode = &newNode
	l.m[key] = newNode
	l.lock.Unlock()

	go l.deleteExpire()

	return
}

func (l *LURcache) Get(key string) (val interface{}, expired bool) {
	l.lock.RLock()
	curNode := l.m[key]
	l.lock.RUnlock()

	if l.button {
		l.lock.Lock()
		curNode.preNode.sufNode = curNode.sufNode
		curNode.sufNode.preNode = curNode.preNode
		l.startNode.sufNode.preNode = &curNode
		curNode.preNode = l.startNode
		curNode.sufNode = l.startNode.sufNode
		l.startNode.sufNode = &curNode
		l.m[key] = curNode
		l.lock.Unlock()
	}

	go l.deleteExpire()

	return curNode.value, curNode.expire.Before(time.Now())
}

//删除所有过期node
func (l *LURcache) deleteExpire() {
	curNode := l.endNode.preNode
	l.lock.Lock()
	defer l.lock.Unlock()
	for {
		if curNode == nil || curNode == l.startNode {
			break
		}
		if time.Now().Before(curNode.expire) {
			break
		} else {
			delete(l.m, curNode.key)
			curNode = curNode.preNode
		}
	}

	l.endNode.preNode = curNode
}
