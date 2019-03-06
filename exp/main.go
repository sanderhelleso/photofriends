package main

import (
	"fmt"
	"../../photofriends/rand"
)

func main() {
	fmt.Println(rand.String(10))
	fmt.Println(rand.RememberToken())
}