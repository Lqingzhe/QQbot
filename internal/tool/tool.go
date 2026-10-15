package tool

import "strings"

func GetRoutKet(Types ...string) string {
	if len(Types) == 0 {
		return ""
	}
	if len(Types) <= 3 {
		return strings.Join(Types, " ")
	} else {
		return strings.Join(Types[:3], " ")
	}
} //一次传入onebot11协议中三级type，如果少于三级可以不传
