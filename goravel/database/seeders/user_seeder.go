package seeders
	

import (
	"github.com/goravel/framework/facades"
	"goravel/app/models"
	"golang.org/x/crypto/bcrypt" 

)


type UserSeeder struct {
}

// Signature The name and signature of the seeder.
func (s *UserSeeder) Signature() string {
	return "UserSeeder"
}

// Run executes the seeder logic.
func (s *UserSeeder) Run() error {
	for _, userData := range []struct{
    	Name, Email, Password string
	}{
		{"Manager", "manager@gmail.com", "qwer1234"},
		{"Admin", "admin@gmail.com", "qwer1234"},
		{"Developer", "developer@gmail.com", "qwer1234"},
		{"Client", "client@gmail.com", "qwer1234"},

	} {
		password, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		user := models.User{
			Name:     userData.Name,
			Email:    userData.Email,
			Password: string(password),
		}
		if err := facades.Orm().Query().Create(&user); err != nil {
			return err
		}
	}

	return nil
}
