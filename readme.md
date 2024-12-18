## wm-wallet-uml

### deposit
``` depoist
                      ,-.                                                                                                      
                      `-'                                                                                                      
                      /|\                                                                                                      
                       |             ,------.                                              ,--------.           ,--.           
                      / \            |Broker|                                              |WmWallet|           |DB|           
                     User            `---+--'                                              `----+---'           `-+'           
                       |                 |                                                      |                 |            
          ____________________________________________________________________________________________________________________ 
          ! CREATE&CANCEL DEPOSIT ORDER  /                                                      |                 |           !
          !_______________1 do deposit__/|                                                      |                 |           !
          !            |---------------->|                                                      |                 |           !
          !            |                 |                                                      |                 |           !
          !            |                 |2 do deposit(source_id,chain,coin,fiat_amount,symbol) |                 |           !
          !            |                 |----------------------------------------------------->|                 |           !
          !            |                 |                                                      |                 |           !
          !            |                 |3 ret order(order_id,chain,coin,exchange_rate,amount) |                 |           !
          !            |                 |<- - - - - - - - - - - - - - - - - - - - - - - - - - -|                 |           !
          !            |                 |                                                      |                 |           !
          !            |                 |                                                      |     4 save      |           !
          !            |                 |                                                      | - - - - - - - ->|           !
          !~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!
                       |                 |                                                      |                 |            
                       |                 |                                                      |                 |            
                 ________________________________________________________________________________________________________      
                 ! DEPOSIT  /            |                                                      |                 |      !     
                 !_________/             |                 5 webhook callback                   |                 |      !     
                 !     |                 |<- - - - - - - - - - - - - - - - - - - - - - - - - - -|                 |      !     
                 !     |                 |                                                      |                 |      !     
                 !     |                 |             6 check transaction exists               |                 |      !     
                 !     |                 |----------------------------------------------------->|                 |      !     
                 !     |                 |                                                      |                 |      !     
                 !     |                 |                                                      |                 |      !     
                 !     |   ____________________________________________________________________________________   |      !     
                 !     |   ! ALT  /      |                                                      |              !  |      !     
                 !     |   !~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!  |      !     
                 !     |   ! [transaction exists]                                               |              !  |      !     
                 !     |   !             |                 7 transaction info                   |              !  |      !     
                 !     |   !             |<- - - - - - - - - - - - - - - - - - - - - - - - - - -|              !  |      !     
                 !     |   !~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!  |      !     
                 !     |   ! [not exists]|                                                      |              !  |      !     
                 !     |   !             |                       8 null                         |              !  |      !     
                 !     |   !             |<- - - - - - - - - - - - - - - - - - - - - - - - - - -|              !  |      !     
                 !     |   !~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!  |      !     
                 !~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!     
                     User            ,---+--.                                              ,----+---.           ,-+.           
                      ,-.            |Broker|                                              |WmWallet|           |DB|           
                      `-'            `------'                                              `--------'           `--'           
                      /|\                                                                                                      
                       |                                                                                                       
                      / \                                                                                                      

```
### withdraw
```withdraw
                                                                                ,.-^^-._                
                      ,-.                                                      |-.____.-|               
                      `-'                                                      |        |               
                      /|\                                                      |        |               
                       |             ,--------.                                |        |       ,------.
                      / \            |WmWallet|                                '-.____.-'       |Wallet|
                     User            `----+---'                                   DB            `---+--'
                       |   1 withdraw     |                                        |                |   
                       |----------------->|                                        |                |   
                       |                  |                                        |                |   
                       |                  |2 check balance && freeze/deduct amount |                |   
                       |                  |--------------------------------------->|                |   
                       |                  |                                        |                |   
                       |                  |                                        |                |   
          _____________________________________________________________________________________     |   
          ! NOT ENOUGH  /                 |                                        |           !    |   
          !____________/                  |         3 not enough, cancel           |           !    |   
          !            |                  |<- - - - - - - - - - - - - - - - - - - -|           !    |   
          !            |                  |                                        |           !    |   
          !            |4 withdraw failed |                                        |           !    |   
          !            |<- - - - - - - - -|                                        |           !    |   
          !~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!    |   
                       |                  |                                        |                |   
                       |                  |                                        |                |   
                       |                  |                                        |                |   
                       |                  |                                        |                |   
                       |                  |                       5 withdraw       |                |   
                       |                  |-------------------------------------------------------->|   
                       |                  |                                        |                |   
                       |                  |                          6 ok          |                |   
                       |                  |<- - - - - - - - - - - - - - - - - - - - - - - - - - - - |   
                       |                  |                                        |                |   
                       |                  |               7 callback withdraw detail                |   
                       |                  |<- - - - - - - - - - - - - - - - - - - - - - - - - - - - |   
                       |                  |                                        |                |   
                       |                  |                                        |                |   
                       |   ____________________________________________________________________     |   
                       |   ! DETAIL SUCC  /                                        |           !    |   
                       |   !_____________/|            8 deduct amount             |           !    |   
                       |   !              |--------------------------------------->|           !    |   
                       |   !~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!    |   
                       |   ! [failed]     |                                        |           !    |   
                       |   !              |         9 unfeeze/add amount           |           !    |   
                       |   !              |--------------------------------------->|           !    |   
                       |   !~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~!    |   
                     User            ,----+---.                                   DB            ,---+--.
                      ,-.            |WmWallet|                                 ,.-^^-._        |Wallet|
                      `-'            `--------'                                |-.____.-|       `------'
                      /|\                                                      |        |               
                       |                                                       |        |               
                      / \                                                      |        |               
                                                                               '-.____.-'               

