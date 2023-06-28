package pkg

import "strconv"

// StringToInt 将字符串转为为int整数
func StringToInt(str string) (int, error) {
	intValue, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return intValue, nil
}
