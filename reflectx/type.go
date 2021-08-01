package reflectx

import (
	"bytes"
	"reflect"
)

var TypeBytes = reflect.TypeOf([]byte(""))

func IsBytes(tpe reflect.Type) bool {
	return tpe.Kind() != reflect.String && tpe.ConvertibleTo(TypeBytes)
}

func FullTypeName(rtype reflect.Type) string {
	buf := bytes.NewBuffer(nil)

	for rtype.Kind() == reflect.Ptr {
		buf.WriteByte('*')
		rtype = rtype.Elem()
	}

	if name := rtype.Name(); name != "" {
		if pkgPath := rtype.PkgPath(); pkgPath != "" {
			buf.WriteString(pkgPath)
			buf.WriteRune('.')
		}
		buf.WriteString(name)
		return buf.String()
	}

	buf.WriteString(rtype.String())
	return buf.String()
}

func Deref(tpe reflect.Type) reflect.Type {
	if tpe.Kind() == reflect.Ptr {
		return Deref(tpe.Elem())
	}
	return tpe
}
