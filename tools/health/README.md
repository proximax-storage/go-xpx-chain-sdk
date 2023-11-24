# Health CLI tool

The tool allows to wait for the list of nodes to synchronize the height and hash of the block at that height.

## Build

```shell
cd cmd && go build -o health
```

## Flags

- `height` - Expected height. Default is 0. Example `-height 1000`;
- `nodes` - List of nodes `<ip:port>=<nodePubKey>`. Example `-nodes 0.0.0.0:7900=10E8A1CCCFE02C4C22C12D42277520F1FC7D471E570C9FE2A2961ECB020BC596`;
- `discover` - Discover connected nodes to list of nodes(`-nodes`). Default is `true`. Example `-discover true`.

## Run

```shell
./health -height <height> -nodes <ip:port>=<nodePubKey> -discover
```

### Example

```shell
./health -height=98945 -nodes 54.151.169.225:7900=10E8A1CCCFE02C4C22C12D42277520F1FC7D471E570C9FE2A2961ECB020BC596 -discovery
```