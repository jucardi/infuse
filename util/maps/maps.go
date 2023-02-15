package maps

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/jucardi/infuse/util/reflectx"
)

var (
	regex = regexp.MustCompile(`.*\[\d*\]$`)
)

// Contains indicates if the given map contains an entry by the given key
func Contains(c map[string]interface{}, key string) bool {
	if !strings.Contains(key, ".") {
		if _, ok := c[key]; ok {
			return ok
		}

		return false
	}

	if _, err := GetValue(c, key); err != nil {
		return false
	}

	return true
}

// GetValue If a map represents a JSON with nested objects. GetValue retrieves the value by the given path. Eg. 'info.database.port'
func GetValue(data map[string]interface{}, key string) (interface{}, error) {
	split := strings.Split(key, ".")
	v := reflect.ValueOf(data)
	for i, s := range split {
		if s == "" {
			continue
		}

		isArray := regex.MatchString(s)
		index := 0

		if isArray {
			split := strings.Split(s, "[")
			s = split[0]
			index, _ = strconv.Atoi(split[1][:len(split[1])-1])
		}

		current := v.MapIndex(reflect.ValueOf(s))

		if !current.IsValid() {
			return nil, fmt.Errorf("unable to get value by the key '%s'. The value for '%s' is not present", key, s)
		}

		if reflectx.IsNil(current) {
			if i < len(split)-1 {
				return nil, fmt.Errorf("unable to get value by the key '%s'. The value for '%s' is null", key, s)
			} else {
				return nil, nil
			}
		}

		if isArray {
			for current = current.Elem(); current.IsValid() && current.Kind() != reflect.Slice && current.Kind() != reflect.Array; {
			}
			if current.Len() <= index {
				return nil, fmt.Errorf("failed to retrieve value at key '%s'. Index out of range for field '%s' (index: %d | length: %d)", key, s, index, current.Len())
			}
			current = current.Index(index)
		}

		if i < len(split)-1 {
			if v.Kind() != reflect.Map {
				return nil, fmt.Errorf("unable to get value by path: '%s' | The piece '%s' does not represent an object", key, s)
			}

			// m, err := ConvertMap(current.Interface())
			m, err := reflectx.GetNonPointerValue(current)
			if err != nil {
				return nil, err
			}

			v = m
		} else {
			v = current
		}
	}

	if !v.IsValid() {
		return nil, fmt.Errorf("value by the key '%s' is not present", key)
	}

	if reflectx.IsNil(v) {
		return nil, nil
	}

	return v.Interface(), nil
}

// GetOrDefault gets the value by the given key, if the value is not present, or an error occurs while retrieving the value, returns what was specified as `defaultVal`
func GetOrDefault(data map[string]interface{}, key string, defaultVal interface{}) interface{} {
	if v, _ := GetValue(data, key); v == nil {
		return defaultVal
	} else {
		return v
	}
}

// SetValue if a map represents a JSON with nested objects. SetValue assigns a value to the given path. Eg. 'info.database.port'
// 'makeEmpty' indicates that if a piece of the path is missing (Eg. 'info.database' is nil) an empty object should be created to continue the assignment.
func SetValue(data map[string]interface{}, key string, value interface{}, makeEmpty bool) error {
	split := strings.Split(key, ".")
	v := reflect.ValueOf(data)
	for i, s := range split {
		if s == "" {
			continue
		}
		if i == len(split)-1 {
			v.SetMapIndex(reflect.ValueOf(s), reflect.ValueOf(value))
		} else {
			if v.Kind() != reflect.Map {
				return fmt.Errorf("unable to get value by path: '%s' | The piece '%s' does not represent an object", key, s)
			}

			val := v.MapIndex(reflect.ValueOf(s))

			if !val.IsValid() {
				if makeEmpty {
					val = reflect.ValueOf(make(map[string]interface{}))
					v.SetMapIndex(reflect.ValueOf(s), val)
				} else {
					return fmt.Errorf("unable to get value by path: '%s' | The piece '%s' does not represent an object", key, s)
				}
			}
			v = reflect.ValueOf(val.Interface().(map[string]interface{}))
		}
	}
	return nil
}

func ConvertMap(val interface{}) (map[string]interface{}, error) {
	if m, ok := val.(map[string]interface{}); ok {
		for k, v := range m {
			if mapValue, ok := v.(map[interface{}]interface{}); ok {
				if newVal, err := ConvertMap(mapValue); err != nil {
					return nil, fmt.Errorf("failed to convert '%s', %s ", k, err.Error())
				} else {
					m[k] = newVal
				}
			}
		}
		return m, nil
	}

	if m, ok := val.(map[interface{}]interface{}); ok {
		ret := map[string]interface{}{}
		for k, v := range m {
			key, ok := k.(string)
			if !ok {
				return nil, fmt.Errorf("all keys must be strings when mapping to a struct, detected key: '%+v'", k)
			}
			if mapValue, ok := v.(map[interface{}]interface{}); ok {
				if newVal, err := ConvertMap(mapValue); err != nil {
					return nil, fmt.Errorf("failed to convert '%s', %s ", key, err.Error())
				} else {
					ret[key] = newVal
				}
			} else {
				ret[key] = v
			}
		}
		return ret, nil
	}

	return nil, fmt.Errorf("unexpected object type: %+v", val)
}

func StringMapEqual(m1 map[string]string, m2 map[string]string) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k1, v1 := range m1 {
		if v2, ok := m2[k1]; ok {
			if v1 != v2 {
				return false
			}
		} else {
			return false
		}
	}
	return true
}
