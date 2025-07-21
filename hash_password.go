package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "test123"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err \!= nil {
		panic(err)
	}
	fmt.Printf("%s\n", hash)
}
EOF < /dev/null