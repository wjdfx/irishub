# iriscli status

## 描述

查询远端节点状态

## 用法

```shell
iriscli status <flags>
```

## 标志

| 名称, 速记             | 默认值                | 描述                                | 必需     |
| --------------------- | --------------------- | ----------------------------------- | -------- |
| --help, -h            |                       | 状态命令帮助                         |          |
| --node, -n            | tcp://localhost:26657 | Tendermint节点的rpc地址              |          |

## 例子

### 查询远端节点状态

```shell
iriscli status
```

查询到的结果如下所示：

```json
{
	"node_info": {
		"protocol_version": {
			"p2p": "4",
			"block": "7",
			"app": "0"
		},
		"id": "959185fdc3d14bdc7be1af40c5290d25042a454c",
		"listen_addr": "tcp://0.0.0.0:26656",
		"network": "test",
		"version": "0.26.0",
		"channels": "4020212223303800",
		"moniker": "node0",
		"other": {
			"tx_index": "on",
			"rpc_address": "tcp://0.0.0.0:26657"
		}
	},
	"sync_info": {
		"latest_block_hash": "04A6B890A61F503A64F254CF8479C8FB9012A9C9494249DC76F81B6453ADF6A1",
		"latest_app_hash": "B3549258BBC34860630BB5721364104DAC241EB243A8B0BCA0AA4968A64A1A6B",
		"latest_block_height": "2647",
		"latest_block_time": "2018-11-16T03:12:46.701163933Z",
		"catching_up": false
	},
	"validator_info": {
		"address": "91679AB00C0A09B006F9A812AAF686092657F658",
		"pub_key": {
			"type": "tendermint/PubKeyEd25519",
			"value": "r4r9TUJgKF8xxANw+8aMy9OP6rdwIFM6iUa8KVUaofo="
		},
		"voting_power": "100"
	}
}
```
