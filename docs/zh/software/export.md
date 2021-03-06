# 导出区块链状态

## 介绍

IRISnet支持导出区块链状态，并以json格式返回给用户。如果把返回的json字符串保存到一个json文件里，这个json文件可以作为一个新区块链网络的创世块。导出区块链状态所用的命令为`iris export`。

默认情况下，IRISnet存储每10,000个块和最后100个块的快照。可以导出任一现有快照高度的区块链状态。

如果想导出不存在的快照高度的状态，首先需要[reset](reset.md)区块链状态到指定的高度。

## 用法

```
iris export <flags>
```

## 标志

| 名称，速记          | 类型   | 是否必填 | 默认值   | 介绍    |
| ------------------- | -----  | -------- | -------- | -------------- |
| --for-zero-height   | bool   | false    | false    | 导出数据之前做一些清理性的工作，如果不想以导出的数据启动一条新链，可以不加这个标志 |
| --height            | uint   | false    | 0        | 从指定的高度导出，默认值为0表示导出当前高度状态 |	
| --home              | string | false    | $HOME/.iris | 指定存储配置和区块链数据的目录 |
| --output-file       | string | false    | genesis.json |  存储导出状态的文件 |

## 示例

1. 导出当前区块链状态
```
iris export --home=<path_to_your_home>
```

2. 从指定高度导出区块链状态，该高度必须是现有快照高度
```
iris export --height 10000 --home=<path_to_your_home>
```

3. 如果想导出指定高度的区块链状态，并且以这个状态启动一条新链，可以尝试这个命令
```
iris export --height 10000 --for-zero-height --home=<path_to_your_home>
```