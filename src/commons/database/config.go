package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectPostgresql(env string) {
	_, currentFile, _, _ := runtime.Caller(0)

	projectRoot := filepath.Join(
		filepath.Dir(currentFile),
		"..",
		"..",
		"..",
	)

	envPath := filepath.Join(projectRoot, env)

	_, err := os.Stat(envPath)
	if err == nil {
		log.Println("Memuat env:", envPath)

		if err := godotenv.Load(envPath); err != nil {
			log.Println("Gagal memuat .env:", err)
		}
	} else {
		log.Println("File env tidak ditemukan, menggunakan environment variable")
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, password, dbname, port, sslmode)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("gagal memuat database")
	}

	log.Println("Database berhasil tersambung")
}
