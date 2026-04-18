package main

import (
	"adult-short-videos/common/utils"
	"fmt"
)

func main() {
	hash, _ := utils.PasswordHash("password123")
	fmt.Println("加密后：" + hash)
	fmt.Println("验证结果：", utils.PasswordVerify(hash, "password123"))

	token, _ := utils.GenerateToken(1, "testuser", "my-secret", 3600)
	fmt.Println("生成的 Token：" + token)

	claims, _ := utils.ParseToken(token, "my-secret")
	fmt.Println("用户ID：", claims.UserID)
}
