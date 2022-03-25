package config

type DenomConfig struct {
	CoinID         string
	DisplayDecimal int
}

var StableCoinDenoms = map[string]bool{
	"peggy0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48":                      true,
	"peggy0xdAC17F958D2ee523a2206206994597C13D831ec7":                      true,
	"ibc/B448C0CA358B958301D328CCDC5D5AD642FC30A6D3AE106FF721DB315F3DDE5C": true,
}

// First we can hardcod it
// TODO: Use env var+json
var DenomConfigs = map[string]*DenomConfig{
	"peggy0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2": {
		CoinID:         "weth",
		DisplayDecimal: 3,
	},
	"peggy0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48": {
		CoinID:         "usd-coin",
		DisplayDecimal: 2,
	},
	"inj": {
		CoinID:         "injective-protocol",
		DisplayDecimal: 2,
	},
	"peggy0xdAC17F958D2ee523a2206206994597C13D831ec7": {
		CoinID:         "tether",
		DisplayDecimal: 1,
	},
	"peggy0x514910771AF9Ca656af840dff83E8264EcF986CA": {
		CoinID:         "chainlink",
		DisplayDecimal: 2,
	},
	"peggy0x7Fc66500c84A76Ad7e9c93437bFc5Ac33E2DDaE9": {
		CoinID:         "aave",
		DisplayDecimal: 3,
	},
	"peggy0x7D1AfA7B718fb893dB30A3aBc0Cfc608AaCfeBB0": {
		CoinID:         "matic-network",
		DisplayDecimal: 2,
	},
	"peggy0x1f9840a85d5aF5bf1D1762F925BDADdC4201F984": {
		CoinID:         "uniswap",
		DisplayDecimal: 2,
	},
	"peggy0x6B3595068778DD592e39A122f4f5a5cF09C90fE2": {
		CoinID:         "sushi",
		DisplayDecimal: 2,
	},
	"peggy0xc944E90C64B2c07662A292be6244BDf05Cda44a7": {
		CoinID:         "the-graph",
		DisplayDecimal: 2,
	},
	"peggy0xC011a73ee8576Fb46F5E1c5751cA3B9Fe0af2a6F": {
		CoinID:         "havven",
		DisplayDecimal: 2,
	},
	"peggy0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599": {
		CoinID:         "wrapped-bitcoin",
		DisplayDecimal: 4,
	},
	"peggy0xBB0E17EF65F82Ab018d8EDd776e8DD940327B28b": {
		CoinID:         "axie-infinity",
		DisplayDecimal: 2,
	},
	"ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9": {
		CoinID:         "cosmos",
		DisplayDecimal: 2,
	},
	"peggy0xAaEf88cEa01475125522e117BFe45cF32044E238": {
		CoinID:         "guildfi",
		DisplayDecimal: 2,
	},
	"ibc/B448C0CA358B958301D328CCDC5D5AD642FC30A6D3AE106FF721DB315F3DDE5C": {
		CoinID:         "terrausd",
		DisplayDecimal: 1,
	},
	"ibc/B8AF5D92165F35AB31F3FC7C7B444B9D240760FA5D406C49D24862BD0284E395": {
		CoinID:         "terra-luna",
		DisplayDecimal: 3,
	},
	"ibc/E7807A46C0B7B44B350DA58F51F278881B863EC4DCA94635DAB39E52C30766CB": {
		CoinID:         "chihuahua-token",
		DisplayDecimal: 1,
	},
}
