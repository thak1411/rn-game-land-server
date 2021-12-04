package main

import (
	"fmt"
	"regexp"
)

// /^[ㄱ-ㅎ|가-힣|ㅏ-ㅣ|a-z|A-Z|0-9|]{2,6}$/, /^[a-z|A-Z|0-9|]{6,12}$/, /^[a-z|A-Z|0-9|]{6,12}$/

func main() {
	r, _ := regexp.Compile(`^[ㄱ-ㅎ|가-힣|ㅏ-ㅣ|a-z|A-Z|0-9|]{2,6}$`)
	fmt.Println(r.MatchString("안녕하세요"))
}
