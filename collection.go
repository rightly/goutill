package goutill

import "reflect"

type collection struct{}

var Collection = collection{}

type Map struct {
	reflect.Value
}

func (collection) MakeMap(key, value interface{}) *Map {
	keyType := reflect.TypeOf(key)
	valueType := reflect.TypeOf(value)
	mapType := reflect.MapOf(keyType, valueType)
	mapSize := 0
	aMap := reflect.MakeMapWithSize(mapType, mapSize)
	return &Map{
		aMap,
	}
}

func (r *Map) Add(key, value interface{}) {
	r.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(value))
}
