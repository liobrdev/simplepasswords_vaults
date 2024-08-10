package setup

import (
	"testing"

	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_vaults/models"
)

func populateTestDB(t *testing.T, db *gorm.DB) (
	users []models.User,
	vaults []models.Vault,
	entries []models.Entry,
	secrets []models.Secret,
) {
	users = createTestUsers(t, db)
	vaults = createTestVaults(&users, t, db)
	entries = createTestEntries(&users, &vaults, t, db)
	secrets = createTestSecrets(&users, &vaults, &entries, t, db)
	return users, vaults, entries, secrets
}
