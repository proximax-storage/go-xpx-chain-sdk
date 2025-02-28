# Replicators cleanup CLI tool

Version 1: Allows to remove replicators that are not bound with nodes.

Version 2: Allows to remove drive info from replicators.

## Usage

### Flags

| Name               | Description                                           | Type    | Default               |
|:-------------------|:------------------------------------------------------|:--------|:----------------------|
| `url`              | ProximaX Chain REST Url                               | string  | http://127.0.0.1:3000 |
| `feeStrategy`      | fee calculation strategy (`low`, `middle`, `high`)    | string  | `middle`              |
| `signerPrivateKey` | Transaction signer private key                        | string  | -                     |
| `replicatorKeys`   | List of replicator public keys divided by whitespaces | string  | -                     |
| `version`          | Transaction version (1 or 2)                          | integer | 1                     |

### Example

```shell
./onboarding -url=http://127.0.0.1:3000 -feeStrategy=middle -signerPrivateKey=0000000000000000000000000000000000000000000000000000000000000000 replicatorKeys="0000000000000000000000000000000000000000000000000000000000000000 0000000000000000000000000000000000000000000000000000000000000000" -version=2
```
