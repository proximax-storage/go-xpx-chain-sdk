# Liquidity Provider CLI tool

Allows to create and change liquidity provider in the ProximaX BC.

## Usage

#### Common Flags

| Name          | Description                                          | Type   | Default               |
|:--------------|:-----------------------------------------------------|:-------|:----------------------|
| `sender`      | private account of transaction sender (**required**) | string | -                     |
| `url`         | ProximaX Chain REST Url                              | string | http://127.0.0.1:3000 |
| `feeStrategy` | fee calculation strategy (`low`, `middle`, `high`)   | string | `middle`              |

```shell
./lp <command> []<flag>
```

### Create Command

#### Flags

| Name             | Description                                         | Type   |
|:-----------------|:----------------------------------------------------|:-------|
| `mosaic`         | provider mosaic id, e.g. 0x6C5D687508AC9D75         | string |
| `initial`        | amount of initial mosaics minting                   | uint64 |
| `deposit`        | amount of currency deposit                          | uint64 |
| `slashingPeriod` | slashing period                                     | uint32 |
| `slashingAcc`    | slashing account public key                         | string |
| `ws`             | window size                                         | unit16 |
| `a`              | alpha                                               | uint32 |
| `b`              | beta                                                | uint32 |

#### Example

```shell
./lp -url=http://127.0.0.1:3000 \
    -feeStrategy=middle \
    -sender=0000000000000000000000000000000000000000000000000000000000000000 \
    -mosaic=26514E2A1EF33824 \
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

| Name                      | Description                                 | Type   |
|:--------------------------|:--------------------------------------------|:-------|
| `mosaic`                  | provider mosaic id, e.g. 0x6C5D687508AC9D75 | string |
| `currencyBalanceIncrease` | currency balance increase                   | bool   |
| `currencyBalanceChange`   | currency balance change                     | uint64 |
| `mosaicBalanceIncrease`   | mosaic balance increase                     | bool   |
| `mosaicBalanceChange`     | mosaic balance change                       | uint64 |

#### Example

```shell
./lp -url=http://127.0.0.1:3000 \
    -feeStrategy=middle \
    -sender=7AA907C3D80B3815BE4B4E1470DEEE8BB83BFEB330B9A82197603D09BA947230 \
    -mosaic=26514E2A1EF33824 \
    -currencyBalanceIncrease \
    -currencyBalanceChange=100 \
    -mosaicBalanceIncrease \
    -mosaicBalanceChange=200 \
    change
```
