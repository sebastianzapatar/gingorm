package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"ejemplo/crud-go-sqlite/internal/models"
)

// Conectar abre (o crea) el archivo SQLite y crea las tablas que falten.
func Conectar(ruta string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(ruta), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err = db.AutoMigrate(&models.Producto{}); err != nil {
		return nil, err
	}
	return db, nil
}
