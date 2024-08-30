# Replicators cleanup CLI tool

Allows to remove replicators that are not bound with nodes.

## Usage

### Flags

| Name                | Description                                           | Type   | Default               |
|:--------------------|:------------------------------------------------------|:-------|:----------------------|
| `url`               | ProximaX Chain REST Url                               | string | http://127.0.0.1:3000 |
| `feeStrategy`       | fee calculation strategy (`low`, `middle`, `high`)    | string | `middle`              |
| `signerPrivateKey`  | Transaction signer private key                        | string | -                     |
| `replicatorKeys`    | List of replicator public keys divided by whitespaces | string | -                     |

### Example

```shell
./onboarding -url=http://127.0.0.1:3000 -feeStrategy=middle -signerPrivateKey=0000000000000000000000000000000000000000000000000000000000000000 replicatorKeys="0000000000000000000000000000000000000000000000000000000000000000 0000000000000000000000000000000000000000000000000000000000000000"
```
