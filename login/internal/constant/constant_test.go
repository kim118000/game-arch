package constant

import (
	"reflect"
	"testing"
)

func TestName(t *testing.T) {
	var a *int = new(int)
	*a = 1

	b:=reflect.TypeOf(a)

	t.Logf("%v, %v", b, b.Elem())

	val := reflect.ValueOf(a)
	t.Logf("%T, %v \n",val, val)
	t.Logf("%T, %v \n",val.Interface(), val.Interface())
	t.Logf("%T, %v \n",val.Elem(), val.Elem())
	t.Logf("%T, %v \n",val.Elem().Interface(), val.Elem().Interface())

	val.Elem().SetInt(10)

	t.Log(*a)
}
