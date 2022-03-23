package guildsapi

import (
	"testing"

	cosmtypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestVerifySigAndDeriveAddress(t *testing.T) {
	cosmtypes.GetConfig().SetBech32PrefixForAccount("inj", "injpub")

	// inj14au322k9munkmx5wrchz9q30juf5wjgz2cfqku

	sigBase64 := "iYNArn1GOux3dNw/nefEZMW/83MRYMa7cdN/kVLNKKkeRZ5wn/caEOHzGfv6ktkSaNgf0rOa02R6gKb6wwEuKA=="
	pubKeyBase64 := "ApNNebT58zlZxO2yjHiRTJ7a7ufjIzeq5HhLrbmtg9Y/"
	txPayloadBase64 := "CoUBCoIBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEmIKKmluajE0YXUzMjJrOW11bmtteDV3cmNoejlxMzBqdWY1d2pnejJjZnFrdRIqaW5qMWhraGRhajJhMmNsbXE1anE2bXNwc2dncXMzMnZ5bnBrMjI4cTNyGggKA2luahIBMRJjCl8KVAotL2luamVjdGl2ZS5jcnlwdG8udjFiZXRhMS5ldGhzZWNwMjU2azEuUHViS2V5EiMKIQKTTXm0+fM5WcTtsox4kUye2u7n4yM3quR4S625rYPWPxIECgIIARiwFhIAGg1pbmplY3RpdmUtODg4IBc="

	addr, err := verifySigAndDeriveAddress(pubKeyBase64, sigBase64, txPayloadBase64)
	assert.NoError(t, err)
	assert.Equal(t, addr.String(), "inj14au322k9munkmx5wrchz9q30juf5wjgz2cfqku")
}
