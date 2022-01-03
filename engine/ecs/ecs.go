package ecs

import "reflect"

type Id uint32

type Component interface {
	ComponentSet(interface{})
}

type BasicStorage struct {
	list map[Id]interface{}
}

func NewBasicStorage() *BasicStorage {
	return &BasicStorage{
		list: make(map[Id]interface{}),
	}
}

func (BasicStorage *BasicStorage) Read(id Id) (interface{}, bool) {
	val, ok := BasicStorage.list[id]
	return val, ok
}

func (BasicStorage *BasicStorage) Write(id Id, val interface{}) {
	BasicStorage.list[id] = val
}

type Engine struct {
	reg       map[string]*BasicStorage
	idCounter Id
}

func NewEngine() *Engine {
	return &Engine{
		reg:       make(map[string]*BasicStorage),
		idCounter: 0,
	}
}

func (engine *Engine) NewId() Id {
	id := engine.idCounter
	engine.idCounter++
	return id
}

func name(t interface{}) string {
	name := reflect.TypeOf(t).String()
	if name[0] == '*' {
		return name[1:]
	}
	return name
}

func GetStorage(engine *Engine, t interface{}) *BasicStorage {
	name := name(t)
	storage, ok := engine.reg[name]
	if !ok {
		engine.reg[name] = NewBasicStorage()
		storage, _ = engine.reg[name]
	}
	return storage
}

func Read(engine *Engine, id Id, val Component) bool {
	storage := GetStorage(engine, val)
	newVal, ok := storage.Read(id)
	if ok {
		val.ComponentSet(newVal)
	}
	return ok
}

func Write(engine *Engine, id Id, val interface{}) {
	storage := GetStorage(engine, val)
	storage.Write(id, val)
}

func Each(engine *Engine, val interface{}, f func(id Id, a interface{})) {
	storage := GetStorage(engine, val)
	for id, a := range storage.list {
		f(id, a)
	}
}
