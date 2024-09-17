package helpers

import (
	"encoding/hex"

	"github.com/liobrdev/simplepasswords_vaults/utils"
)

var HexHash = hex.EncodeToString(utils.HashToken(VALID_EMAIL + VALID_PW))
