package models

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUserCreation(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		t.Fatal(err)
	}

	user := &User{
		Email:        "unit@test.com",
		PasswordHash: "hashedpass",
	}

	if err := db.Create(user).Error; err != nil {
		t.Errorf("failed to create user: %v", err)
	}

	if user.ID == [16]byte{} {
		t.Errorf("expected UUID to be set, got zero value")
	}
}
