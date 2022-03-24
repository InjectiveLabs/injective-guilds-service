package guildsapi

import (
	"fmt"
	"testing"

	cosmtypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestVerifySigAndExtratInfo(t *testing.T) {
	cosmtypes.GetConfig().SetBech32PrefixForAccount("inj", "injpub")

	// inj14au322k9munkmx5wrchz9q30juf5wjgz2cfqku
	sigBase64 := "/57Tewk8mKsou/USmw2bkC84M4hJzqR991I2xID3COVNQwRcpCDsw3bShJibSNivF0toWSpV2nDVdKWBcwKI1w=="
	pubKeyBase64 := "ApNNebT58zlZxO2yjHiRTJ7a7ufjIzeq5HhLrbmtg9Y/"
	txPayloadBase64 := "eyJhY3Rpb24iOiAiZW50ZXItZ3VpbGQiLCAiZXhwaXJlZF9hdCI6IDE3MzQ0ODQ5OTN9"

	addr, msgInfo, err := verifySigAndExtractInfo(pubKeyBase64, sigBase64, txPayloadBase64)
	assert.NoError(t, err)
	assert.Equal(t, addr.String(), "inj14au322k9munkmx5wrchz9q30juf5wjgz2cfqku")
	assert.Equal(t, msgInfo.Action, "enter-guild")
	fmt.Println("info:", msgInfo)
}
