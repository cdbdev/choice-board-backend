package main

import (
	"corner-backend/internal/pkg/api"
)

func main() {
	err := api.Start()
	if err != nil {
		panic(err)
	}
}