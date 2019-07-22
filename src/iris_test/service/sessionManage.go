package service

import (
"container/list"
"github.com/astaxie/session"
"sync"
"time"
)

var pder = &Provider{list: list.New()}

type SessionStore struct {
	Sid          string                      //session id唯一标示
	TimeAccessed time.Time                   //最后访问时间
	Value        map[interface{}]interface{} //session里面存储的值
}

func (st *SessionStore) Set(key, value interface{}) error {
	st.Value[key] = value
	pder.SessionUpdate(st.Sid)
	return nil
}

func (st *SessionStore) Get(key interface{}) interface{} {
	pder.SessionUpdate(st.Sid)
	if v, ok := st.Value[key]; ok {
		return v
	} else {
		return nil
	}
}

func (st *SessionStore) Delete(key interface{}) error {
	delete(st.Value, key)
	pder.SessionUpdate(st.Sid)
	return nil
}

func (st *SessionStore) SessionID() string {
	return st.Sid
}

type Provider struct {
	lock     sync.Mutex               //用来锁
	sessions map[string]*list.Element //用来存储在内存
	list     *list.List               //用来做gc
}

func (pder *Provider) SessionInit(sid string) (session.Session, error) {
	pder.lock.Lock()
	defer pder.lock.Unlock()
	v := make(map[interface{}]interface{}, 0)
	newsess := &SessionStore{Sid: sid, TimeAccessed: time.Now(), Value: v}
	element := pder.list.PushBack(newsess)
	pder.sessions[sid] = element
	return newsess, nil
}

func (pder *Provider) SessionRead(sid string) (session.Session, error) {
	if element, ok := pder.sessions[sid]; ok {
		return element.Value.(*SessionStore), nil
	} else {
		sess, err := pder.SessionInit(sid)
		return sess, err
	}
	return nil, nil
}

func (pder *Provider) SessionDestroy(sid string) error {
	if element, ok := pder.sessions[sid]; ok {
		delete(pder.sessions, sid)
		pder.list.Remove(element)
		return nil
	}
	return nil
}

func (pder *Provider) SessionGC(maxlifetime int64) {
	pder.lock.Lock()
	defer pder.lock.Unlock()

	for {
		element := pder.list.Back()
		if element == nil {
			break
		}
		if (element.Value.(*SessionStore).TimeAccessed.Unix() + maxlifetime) < time.Now().Unix() {
			pder.list.Remove(element)
			delete(pder.sessions, element.Value.(*SessionStore).Sid)
		} else {
			break
		}
	}
}

func (pder *Provider) SessionUpdate(sid string) error {
	pder.lock.Lock()
	defer pder.lock.Unlock()
	if element, ok := pder.sessions[sid]; ok {
		element.Value.(*SessionStore).TimeAccessed = time.Now()
		pder.list.MoveToFront(element)
		return nil
	}
	return nil
}

func init() {
	pder.sessions = make(map[string]*list.Element, 0)
	session.Register("memory", pder)
}