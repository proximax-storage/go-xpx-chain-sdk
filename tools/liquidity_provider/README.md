# Liquidity Provider CLI tool

Allows to create and change liquidity provider in the ProximaX BC.

## Usage

#### Common Flags

| Name           | Description                                                                                       | Type   | Default               |
|:---------------|:--------------------------------------------------------------------------------------------------|:-------|:----------------------|
| `sender`       | private key of transaction sender (**required**, should be one of cosigners if owner is multisig) | string | -                     |
| `owner`        | public key of liquidity provider owner (**optional**, set if owner is a multisig)                 | string | -                     |
| `url`          | ProximaX Chain REST Url                                                                           | string | http://127.0.0.1:3000 |
| `feeStrategy`  | fee calculation strategy (`low`, `middle`, `high`)                                                | string | `middle`              |

```shell
./lp <command> []<flag>
```

### Create Command

#### Flags

| Name             | Description                                   | Type   | Default |
|:-----------------|:----------------------------------------------|:-------|:--------|
| `mosaic`         | provider mosaic name (storage, streaming, sc) | string | -       |
| `initial`        | amount of initial mosaics minting             | uint64 | 100000  |
| `deposit`        | amount of currency deposit                    | uint64 | 100000  |
| `slashingPeriod` | slashing period                               | uint32 | 500     |
| `slashingAcc`    | slashing account public key                   | string | -       |
| `ws`             | window size                                   | unit16 | 5       |
| `a`              | alpha                                         | uint32 | 500     |
| `b`              | beta                                          | uint32 | 500     |

#### Example

```shell
./lp -url=http://127.0.0.1:3000 \
    -feeStrategy=middle \
    -sender=0000000000000000000000000000000000000000000000000000000000000000 \
    -owner=0000000000000000000000000000000000000000000000000000000000000000 \
    -mosaic=storage \
    -initial=100000 \
    -deposit=1000000 \
    -slashingPeriod=500 \
    -ws=5 \
    -a=500 \
    -b=500 \
    -slashingAcc=0000000000000000000000000000000000000000000000000000000000000000 \
    create
```

### Change Command

#### Flags

| Name                      | Description                                   | Type   | Default |
|:--------------------------|:----------------------------------------------|:-------|:--------|
| `mosaic`                  | provider mosaic name (storage, streaming, sc) | string | -       |
| `currencyBalanceIncrease` | currency balance increase                     | bool   | false   |
| `currencyBalanceChange`   | currency balance change                       | uint64 | 0       |
| `mosaicBalanceIncrease`   | mosaic balance increase                       | bool   | false   |
| `mosaicBalanceChange`     | mosaic balance change                         | uint64 | 0       |

#### Example

```shell
./lp -url=http://127.0.0.1:3000 \
    -feeStrategy=middle \
    -sender=0000000000000000000000000000000000000000000000000000000000000000 \
    -owner=0000000000000000000000000000000000000000000000000000000000000000 \
    -mosaic=storage \
    -currencyBalanceIncrease \
    -currencyBalanceChange=100 \
    -mosaicBalanceIncrease \
    -mosaicBalanceChange=200 \
    change
```
