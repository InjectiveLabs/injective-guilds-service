package config

// TODO: Use env var
var GrantRequirements = []string{
	"/injective.exchange.v1beta1.MsgCreateSpotLimitOrder",
	"/injective.exchange.v1beta1.MsgCreateSpotMarketOrder",
	"/injective.exchange.v1beta1.MsgCancelSpotOrder",
	"/injective.exchange.v1beta1.MsgBatchCancelSpotOrders",
	"/injective.exchange.v1beta1.MsgDeposit",
	"/injective.exchange.v1beta1.MsgWithdraw",
	"/injective.exchange.v1beta1.MsgCreateDerivativeLimitOrder",
	"/injective.exchange.v1beta1.MsgCreateDerivativeMarketOrder",
	"/injective.exchange.v1beta1.MsgCancelDerivativeOrder",
	"/injective.exchange.v1beta1.MsgBatchUpdateOrders",
	"/injective.exchange.v1beta1.MsgBatchCancelDerivativeOrders",
}
