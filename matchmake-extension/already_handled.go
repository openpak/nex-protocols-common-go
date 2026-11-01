package matchmake_extension

import "reflect"

// alreadyHandled reports whether the protocol struct behind iface already has a handler assigned
// to the named callback field. The OpenPak additions to this package (the Find* methods) must not
// replace a title's own handler that was assigned before NewCommonProtocol ran.
func alreadyHandled(iface any, field string) bool {
	v := reflect.ValueOf(iface)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return false
	}
	f := v.Elem().FieldByName(field)
	return f.IsValid() && f.Kind() == reflect.Func && !f.IsNil()
}
