package main

import (
	"fmt"
	"lab-ap/config"
	"lab-ap/database"
	"lab-ap/internal/entity"
)

func main() {
	cfg := config.Load()
	db, err := database.Connect(cfg.DSN(), false)
	if err != nil {
		panic(err)
	}
	var configs []entity.Konfigurasi
	db.Find(&configs)
	for _, c := range configs {
		fmt.Printf("Key: %s, Value: %s\n", c.Key, c.Value)
	}
}
