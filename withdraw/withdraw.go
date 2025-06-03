package withdraw

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
		Address     string          `json:"address"`
		Tag         string          `json:"tag"`
		Amount      decimal.Decimal `json:"amount"`
		CallbackUrl string          `json:"callback_url"`
	}

	CreateOrderResp struct {
		SourceId string `json:"source_id"`
		OrderId  string `json:"order_id"`
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
		GasFee       decimal.Decimal `json:"gas_fee"`
		OrderId      string          `json:"order_id"`
		Status       int8            `json:"status"`
		StatusDesc   string          `json:"status_desc"`
	}

	Resp[T CreateOrderResp | GetDetailResponse] struct {
		Data T      `json:"data"`
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
)

const (
	routeBrokerWithdrawOrderCreate = "/v1/api/broker/order/withdraw"
	routeBrokerWithdrawOrderDetail = "/v1/api/broker/order/withdraw-detail"
)

type Withdraw struct {
	w   *sdk.WmWalletClient
	url string
}

func NewWithdraw(w *sdk.WmWalletClient, url string) *Withdraw {
	return &Withdraw{w: w, url: url}
}

func (d *Withdraw) Create(ctx context.Context, req *CreateOrderReq) (*CreateOrderResp, error) {
	r, err := buildReq(ctx, req, d.url, routeBrokerWithdrawOrderCreate)
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
	if err := json.Unmarshal([]byte(tmp), resp); err != nil {
		return nil, err
	}
	if err = getErr(resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (d *Withdraw) Detail(ctx context.Context, req *GetDetailRequest) (*GetDetailResponse, error) {
	r, err := buildReq(ctx, req, d.url, routeBrokerWithdrawOrderDetail)
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
	if err := json.Unmarshal([]byte(tmp), resp); err != nil {
		return nil, err
	}
	if err = getErr(resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func buildReq[T *CreateOrderReq | *GetDetailRequest](ctx context.Context, req T, baseUrl, router string) (*http.Request, error) {
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

func getErr[T CreateOrderResp | GetDetailResponse](resp *Resp[T]) error {
	if resp.Code == 0 {
		return nil
	}
	return errors.New(resp.Msg)
}