```
## Common

### Broker must in header.

### Chain Support:

| chain_id | chain_name |
|----------|------------|
| 1        | TRON       |
| 2        | TON        |

### Coin Support:

| coin_id | coin_name |
|---------|-----------|
| 1       | USDT      |
| 2       | TRX       |
| 3       | TON       |

### Symbol Support

| symbol   | 
|----------|
| USDT/TON |
| USDT/TRX |
| ......   |


### resp always like:

#### resp from gateway:
```json
{
  "code": 0,
  "data": {
    "code": 0,
    "msg": "",
    "data": {
      "source_id": "20241205114600032509",
      "chain_id": 2,
      "coin_id": 1,
      "address": "UQBetutxT6cLe54yHYRxKsUCW_bnvjYa_02GWFWQHaigBYWg",
      "tag": "13",
      "hash": "8To7dqNJ7kJPmTz4F9HSYpU/wvGZbx4q2RdkYutsm8E=",
      "fiat_amount": "2",
      "symbol": "USDT/CNY",
      "exchange_rate": "7.274951607768018",
      "amount": "0.2749159180474064",
      "order_id": "06ceb53730794f62be1a29d8aaf2efeb",
      "url": "http://43.156.157.230/wm/?broker=bacb849551df44e19608b46f9c4db031&no=06ceb53730794f62be1a29d8aaf2efeb&lang=zh",
      "status": 8,
      "status_desc": "订单成功"
    }
  }
}
```
#### api resp body in data:

``` json
{
    "code": 0,
    "msg": "",
    "data": {
        "source_id": "20241205114600032509",
        "chain_id": 2,
        "coin_id": 1,
        "address": "UQBetutxT6cLe54yHYRxKsUCW_bnvjYa_02GWFWQHaigBYWg",
        "tag": "13",
        "hash": "8To7dqNJ7kJPmTz4F9HSYpU/wvGZbx4q2RdkYutsm8E=",
        "fiat_amount": "2",
        "symbol": "USDT/CNY",
        "exchange_rate": "7.274951607768018",
        "amount": "0.2749159180474064",
        "order_id": "06ceb53730794f62be1a29d8aaf2efeb",
        "url": "http://43.156.157.230/wm/?broker=bacb849551df44e19608b46f9c4db031&no=06ceb53730794f62be1a29d8aaf2efeb&lang=zh",
        "status": 8,
        "status_desc": "订单成功"
    }
}
```

## API

### Deposit.Create

path: `/v1/api/broker/order/create`

req:

| name         | type            | comment | require |
|--------------|-----------------|---------|--------|
| source_id    | string          | uniq id | y      |
| chain_id     | int             |         | y      |
| coin_id      | int             |         | y      |
| fiat_amount  | decimal(40,18)  |         | y      |
| symbol       | string          |         | y      |

resp:

| name          | type           | comment       |
|---------------|----------------|---------------|
| source_id     | string         | broker uid    |
| chain_id      | int            |               |
| coin_id       | int            |               |
| fiat_amount   | decimal(40,18) |               |
| exchange_rate | decimal(40,18) |               |
| amount        | decimal(40,18) |               |
| order_id      | string         | wm wallet uid |
| url           | string         | a html view   |

### Deposit.Detail

path: `/v1/api/broker/order/detail`

req:

| name         | type           | comment | require |
|--------------|----------------|---------|---------|
| source_id    | string         | uniq id | y       |

resp:

| name          | type           | comment     |
|---------------|----------------|-------------|
| source_id     | string         |             |
| chain_id      | int            |             |
| coin_id       | int            |             |
| address       | string         |             |
| tag           | string         |             |
| hash          | string         |             |
| fiat_amount   | decimal(40,18) |             |
| symbol        | string         |             |
| exchange_rate | decimal(40,18) |             |
| amount        | decimal(40,18) |             |
| order_id      | string         |             |
| url           | string         | a html view |
| status        | int8           |             |
| status_desc   | string         |             |


| status | status_desc     |
|--------|-----------------|
| 0      | Order Created   |
| 8      | Order Success   |
| 12     | User Cancel     |
| 16     | Expire Cancel   |
| 20     | Less Amount     |

### Deposit.Cancel

path:    `/v1/api/broker/order/cancel`

req:

| name       | type    | comment | require |
|------------|---------|---------|---------|
| source_id  | string  | uniq id | y       |

resp:

| name  | type | comment |
|-------|------|---------|
| null  |      |         |


### Deposit Callback Struct

| name          | type           | comment     |
|---------------|----------------|-------------|
| source_id     | string         |             |
| chain_id      | int            |             |
| coin_id       | int            |             |
| address       | string         |             |
| tag           | string         |             |
| hash          | string         |             |
| fiat_amount   | decimal(40,18) |             |
| symbol        | string         | USDT/TON    |
| exchange_rate | decimal(40,18) |             |
| amount        | decimal(40,18) | 0.123456    |
| order_id      | string         |             |
| status        | int8           |             |
| status_desc   | string         |             |

| status | status_desc      |
|--------|------------------|
| 0      | Withdraw Doing   |
| 4      | Withdraw Success |
| 8      | Withdraw Fail    |

### Withdraw.Create

path:    `/v1/api/broker/order/withdraw`

req:

| name      | type           | comment | require |
|-----------|----------------|---------|---------|
| source_id | string         | uniq id | y       |
| chain_id  | int            |         | y       |
| coin_id   | int            |         | y       |
| address   | string         |         | y       |
| tag       | string         |         | y       |
| amount    | decimal(40,18) |         | y       |

resp:

| name            | type            | comment                                |
|-----------------|-----------------|----------------------------------------|
| source_id       | string          | uniq id                                |
| order_id        | string          |                                        |


### Withdraw.Detail

path:    `/v1/api/broker/order/withdraw-detail`

req:

| name      | type           | comment | require |
|-----------|----------------|---------|---------|
| source_id | string         | uniq id | y       |

resp:

| name          | type            | comment |
|---------------|-----------------|---------|
| source_id     | string          | uniq id |
| chain_id      | int             |         |
| coin_id       | int             |         |
| address       | string          |         |
| tag           | string          |         |
| hash          | string          |         |
| fiat_amount   | decimal(40,18)  |         |
| symbol        | string          |         |
| exchange_rate | decimal(40,18)  |         |
| amount        | decimal(40,18)  |         |
| service_fee   | decimal(40,18)  |         |
| order_id      | string          |         |
| status        | int8            |         |
| status_desc   | string          |         |


### Withdraw Callback Struct

| name        | type           | comment |
|-------------|----------------|---------|
| source_id   | string         |         |
| order_id    | int            |         |
| chain_id    | int            |         |
| coin_id     | int            |         |
| tag         | string         |         |
| amount      | decimal(40,18) |         |
| service_fee | decimal(40,18) |         |
| status      | int8           |         |
| status_desc | string         |         |


## sdk example

```go
package thirdparty

