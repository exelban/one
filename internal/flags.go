package internal

import (
	"strings"
)

func StringFlag(args *[]string, flags ...string) string {
	for i, arg := range *args {
		for _, f := range flags {
			if strings.HasPrefix(arg, f+"=") {
				*args = remove(*args, i)
				return strings.TrimPrefix(arg, f+"=")
			}
			if strings.HasPrefix(arg, f) && len(*args) > i+1 {
				value := (*args)[i+1]
				*args = remove(*args, i)
				*args = remove(*args, i)
				return value
			}
		}
	}
	return ""
}
func BoolFlag(args *[]string, flags ...string) bool {
	for i, arg := range *args {
		for _, f := range flags {
			if arg == f {
				*args = remove(*args, i)
				return true
			}
		}
	}
	return false
}

func remove(args []string, i int) []string {
	return append(args[:i], args[i+1:]...)
}
