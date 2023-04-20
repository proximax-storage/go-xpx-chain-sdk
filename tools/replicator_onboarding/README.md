# Replicator Onboarding CLI tool

Allows to onboard a new replicator.

## Usage

### Flags

| Name          | Description                                        | Type   | Default               |
|:--------------|:---------------------------------------------------|:-------|:----------------------|
| `url`         | ProximaX Chain REST Url                            | string | http://127.0.0.1:3000 |
| `feeStrategy` | fee calculation strategy (`low`, `middle`, `high`) | string | `middle`              |
| `capacity`    | capacity of replicator (GB)                        | uint64 | -                     |
| `privateKey`  | Replicator private key                             | string | -                     |

### Example

```shell
./onboarding -url=http://127.0.0.1:3000 -feeStrategy=middle -capacity=1000 -privateKey=0000000000000000000000000000000000000000000000000000000000000000
```
