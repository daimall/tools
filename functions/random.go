package functions

import (
	"math/rand"
	"strconv"
	"time"
)

// 生成指定位数的随机数字
func CreateRandomNumber(count int) string {
	var numbers = []int{1, 2, 3, 4, 5, 7, 8, 9}
	var container string
	length := len(numbers)
	for i := 1; i <= count; i++ {
		source := rand.NewSource(time.Now().Unix() + int64(rand.Intn(i+50000))) // 使用当前时间作为随机种子
		randomGenerator := rand.New(source)
		random := randomGenerator.Intn(length)
		// rand.Seed(time.Now().UnixNano())
		// random := rand.Intn(length)
		container += strconv.Itoa(numbers[random])
	}
	return container
}

// 生成指定长度的随机字符串
func CreateRandomString(count int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	length := len(letters)
	b := make([]rune, count)
	for i := 0; i < count; i++ {
		source := rand.NewSource(time.Now().Unix() + int64(rand.Intn(i+50000))) // 使用当前时间作为随机种子
		randomGenerator := rand.New(source)
		randomInt := randomGenerator.Intn(length)
		b[i] = letters[randomInt]
	}
	return string(b)
}

// 判断一个字符串是纯数字
func IsNumeric(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}
