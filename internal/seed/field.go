package seed

import (
	"log"

	"github.com/qullDev/BookMyField/internal/config"
	"github.com/qullDev/BookMyField/internal/models"
)

func SeedFields() {
	fields := []models.Field{
		{Name: "Lapangan Futsal A", Location: "Jakarta", Price: 200000},
		{Name: "Lapangan Basket B", Location: "Bandung", Price: 150000},
		{Name: "Lapangan Badminton C", Location: "Surabaya", Price: 100000},
	}

	for _, f := range fields {
		// Check if field already exists
		var existingField models.Field
		if err := config.DB.Where("name = ? AND location = ?", f.Name, f.Location).First(&existingField).Error; err != nil {
			// Field doesn't exist, create new one
			if err := config.DB.Create(&f).Error; err != nil {
				log.Printf("❌ Gagal seed field %s: %v", f.Name, err)
			} else {
				log.Printf("✅ Field %s berhasil dibuat", f.Name)
			}
		} else {
			log.Printf("⚠️ Field %s sudah ada, skip seeding", f.Name)
		}
	}
	log.Println("✅ Seed data fields berhasil")
}