import (
	"context"
	sdk "github.com/wmwallet/wm-wallet-sdk"
	"github.com/wmwallet/wm-wallet-sdk/deposit"
	"github.com/wmwallet/wm-wallet-sdk/withdraw"
)

const CustomerName = "test"

const Url = "http://wmwallet.pro"

type WmWalletSDK struct {
	wmWalletSDK *sdk.WmWalletClient
}

func NewWmWalletSDK() *WmWalletSDK {
	ops := []sdk.Option{
		sdk.WithSecretPath("wallet/public_key.pem"),
		sdk.WithCertPath("wallet/ca.crt", "wallet/client.crt", "wallet/client.key"),
		sdk.WithCustomer(CustomerName),
	}
	wmWalletSDK, err := sdk.Init(ops...)
	if err != nil {
		panic(err)
	}
	return &WmWalletSDK{
		wmWalletSDK: wmWalletSDK,
	}
}

func (wws *WmWalletSDK) DepositOrderCreate(req *deposit.CreateOrderReq) (resp *deposit.CreateOrderResp, err error) {
	d := deposit.NewDeposit(wws.wmWalletSDK, Url)
	resp, err = d.Create(context.Background(), req)
	if err != nil {
		return
	}
	return
}

func (wws *WmWalletSDK) DepositOrderDetail(req *deposit.GetDetailRequest) (resp *deposit.GetDetailResponse, err error) {
	d := deposit.NewDeposit(wws.wmWalletSDK, Url)
	resp, err = d.Detail(context.Background(), req)
	if err != nil {
		return
	}
	return
}

func (wws *WmWalletSDK) DepositOrderCancel(req *deposit.CancelOrderReq) (resp *deposit.CancelOrderResp, err error) {
	d := deposit.NewDeposit(wws.wmWalletSDK, Url)
	resp, err = d.Cancel(context.Background(), req)
	if err != nil {
		return
	}
	return
}

func (wws *WmWalletSDK) WithdrawOrderCreate(req *withdraw.CreateOrderReq) (resp *withdraw.CreateOrderResp, err error) {
	w := withdraw.NewWithdraw(wws.wmWalletSDK, Url)
	resp, err = w.Create(context.Background(), req)
	if err != nil {
		return
	}
	return
}

func (wws *WmWalletSDK) WithdrawOrderDetail(req *withdraw.GetDetailRequest) (resp *withdraw.GetDetailResponse, err error) {
	w := withdraw.NewWithdraw(wws.wmWalletSDK, Url)
	resp, err = w.Detail(context.Background(), req)
	if err != nil {
		return
	}
	return
}

```
