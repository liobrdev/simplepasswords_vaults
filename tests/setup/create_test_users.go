package setup

import (
	"testing"

	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/tests/helpers"
)

func createTestUsers(t *testing.T, db *gorm.DB) (users []models.User) {
	users = []models.User{
		{
			Slug: helpers.NewSlug(t),
		},
		{
			Slug: helpers.NewSlug(t),
		},
	}

	if result := db.Create(&users); result.Error != nil {
		t.Fatalf("Create test users failed: %s", result.Error.Error())
	}

	return
}
