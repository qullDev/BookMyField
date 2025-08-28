package seed

import (
	"log"

	"github.com/google/uuid"
	"github.com/qullDev/BookMyField/internal/config"
	"github.com/qullDev/BookMyField/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func SeedRegularUser() {
	var user models.User
	if err := config.DB.First(&user, "role = ?", "user").Error; err == nil {
		log.Println("⚠️ Regular user sudah ada, skip seeding")
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)

	regularUser := models.User{
		ID:       uuid.New(),
		Name:     "Regular User",
		Email:    "user@bookmyfield.com",
		Password: string(hashedPassword),
		Role:     "user",
	}

	if err := config.DB.Create(&regularUser).Error; err != nil {
		log.Printf("❌ Gagal seed regular user: %v", err)
		return
	}

	log.Println("✅ Seed regular user berhasil (email: user@bookmyfield.com, password: user123)")
}
