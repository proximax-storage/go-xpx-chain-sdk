# Replicator Onboarding CLI tool

Allows to onboard a new replicator.

## Usage

### Flags

| Name          | Description                                        | Type   | Default               |
|:--------------|:---------------------------------------------------|:-------|:----------------------|
| `url`         | ProximaX Chain REST Url                            | string | http://127.0.0.1:3000 |
| `feeStrategy` | fee calculation strategy (`low`, `middle`, `high`) | string | `middle`              |
| `capacity`    | capacity of replicator (MB)                        | uint64 | -                     |
| `sender`      | Sender private key                                 | string | -                     |
| `receiver`    | Receiver public key                                | string | -                     |
| `mosaic`      | Id of transfer mosaic                              | string | -                     |
| `amount`      | Amount of transfer mosaic                          | string | -                     |

### Example

```shell
./transfer -url=http://127.0.0.1:3000 \
 -feeStrategy=middle \
 -sender=0000000000000000000000000000000000000000000000000000000000000000 \
 -receiver=0000000000000000000000000000000000000000000000000000000000000000 \
 -mosaic=6C5D687508AC9D75 \
 -amount=10000000
```
