package functions

import "strconv"

// 公共方法
func Str2Uint(idStr string) (id uint, err error) {
	var v int
	if v, err = strconv.Atoi(idStr); err != nil {
		return
	}
	id = uint(v)
	return
}
