package webhook

import (
	"github.com/shopspring/decimal"
)

type (
	DepositCallbackReq struct {
		SourceId     string          `json:"source_id"`
		OrderId      string          `json:"order_id"`
		FiatAmount   decimal.Decimal `json:"fiat_amount"`
		ExchangeRate decimal.Decimal `json:"exchange_rate"`
		Symbol       string          `json:"symbol"`
		Amount       decimal.Decimal `json:"amount"`
		Status       int8            `json:"status"`
	}
	DepositCallbackResp struct{}

	WithdrawCallbackReq struct {
		SourceId string          `json:"source_id"`
		OrderId  string          `json:"order_id"`
		ChainId  int             `json:"chain_id"`
		CoinId   int             `json:"coin_id"`
		Tag      string          `json:"tag"`
		Amount   decimal.Decimal `json:"amount"`
		Status   int8            `json:"status"`
	}

	WithdrawCallbackResp struct{}

	Resp[T DepositCallbackResp | WithdrawCallbackResp] struct {
		Data T      `json:"data"`
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
)
