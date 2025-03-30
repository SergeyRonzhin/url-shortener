package main

import (
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func main() {

	var users []User
	client := resty.New()

	_, err := client.R().
		SetResult(&users).
		Get("https://jsonplaceholder.typicode.com/users")

	if err != nil {
		panic(err)
	}

	names := make([]string, len(users))

	for _, u := range users {
		names = append(names, u.Username)
	}

	fmt.Println(strings.Join(names, " "))
}
