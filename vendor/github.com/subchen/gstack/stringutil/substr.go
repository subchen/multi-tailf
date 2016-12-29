package stringutil

import (
	"strings"
)

func SubstrBefore(str string, find string) string {
	if len(find) == 0 {
		return ""
	}
	pos := strings.Index(str, find)
	if pos == -1 {
		return str
	}
	return str[:pos]
}

func SubstrAfter(str string, find string) string {
	if len(find) == 0 {
		return str
	}
	pos := strings.Index(str, find)
	if pos == -1 {
		return ""
	}
	return str[pos+len(find):]
}

func SubstrBeforeLast(str string, find string) string {
	if len(find) == 0 {
		return str
	}
	pos := strings.LastIndex(str, find)
	if pos == -1 {
		return str
	}
	return str[:pos]
}

func SubstrAfterLast(str string, find string) string {
	if len(find) == 0 {
		return ""
	}
	pos := strings.LastIndex(str, find)
	if pos == -1 {
		return ""
	}
	return str[pos+len(find):]
}
