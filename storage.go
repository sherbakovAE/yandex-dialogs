package dialogs

import (
	"time"
)

type State int

type Storage interface {
	GetState(userId string) State
	SetState(userId string, state State) State
	GetData(userId string) interface{}
	SetData(userId string, data interface{}) interface{}
	Delete(userId string)
}

type MemoryStorageItem struct {
	data  interface{}
	state State
	time  time.Time
}
type MemoryStorage map[string]*MemoryStorageItem

// удаление данных старых сессий
func (m MemoryStorage) ClearOld() {
	timeout := time.Duration(6 * time.Hour)
	for key, v := range m {
		if time.Now().Sub(v.time) > timeout {
			delete(m, key)
		}
	}
}

// чтение состояния
func (m MemoryStorage) SetState(userID string, s State) State {

	if _, ok := m[userID]; !ok {
		m.ClearOld() // при создании нового элемента выполняется очистка старых
		m[userID] = &MemoryStorageItem{data: nil, state: s, time: time.Now()}
		return s
	} else {
		m[userID].state = s
		m[userID].time = time.Now()
		return m[userID].state
	}
}

// установка состояния
func (m MemoryStorage) GetState(userID string) State {
	if _, ok := m[userID]; !ok {
		m[userID] = &MemoryStorageItem{data: nil, state: 0, time: time.Now()}
		return 0
	} else {
		m[userID].time = time.Now()
		return m[userID].state
	}
}

// чтение данных сессии
func (m MemoryStorage) SetData(userID string, d interface{}) interface{} {

	if _, ok := m[userID]; !ok {
		m[userID] = &MemoryStorageItem{data: d, state: 0, time: time.Now()}
		return m[userID].data
	} else {
		m[userID].data = d
		m[userID].time = time.Now()
		return m[userID].data
	}
}

// сохранение данных сессии
func (m MemoryStorage) GetData(userID string) interface{} {

	if _, ok := m[userID]; !ok {
		m[userID] = &MemoryStorageItem{data: nil, state: 0, time: time.Now()}
		return m[userID].data
	} else {
		m[userID].time = time.Now()
		return m[userID].data
	}
}

// удаление данных сессии
func (m MemoryStorage) Delete(userID string) {
	if _, ok := m[userID]; ok {
		delete(m, userID)
	}
}

// создание нового хранилища in-memory
func NewMemoryStorage() MemoryStorage {
	m := make(map[string]*MemoryStorageItem, 0)
	return m
}
