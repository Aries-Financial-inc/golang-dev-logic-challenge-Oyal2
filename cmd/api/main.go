package main

import (
	"fmt"

	"github.com/Aries-Financial-inc/golang-dev-logic-challenge-Oyal2/internal/server"
)

func main() {
	server := server.NewServer()

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
