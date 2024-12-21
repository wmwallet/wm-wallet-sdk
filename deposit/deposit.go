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
		SourceId    string          `json:"source_id"`
		ChainId     int             `json:"chain_id"`
		CoinId      int             `json:"coin_id"`
		FiatAmount  decimal.Decimal `json:"fiat_amount"`
		Symbol      string          `json:"symbol"`
		CallbackUrl string          `json:"callback_url"`
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

	GetDetailRequest struct {
		SourceId string `json:"source_id"`
	}

	GetDetailResponse struct {
		SourceId     string          `json:"source_id"`
		ChainId      int             `json:"chain_id"`
		CoinId       int             `json:"coin_id"`
		Address      string          `json:"address"`
		Tag          string          `json:"tag"`
		Hash         string          `json:"hash"`
		FiatAmount   decimal.Decimal `json:"fiat_amount"`
		Symbol       string          `json:"symbol"`
		ExchangeRate decimal.Decimal `json:"exchange_rate"`
		Amount       decimal.Decimal `json:"amount"`
		OrderId      string          `json:"order_id"`
		Url          string          `json:"url"`
		Status       int8            `json:"status"`
		StatusDesc   string          `json:"status_desc"`
	}

	CancelOrderReq struct {
		SourceId string `json:"source_id"`
	}

	CancelOrderResp struct{}

	Resp[T CreateOrderResp | GetDetailResponse | CancelOrderResp] struct {
		Data T      `json:"data"`
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
)

const (
	routeBrokerOrderCreate = "/v1/api/broker/order/create"
	routeBrokerOrderDetail = "/v1/api/broker/order/detail"
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

func (d *Deposit) Detail(ctx context.Context, req *GetDetailRequest) (*GetDetailResponse, error) {
	r, err := buildReq(ctx, req, d.url, routeBrokerOrderDetail)
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
	var resp = &Resp[GetDetailResponse]{}
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

func buildReq[T *CreateOrderReq | *GetDetailRequest | *CancelOrderReq](ctx context.Context, req T, baseUrl, router string) (*http.Request, error) {
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

func getErr[T CreateOrderResp | GetDetailResponse | CancelOrderResp](resp *Resp[T]) error {
	if resp.Code == 0 {
		return nil
	}
	return errors.New(resp.Msg)
}
