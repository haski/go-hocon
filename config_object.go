package hocon

import (
	"strings"
)

type Content map[string]*Value

type ConfigObject struct {
	content *Content
}

func (c *ConfigObject) setString(path string, value string) {
	prefix, key := pathPrefix(splitPath(path))
	resultKey := setObjectKey(prefix, c)
	(*resultKey.content)[key] = MakeStringValue(value)
}

func (c *ConfigObject) setInt(path string, value string) {
	prefix, key := pathPrefix(splitPath(path))
	resultKey := setObjectKey(prefix, c)
	(*resultKey.content)[key] = MakeNumericValue(value)
}

func (c *ConfigObject) setObject(path string, value *ConfigObject) {
	prefix, key := pathPrefix(splitPath(path))
	resultKey := setObjectKey(prefix, c)
	(*resultKey.content)[key] = &Value{
		Type:     ObjectType,
		RefValue: value,
	}
}

func (c *ConfigObject) setArray(path string, value *ConfigArray) {
	prefix, key := pathPrefix(splitPath(path))
	resultKey := setObjectKey(prefix, c)
	(*resultKey.content)[key] = &Value{
		Type:     ArrayType,
		RefValue: value,
	}
}

func (c *ConfigObject) setReference(path string, value string) {
}

func NewConfigObject() *ConfigObject {
	m := make(Content)
	co := ConfigObject{
		content: &m,
	}
	return &co
}

func setObjectKey(keys []string, obj *ConfigObject) *ConfigObject {
	for _, key := range keys {
		if v, exists := (*obj.content)[key]; exists {
			switch v.Type {
			case ObjectType:
				obj = v.RefValue.(*ConfigObject)
			default:
				panic("Wrong path")
			}
			continue
		}

		newObject := NewConfigObject()
		(*obj.content)[key] = &Value{
			Type:     ObjectType,
			RefValue: newObject,
		}
		obj = newObject
	}
	return obj
}

func traversePath(o *ConfigObject, path string) (*ConfigObject, string) {
	obj := o
	paths := strings.Split(path, ".")
	for _, p := range paths[:len(paths)-1] {
		if d := (*obj.content)[p]; d == nil {
			return nil, ""
		} else {
			switch d.Type {
			case ObjectType:
				obj = d.RefValue.(*ConfigObject)
			default:
				return nil, ""
			}
		}
	}
	return obj, paths[len(paths)-1]
}

func (o *ConfigObject) GetString(path string) (res string) {
	if obj, key := traversePath(o, path); obj != nil {
		if v, ok := (*obj.content)[key]; ok {
			res = v.RefValue.(string)
		}
	}
	return res
}

func (o *ConfigObject) GetInt(path string) (res int) {
	if obj, key := traversePath(o, path); obj != nil {
		if v, ok := (*obj.content)[key]; ok {
			res = v.RefValue.(int)
		}
	}
	return res
}

func (o *ConfigObject) GetObject(path string) (res *ConfigObject) {
	if obj, key := traversePath(o, path); obj != nil {
		if v, ok := (*obj.content)[key]; ok {
			res = v.RefValue.(*ConfigObject)
		}
	}
	return res
}

func (o *ConfigObject) GetArray(path string) (res *ConfigArray) {
	if obj, key := traversePath(o, path); obj != nil {
		if v, ok := (*obj.content)[key]; ok {
			res = v.RefValue.(*ConfigArray)
		}
	}
	return res
}

func (o *ConfigObject) GetKeys() []string {
	res := make([]string, len(*o.content))
	i := 0
	for k, _ := range *o.content {
		res[i] = k
		i++
	}
	return res
}
