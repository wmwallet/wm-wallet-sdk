package deposit

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	sdk "github.com/wmwallet/wm-wallet-sdk"
)

type (
	CreateOrderReq struct {
		SourceId   string          `json:"source_id"`
		ChainId    int             `json:"chain_id"`
		CoinId     int             `json:"coin_id"`
		FiatAmount decimal.Decimal `json:"fiat_amount"`
		Symbol     string          `json:"symbol"`
	}

	CreateOrderResp struct {
		SourceId     string          `json:"source_id"`
		ChainId      int             `json:"chain_id"`
		CoinId       int             `json:"coin_id"`
		FiatAmount   decimal.Decimal `json:"fiat_amount"`
		Symbol       string          `json:"symbol"`
		OrderId      string          `json:"order_id"`
		ExchangeRate decimal.Decimal `json:"exchange_rate"`
		Amount       decimal.Decimal `json:"amount"`
		Url          string          `json:"url"`
	}

	CancelOrderReq struct {
		SourceId string `json:"source_id"`
	}

	CancelOrderResp struct{}

	Resp[T CreateOrderResp | CancelOrderResp] struct {
		Data T      `json:"data"`
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
)

const (
	routeBrokerOrderCreate = "/v1/api/broker/order/create"
	routeBrokerOrderCancel = "/v1/api/broker/order/cancel"
)

type Deposit struct {
	w   *sdk.WmWalletClient
	url string
}

func NewDeposit(w *sdk.WmWalletClient, url string) *Deposit {
	return &Deposit{w: w, url: url}
}

func (d *Deposit) Create(ctx context.Context, req *CreateOrderReq) (*CreateOrderResp, error) {
	r, err := buildReq(ctx, req, d.url, routeBrokerOrderCreate)
	if err != nil {
		return nil, err
	}
	body, err := d.w.Post(ctx, r)
	if err != nil {
		return nil, err
	}
	var tmp string
	if err := json.Unmarshal(body, &tmp); err != nil {
		return nil, err
	}
	var resp = &Resp[CreateOrderResp]{}
	if err := json.Unmarshal([]byte(tmp), &resp); err != nil {
		return nil, err
	}
	if err = getErr(resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (d *Deposit) Cancel(ctx context.Context, req *CancelOrderReq) (*CancelOrderResp, error) {
	r, err := buildReq(ctx, req, d.url, routeBrokerOrderCancel)
	if err != nil {
		return nil, err
	}
	body, err := d.w.Post(ctx, r)
	if err != nil {
		return nil, err
	}

	var tmp string
	if err := json.Unmarshal(body, &tmp); err != nil {
		return nil, err
	}
	var resp = &Resp[CancelOrderResp]{}
	if err := json.Unmarshal([]byte(tmp), resp); err != nil {
		return nil, err
	}
	if err = getErr(resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func buildReq[T *CreateOrderReq | *CancelOrderReq](ctx context.Context, req T, baseUrl, router string) (*http.Request, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	r, err := http.NewRequest(http.MethodPost, baseUrl+router, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	r = r.WithContext(ctx)
	return r, nil
}

func getErr[T CreateOrderResp | CancelOrderResp](resp *Resp[T]) error {
	if resp.Code == 0 {
		return nil
	}
	return errors.New(resp.Msg)
}
