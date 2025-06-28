package goreloaded

import (
	"strconv"
)

func hexbin(s string ,base string) string {
	res := ""
	if base == "(hex)"{
		num, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		return  s
	}
	res = strconv.Itoa(int(num))
	return res

	}else if base == "(bin)"{
		num, err := strconv.ParseInt(s, 2, 64)
	if err != nil {
		return s
	}
	res = strconv.Itoa(int(num))
	
	}
	return res
}
