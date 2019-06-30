package utils

import "reflect"

// StructToMap 将一个stuct转换成map
func StructToMap(i interface{}, fieldTag string, excludeTag ...string) map[string]interface{} {
	v_t := reflect.ValueOf(i)
	u_t := reflect.TypeOf(i)
	v_kind := v_t.Kind()
	var field_numbers int
	if v_kind == reflect.Struct {
		field_numbers = u_t.NumField()
	}
	if v_kind == reflect.Ptr {
		u_t = u_t.Elem()
		field_numbers = u_t.NumField()
		v_t = v_t.Elem()
	}
	field_map := map[string]interface{}{}
	for i := 0; i < field_numbers; i++ {
		struct_field := u_t.Field(i)
		if excludeTag != nil && struct_field.Tag.Get(excludeTag[0]) == "true" {
			continue
		}
		if struct_field.Anonymous && v_t.Field(i).Kind() == reflect.Struct {
			fv := v_t.Field(i)
			fi := fv.Interface()
			m1 := StructToMap(fi, fieldTag, excludeTag...)
			for k, v := range m1 {
				field_map[k] = v
			}
		} else {
			field_tag := struct_field.Tag.Get(fieldTag)
			field_val := v_t.Field(i)
			filed_interface := field_val.Interface()
			field_map[field_tag] = filed_interface
		}

	}
	return field_map
}
