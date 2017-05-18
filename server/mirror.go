package mirror

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Mirror is a record of a reflect variable
type Mirror struct {
	Key        string   `json:"key"`
	Name       string   `json:"name"`
	Value      string   `json:"value"`
	StaticType string   `json:"staticType"`
	Children   []Mirror `json:"children"`
}

// MirrorMap hold all of the reflect variables: a tree
var mirrorMap = make(map[string]Mirror)

// valueMap hold all of the raw variables {key:reflect.Value}: a flat hash table
var valueMap = make(map[string]reflect.Value)

// PutInMirror reflect the specific variable and put in MirrorMap
// 使用动态反射，所以现在只需要反射一层
func PutInMirror(name string, value interface{}) {
	if _, ok := valueMap[name]; ok {
		fmt.Printf("WARN: %s is already exists\n", name)
		return
	}

	var v = reflect.ValueOf(value)
	if v.IsValid() {
		valueMap[name] = v
	} else {
		valueMap[name] = reflect.ValueOf(name + " is nil,can't be reflected!")
	}
	//mirrorMap[name] = startReflect(v, name, "", 1)
}

func startReflect(v reflect.Value, selfName string, parentName string, level int) Mirror {
	var m = Mirror{}

	if !v.IsValid() {
		m.Key = createKey(selfName, parentName)
		m.StaticType = "string"
		m.Value = selfName + " is nil,can't be reflected!"
		m.Children = nil
		m.Name = selfName + ": " + m.Value + " [" + m.StaticType + "]"
		return m
	}

	switch v.Kind() {
	case reflect.Struct:
		m = reflectStruct(v, selfName, parentName, level)
	case reflect.Map:
		m = reflectMap(v, selfName, parentName, level)
	case reflect.Slice, reflect.Array:
		m = reflectSlice(v, selfName, parentName, level)
	case reflect.Ptr:
		m = startReflect(v.Elem(), selfName, parentName, level)
	case reflect.Interface:
		m = startReflect(reflect.ValueOf(v), selfName, parentName, level)
	default:
		m = reflectAtom(v, selfName, parentName, level)
	}
	return m
}

func createKey(selfName string, parentName string) string {
	if parentName == "" {
		return selfName
	}
	return parentName + "." + selfName
}

func reflectByKey(key string) Mirror {
	var m = Mirror{}
	var v = valueMap[key]
	var lastDot = strings.LastIndex(key, ".")
	var selfName = key
	var parentName = ""
	if lastDot > -1 {
		selfName = key[lastDot+1:]
		parentName = key[0:lastDot]
	}
	m = startReflect(v, selfName, parentName, 1)
	return m
}

func reflectSlice(value reflect.Value, selfName string, parentName string, level int) Mirror {
	var m = Mirror{}
	m.Key = createKey(selfName, parentName)
	m.StaticType = value.Type().Name()
	m.Value = getReflectValue(value)
	m.Name = selfName + ": " + m.Value + " [" + m.StaticType + "]"
	valueMap[m.Key] = value
	if level <= 0 {
		m.Children = []Mirror{}
	} else {
		for i := 0; i < value.Len(); i++ {
			var name = strconv.Itoa(i)
			var v = value.Index(i)
			m.Children = append(m.Children, startReflect(v, name, m.Key, level-1))
		}
	}
	return m
}

func reflectMap(value reflect.Value, selfName string, parentName string, level int) Mirror {
	var m = Mirror{}
	m.Key = createKey(selfName, parentName)
	m.StaticType = value.Type().Name()
	m.Value = getReflectValue(value)
	m.Name = selfName + ": " + m.Value + " [" + m.StaticType + "]"
	valueMap[m.Key] = value
	if level <= 0 {
		m.Children = []Mirror{}
	} else {
		var mapKeys = value.MapKeys()
		for i := 0; i < len(mapKeys); i++ {
			var k = mapKeys[i]
			var v = value.MapIndex(k)
			var kAsName = getReflectValue(k)
			m.Children = append(m.Children, startReflect(v, kAsName, m.Key, level-1))
		}
	}
	return m
}

func reflectStruct(value reflect.Value, selfName string, parentName string, level int) Mirror {
	var m = Mirror{}
	m.Key = createKey(selfName, parentName)
	m.StaticType = value.Type().Name()
	m.Value = getReflectValue(value)
	m.Name = selfName + ": " + m.Value + " [" + m.StaticType + "]"
	valueMap[m.Key] = value
	if level <= 0 {
		m.Children = []Mirror{}
	} else {
		for i := 0; i < value.NumField(); i++ {
			var cv = value.Field(i)
			var cvName = value.Type().Field(i).Name
			m.Children = append(m.Children, startReflect(cv, cvName, m.Key, level-1))
		}
	}
	return m
}

func reflectAtom(value reflect.Value, selfName string, parentName string, level int) Mirror {
	var m = Mirror{}
	m.Key = createKey(selfName, parentName)
	m.StaticType = value.Type().Name()
	m.Value = getReflectValue(value)
	m.Children = nil
	m.Name = selfName + ": " + m.Value + " [" + m.StaticType + "]"
	valueMap[m.Key] = value
	return m
}

//inspired from https://github.com/golang/go/blob/049b89dc6f6b6f1001672dd5456197b74a97cbec/src/fmt/print.go#L845
func getReflectValue(value reflect.Value) string {
	var r string

	/*if !value.IsValid() {
		r = "nil"
		return r
	}*/

	switch value.Kind() {
	case reflect.Bool:
		r = strconv.FormatBool(value.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		r = strconv.FormatInt(value.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		r = strconv.FormatUint(value.Uint(), 10)
	case reflect.Float32:
		r = strconv.FormatFloat(value.Float(), 'E', -1, 32)
	case reflect.Float64:
		r = strconv.FormatFloat(value.Float(), 'E', -1, 64)
	case reflect.Complex64:
		r = "value-Complex64"
	case reflect.Complex128:
		r = "value-Complex128"
	case reflect.String:
		r = value.String()
	case reflect.Map:
		r = "value-Map"
	case reflect.Array:
		r = "value-Array"
	case reflect.Chan:
		r = "value-Chan"
	case reflect.Func:
		r = "value-Func"
	case reflect.Interface:
		r = "value-Interface"
	case reflect.Ptr:
		r = "value-Ptr"
	case reflect.Slice:
		r = "value-Slice"
	case reflect.Struct:
		r = "value-Struct"
	case reflect.UnsafePointer:
		r = "value-UnsafePointer"
	default:
		r = "God knows"
	}
	return r
}
