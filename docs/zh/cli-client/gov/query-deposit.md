# iriscli gov query-deposit

## 描述

查询保证金的充值明细

## 使用方式

```
iriscli gov query-deposit <flags>
```
打印帮助信息:

```
iriscli gov query-deposit --help
```
## 标志

| 名称, 速记       | 默认值                 | 描述                                                                                                                                                 | 是否必须  |
| --------------- | --------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --depositor     |                       | bech32编码的存款人地址                                                                                                                    | Yes      |
| --proposal-id   |                       | 提议ID                                                                                                        | Yes      |

## 例子

### 查询充值保证金

```shell
iriscli gov query-deposit --chain-id=<chain-id> --proposal-id=<proposal-id> --depositor=<depositor_address>
```

通过指定提议、指定存款人查询保证金充值详情，得到结果如下：

```txt
Deposit by iaa1c4kjt586r3t353ek9jtzwxum9x9fcgwent790r on Proposal 90 is for the amount 995iris
```
