# Replicator tree rebuild CLI tool

Rebuilds AVL tree of replicators stored in Replicator and Queue caches to ensure its validity.

## Usage

All public keys of replicators that exist in Replicator cache should be supplied in `replicatorKeys`; otherwise
the AVL tree of replicators may end up in an invalid state. If this happpens, it can be fixed by simply executing this 
tool again with correct list of replicator keys.

### Flags

| Name                | Description                                           | Type   | Default               |
|:--------------------|:------------------------------------------------------|:-------|:----------------------|
| `url`               | ProximaX Chain REST Url                               | string | http://127.0.0.1:3000 |
| `feeStrategy`       | Fee calculation strategy (`low`, `middle`, `high`)    | string | `middle`              |
| `signerPrivateKey`  | Transaction signer private key                        | string | -                     |
| `replicatorKeys`    | List of replicator public keys separated by whitespaces | string | -                     |

### Example

```shell
./replicator_tree_rebuild -url=http://127.0.0.1:3000 -feeStrategy=middle -signerPrivateKey=0000000000000000000000000000000000000000000000000000000000000000 replicatorKeys="0000000000000000000000000000000000000000000000000000000000000000 0000000000000000000000000000000000000000000000000000000000000000"
```
