package user

import (
	"encoding/json"
	"errors"
	"go-auth-service/internal/shared"
	"io/ioutil"
	"log"
	"os"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) MigrateDatabase() {
	err := s.db.AutoMigrate(&User{})
	if err != nil {
		log.Printf("❌ Failed at AutoMigrate: %v", err)
	}
	log.Println("🎉 NodeType - Database migrate successfully")
}

func (s *UserService) FindUserByUsername(username, siteId string) (*shared.UserDTO, error) {
	var user User
	if err := s.db.Where("username = ? AND site = ?", username, siteId).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}
	userDto := user.ToDTO()
	return &userDto, nil
}

func (s *UserService) CreateNewUser(user *shared.UserDTO) (*shared.UserDTO, error) {
	var count int64
	s.db.Model(&User{}).Where("username = ?", user.Username).Count(&count)
	if count > 0 {
		return nil, errors.New("username exist")
	}

	newUser := User{
		ID:           shared.RandomID(),
		Username:     user.Username,
		Password:     user.Password,
		PhoneNumber:  user.PhoneNumber,
		Email:        user.Email,
		Role:         user.Role,
		Site:         user.Site,
		TokenVersion: user.TokenVersion,
	}

	if err := s.db.Create(&newUser).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) SeedUsers(filePath string) error {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // File not found is fine, just skip seeding
	}

	// Read file
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	var users []User
	if err := json.Unmarshal(fileData, &users); err != nil {
		return err
	}

	for _, user := range users {
		// Check if user exists
		var count int64
		s.db.Model(&User{}).Where("username = ? AND site = ?", user.Username, user.Site).Count(&count)
		if count > 0 {
			continue
		}

		user.ID = shared.RandomID()
		// Password in JSON is already hashed (bcrypt), so we don't need to hash it again if it starts with $
		// But in this specific userData.json, they seem already hashed.
		// If they were plain text, we would need to hash them.
		// Based on the file content provided: "$2a$10$..." so it is hashed.

		if err := s.db.Create(&user).Error; err != nil {
			log.Printf("Failed to seed user %s: %v", user.Username, err)
		}
	}
	log.Println("🎉 Seeding completed")
	return nil
}

func (s *UserService) FindUsersBySite(siteId string) *[]shared.UserDTO {
	var users []User
	s.db.Where("site = ?", siteId).Find(&users)

	var result []shared.UserDTO
	for _, user := range users {
		result = append(result, user.ToDTO())
	}
	return &result
}
