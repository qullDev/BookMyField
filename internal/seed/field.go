package seed

import (
	"log"

	"github.com/google/uuid"
	"github.com/qullDev/BookMyField/internal/config"
	"github.com/qullDev/BookMyField/internal/models"
)

func SeedFields() {
	fields := []models.Field{
		{ID: uuid.New(), Name: "Lapangan Futsal A", Location: "Jakarta", Price: 200000},
		{ID: uuid.New(), Name: "Lapangan Basket B", Location: "Bandung", Price: 150000},
		{ID: uuid.New(), Name: "Lapangan Badminton C", Location: "Surabaya", Price: 100000},
	}

	for _, f := range fields {
		if err := config.DB.FirstOrCreate(&models.Field{}, models.Field{Name: f.Name, Location: f.Location}).Error; err != nil {
			log.Printf("❌ Gagal seed field %s: %v", f.Name, err)
		}
	}
	log.Println("✅ Seed data fields berhasil")
}
