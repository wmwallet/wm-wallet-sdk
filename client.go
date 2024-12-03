package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const (
	Wsign   = "w-sign"
	Wbroker = "w-broker"
	Wts     = "w-ts"
	Wnonce  = "w-nonce"
	Wsecret = "w-secret"

	wtest = "w-test"
)

type GWResp struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

func (w *WmWalletClient) postWithEncrypt(ctx context.Context, req *http.Request) ([]byte, error) {
	cli := w.client
	encrypt := w.encrypt
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(req.Body)
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, errors.WithMessage(err, "read request body")
	}

	secret := generateRandomString(letters, 16)
	cipher, err := encrypt.AESEncryptECB([]byte(secret), body)
	if err != nil {
		return nil, errors.WithMessage(err, "encrypt body err")
	}

	cipherS, err := w.encrypt.Encrypt([]byte(secret))
	if err != nil {
		return nil, errors.WithMessage(err, "encrypt secret err")
	}

	req.Header.Set(Wsecret, cipherS)
	req.Header.Set(Wnonce, generateRandomString(digits, 6))
	req.Header.Set(Wts, strconv.Itoa(int(time.Now().UnixMilli())))
	req.Header.Set(Wbroker, w.customer)
	sign(req, body)

	req.Body = io.NopCloser(bytes.NewReader(cipher))
	req.ContentLength = int64(len(cipher))
	req.Header.Set("Content-Length", strconv.Itoa(len(cipher)))
	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return nil, errors.WithMessage(err, "http post request err")
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithMessage(err, "read response body")
	}
	gwResp := &GWResp{}
	err = json.Unmarshal(b, gwResp)
	if err != nil {
		return nil, errors.WithMessage(err, "json unmarshal err")
	}
	if gwResp.Code != 0 {
		return nil, errors.New(gwResp.Msg)
	}
	return gwResp.Data, nil
}

func (w *WmWalletClient) postWithoutEncrypt(ctx context.Context, req *http.Request) ([]byte, error) {
	cli := w.client

	req.Header.Set(Wbroker, w.customer)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(wtest, "1")

	resp, err := cli.Do(req)
	if err != nil {
		return nil, errors.WithMessage(err, "http post request err")
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithMessage(err, "read response body")
	}
	gwResp := &GWResp{}
	err = json.Unmarshal(b, gwResp)
	if err != nil {
		return nil, errors.WithMessage(err, "json unmarshal err")
	}
	if gwResp.Code != 0 {
		return nil, errors.New(gwResp.Msg)
	}
	return gwResp.Data, nil
}

func (w *WmWalletClient) Post(ctx context.Context, req *http.Request) ([]byte, error) {
	if w.isTest {
		return w.postWithoutEncrypt(ctx, req)
	}
	return w.postWithEncrypt(ctx, req)
}
