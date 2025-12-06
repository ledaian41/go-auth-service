package user

import (
	"errors"
	"go-auth-service/internal/shared/dto"
	"go-auth-service/internal/shared/utils"
	"log"

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

func (s *UserService) FindUserByUsername(username, siteId string) (*shared_dto.UserDTO, error) {
	filteredUsers := shared_utils.Filter(UserList, func(user User) bool {
		return user.Site == siteId
	})
	for _, user := range filteredUsers {
		if user.Username == username {
			userCopy := user.ToDTO()
			return &userCopy, nil
		}
	}
	return nil, errors.New("user not found")
}

func (s *UserService) CreateNewUser(user *shared_dto.UserDTO) (*shared_dto.UserDTO, error) {
	if IsUsernameExist(user.Username) {
		return nil, errors.New("username exist")
	}

	newUser := User{
		ID:           shared_utils.RandomID(),
		Username:     user.Username,
		Password:     user.Password,
		PhoneNumber:  user.PhoneNumber,
		Email:        user.Email,
		Role:         user.Role,
		Site:         user.Site,
		TokenVersion: user.TokenVersion,
	}
	UserList = append(UserList, newUser)
	return user, nil
}

func IsUsernameExist(username string) bool {
	for _, user := range UserList {
		if user.Username == username {
			return true
		}
	}
	return false
}

func (s *UserService) FindUsersBySite(siteId string) *[]shared_dto.UserDTO {
	filteredUsers := shared_utils.Filter(UserList, func(user User) bool {
		return user.Site == siteId
	})
	result := shared_utils.Map(filteredUsers, func(user User) shared_dto.UserDTO {
		return shared_dto.UserDTO{
			Username:    user.Username,
			Name:        user.Name,
			PhoneNumber: user.PhoneNumber,
			Email:       user.Email,
			Role:        user.Role,
		}
	})
	return &result
}
