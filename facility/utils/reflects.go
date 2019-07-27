package utils

import "reflect"

func InsertMap(i interface{}, fieldTag string) map[string]interface{} {
	return StructToMap(i, fieldTag, "omitinsert")
}

// StructToMap 将一个stuct转换成map
func StructToMap(i interface{}, fieldTag string, excludeTag ...string) map[string]interface{} {
	value := reflect.ValueOf(i)
	t := reflect.TypeOf(i)
	vkind := value.Kind()
	var fn int
	if vkind == reflect.Struct {
		fn = t.NumField()
	}
	if vkind == reflect.Ptr {
		t = t.Elem()
		fn = t.NumField()
		value = value.Elem()
	}
	fm := map[string]interface{}{}
	for i := 0; i < fn; i++ {
		sf := t.Field(i)
		if excludeTag != nil && sf.Tag.Get(excludeTag[0]) == "true" {
			continue
		}
		if sf.Anonymous && value.Field(i).Kind() == reflect.Struct {
			fv := value.Field(i)
			fi := fv.Interface()
			m1 := StructToMap(fi, fieldTag, excludeTag...)
			for k, v := range m1 {
				fm[k] = v
			}
		} else {
			ft := sf.Tag.Get(fieldTag)
			fv := value.Field(i)
			fif := fv.Interface()
			fm[ft] = fif
		}

	}
	return fm
}

func TagValues(i interface{}, tag string, excludeTags ...string) []string {
	value := reflect.ValueOf(i)
	t := reflect.TypeOf(i)
	vkind := value.Kind()
	var fn int
	if vkind == reflect.Struct {
		fn = t.NumField()
	}
	if vkind == reflect.Ptr {
		t = t.Elem()
		fn = t.NumField()
		value = value.Elem()
	}
	tagVals := []string{}
	for i := 0; i < fn; i++ {
		sf := t.Field(i)
		if excludeTags != nil && sf.Tag.Get(excludeTags[0]) == "true" {
			continue
		}
		ft := sf.Tag.Get(tag)
		tagVals = append(tagVals, ft)

	}
	return tagVals
}
