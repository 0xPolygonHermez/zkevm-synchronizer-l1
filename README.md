# zkEVM Synchronizer-L1
zkEVM Synchronizer-L1
This is a library to synchronize L1InfoTree events from L1

## Usage

### Import it
```
"github.com/0xPolygonHermez/zkevm-synchronizer-l1/synchronizer"
```
### Create object
You can choose to create from a config file object calling `NewSynchronizer` or from a **toml** config file using `NewSynchronizerFromConfigfile`


## Executing
The function `Sync` if a blocking function that starts the synchronization. You can run as a go rutine:
```
go mySync.Sync(false)
```
The parameter `returnOnSync` allows to return when the synchornization is finish so you can run the first time as a blocking call and, once the system is synchronized, run another time in the background


## Interacting
You can call to `IsSynced` to known if the system is syncrhonized

Check [usage.md](docs/usage.md) for a client example

test