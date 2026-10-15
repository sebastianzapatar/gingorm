package main

import (
	"log"

	"ejemplo/crud-go-sqlite/internal/database"
	"ejemplo/crud-go-sqlite/internal/handlers"
	"ejemplo/crud-go-sqlite/internal/routes"
)

func main() {
	db, err := database.Conectar("productos.db")
	if err != nil {
		log.Fatal(err)
	}
	r := routes.Nuevo(handlers.NuevoProductos(db))
	log.Println("API en http://localhost:8080")
	log.Fatal(r.Run("127.0.0.1:8080"))
}
