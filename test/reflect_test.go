package test

import (
	"fmt"
	"gotrue/facility/utils"
	"reflect"
	"testing"
)

type User1 struct {
	Name    string `ttt:"name" exclude:"true"`
	TestAge int    `ttt:"test_age"`
	Addr
}

type Addr struct {
	Hello string `ttt:"address"`
}

func TestReflect(t *testing.T) {
	var u interface{} = &User1{
		Name:    "lxc",
		TestAge: 15,
	}
	v_t := reflect.ValueOf(u)
	u_t := reflect.TypeOf(u)
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
	t.Logf("type of u is %v, value of u is %v, kind of u is %v,field numbers is %d",
		u_t, v_t, v_kind, field_numbers)

	field_map := map[string]interface{}{}
	for i := 0; i < field_numbers; i++ {
		struct_field := u_t.Field(i)
		if struct_field.Tag.Get("exclude") == "true" {
			t.Logf("field %s exclueded to map", struct_field.Name)
			continue
		}
		tag := struct_field.Tag.Get("db")
		field_kind := struct_field.Type.Kind()
		field_val := v_t.Field(i)
		filed_interface := field_val.Interface()
		t.Logf("struct_field of index %d is %v, kind is %v, tag is %s ",
			i, struct_field, field_kind, tag)
		field_map[tag] = filed_interface
	}
	t.Logf("The struct User1's map is %v", field_map)
}

func TestMyStruct2Map(t *testing.T) {
	ad := Addr{
		Hello: "address",
	}
	var u interface{} = &User1{
		Name:    "lxc",
		TestAge: 15,
		Addr:    ad,
	}
	fmt.Printf("transformed map is %v \n", utils.StructToMap(u, "ttt", "exclude"))
}
