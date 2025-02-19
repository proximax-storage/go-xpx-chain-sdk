# Replicator Offboarding CLI tool

Allows to offboard a replicator.

## Usage

### Flags

| Name                   | Description                                        | Type   | Default               |
|:-----------------------|:---------------------------------------------------|:-------|:----------------------|
| `url`                  | ProximaX Chain REST Url                            | string | http://127.0.0.1:3000 |
| `feeStrategy`          | fee calculation strategy (`low`, `middle`, `high`) | string | `middle`              |
| `replicatorPrivateKey` | Replicator private key                             | string | -                     |
| `drivePublicKey`       | Drive public key                                   | string | -                     |

### Example

```shell
./offboarding -url=http://127.0.0.1:3000 -feeStrategy=middle -replicatorPrivateKey=0000000000000000000000000000000000000000000000000000000000000000 -drivePublicKey=0000000000000000000000000000000000000000000000000000000000000000
```