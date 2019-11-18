package reflect_util

import (
	"reflect"
)

// StructToMap tranfer a struct to map with fields tagged with given tag
func StructToMap(i interface{}, fieldTag string, omitTags ...string) map[string]interface{} {
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
		if omitTags != nil {
			ok := false
			for _, omit := range omitTags {
				if _, found := sf.Tag.Lookup(omit); found {
					ok = true
					break
				}
			}
			if ok {
				continue
			}
		}
		if sf.Anonymous && value.Field(i).Kind() == reflect.Struct {
			fv := value.Field(i)
			fi := fv.Interface()
			m1 := StructToMap(fi, fieldTag, omitTags...)
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

// TagValues will return tag values for a given tag
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
