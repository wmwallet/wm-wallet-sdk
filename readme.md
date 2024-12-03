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

### Chain support:

| chain_id | chain_name |
|----------|------------|
| 1        | TRON       |
| 2        | TON        |

### Coin support:

| coin_id | coin_name |
|---------|-----------|
| 1       | USDT      |
| 2       | TRX       |
| 3       | TON       |


### resp always like:

#### resp from gateway:
```json
{
    "code": 0,
    "data": {
        "code": 0,
        "msg": "",
        "data": {
            "chain": "TON",
            "hash": "og6hrZpiFxjsKsfRAr+CIQd",
            "address": "EQAkhkg79yqAqbG67Ch5m8j",
            "tag": "2",
            "coin": "USDT",
            "amount": "0.222222",
            "blockNo": "",
            "txId": "da2c87300d993cc924b9085af5f5183a"
        }
    }
}
```
#### api resp body in data:

``` json
{
    "msg": "",
    "code": 0,
    "data": {
        "chain": "TON",
        "hash": "og6hrZpiFxjsKsfRAr+CIQdZp",
        "address": "EQAkhkg79yqAqbG67Ch5m8j",
        "tag": "2",
        "coin": "USDT",
        "amount": "0.222222",
        "blockNo": "",
        "txId": "da2c87300d993cc924b9085af5f5183a"
    }
}
```

## API

### Deposit.Create

path: `/v1/api/broker/order/create`

req:

| name         | type           | comment | require |
|--------------|----------------|---------|---------|
| source_id    | string         | uniq id | y       |
| chain_id     | int            |         | y       |
| coin_id      | int            |         | y       |
| fiat_amount  | decimal(40,18) |         | y       |

resp:

| name          | type           | comment     |
|---------------|----------------|-------------|
| source_id     | string         |             |
| chain_id      | int            |             |
| coin_id       | int            |             |
| fiat_amount   | decimal(40,18) |             |
| order_id      | string         |             |
| exchange_rate | decimal(40,18) |             |
| amount        | decimal(40,18) |             |
| url           | string         | a html view |

### Deposit.Cancel

path:    `/v1/api/broker/order/cancel`

req:

| name       | type    | comment | require |
|------------|---------|---------|---------|
| source_id  | string  | uniq id | y       |

resp:

| name | type | comment          |
|------|------|------------------|
| null |      |                  |


### Deposit Callback Struct

| name          | type           | comment            |
|---------------|----------------|--------------------|
| source_id     | string         |                    |
| order_id      | string         |                    |
| fiat_amount   | decimal(40,18) |                    |
| exchange_rate | decimal(40,18) |                    |
| symbol        | string         | coin, USDT/TON/... |
| amount        | decimal(40,18) | 0.123456           |
| status        | int8           |                    |

### Withdraw.Create

path:    `/v1/api/broker/order/withdraw`

req:

| name      | type       | comment | require |
|-----------|------------|---------|---------|
| source_id | string     | uniq id | y       |
| chain_id  | string     |         | y       |
| coin_id   | string     |         | y       |
| address   | string     |         | y       |
| tag       | string     |         | y       |
| amount    | string     |         | y       |

resp:

| name            | type            | comment                                |
|-----------------|-----------------|----------------------------------------|
| source_id       | string          | uniq id                                |
| order_id        | string          |                                        |


### Withdraw Callback Struct

| name        | type           | comment |
|-------------|----------------|---------|
| source_id   | string         |         |
| order_id    | int            |         |
| chain_id    | int            |         |
| coin_id     | int            |         |
| tag         | string         |         |
| amount      | decimal(40,18) |         |
| status      | int8           |         |


## sdk example

```go
// prod
ops := []sdk.Option{
    sdk.WithSecretPath("config/public_key.pem"),
    sdk.WithCertPath(("config/test_ca.crt"), ("config/test_client.crt"), ("config/test_client.key")),
    sdk.WithCustomer("test"),
}
w, err := sdk.Init(ops...)
if err != nil {
t.Fatal(err)
}

d := NewDeposit(w, URL)
resp, err := d.GetNewAddress(context.Background(), &GetNewAddrReq{
Network:   "TON",
RequestId: "12345",
})

```

```go
// test
ops := []sdk.Option{
sdk.WithCustomer("a"),
sdk.WithTest(true),
}
w, err := sdk.Init(ops...)
if err != nil {
t.Fatal(err)
}

d := NewDeposit(URL)
resp, err := d.GetNewAddress(context.Background(), &GetNewAddrReq{
Network:   "TON",
RequestId: "12345",
})
if err != nil {
t.Fatal(err)
}
```
