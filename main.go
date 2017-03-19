package main

import (
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	log.SetFlags(0)
	log.Println("Starting...")
	log.Fatal(StartServer())
}
