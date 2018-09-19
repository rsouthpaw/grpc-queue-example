package main

import "gdma_handler/handler"

const (
	PORT = ":50051"
)

func main() {
	handler.StartServer(PORT)
}
