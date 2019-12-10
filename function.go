package goutill

import "reflect"

type function struct{}

var Function = function{}

// interface, function name, arguments 를 인자로 해당 인터페이스의 function 을 수행, 인자값을 리턴한다.
func (function) Invoke(v interface{}, name string, args ...interface{}) []reflect.Value {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}

	return reflect.ValueOf(v).MethodByName(name).Call(inputs)
}
