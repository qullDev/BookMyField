package seed

import (
	"log"

	"github.com/google/uuid"
	"github.com/qullDev/BookMyField/internal/config"
	"github.com/qullDev/BookMyField/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func SeedAdminUser() {
	// cek apakah admin sudah ada
	var user models.User
	if err := config.DB.First(&user, "role = ?", "admin").Error; err == nil {
		log.Println("⚠️ Admin user sudah ada, skip seeding")
		return
	}

	// hash password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	admin := models.User{
		ID:       uuid.New(),
		Name:     "Super Admin",
		Email:    "admin@admin.com",
		Password: string(hashedPassword),
		Role:     "admin",
	}

	if err := config.DB.Create(&admin).Error; err != nil {
		log.Printf("❌ Gagal seed admin user: %v", err)
		return
	}

	log.Println("✅ Seed admin user berhasil (email: admin@admin.com, password: password123)")
}
