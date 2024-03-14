# atomicals-core ⚛️

Atomicals: 是一个使用染色方法在BTC链上发行资产的协议，作者[Arthur's twitter(X)](https://twitter.com/atomicalsxyz)，该协议包括：域名机制（realm），nft，ft（arc-20）

- 目前Arthur 并未以文档或protocal形式披露atomicals的具体内容。但是提供了一个python版本的实现，包括：
    - atomicals索引器[atomicals-electrumx](https://github.com/atomicals/atomicals-electrumx)
    - atomicals交易发送工具[atomicals-js](https://github.com/atomicals/atomicals-js)命令工具

Atomicals-core: 是atomicals索引器atomicals-electrumx的golang版本，并以文本方式提供了atomicals协议的详细内容（在本仓库的doc目录下）

- 请注意，如果Arthur想要对atomicals-electrumx进行协议升级或更新，Atomicals-core是滞后的。
- 在未来一段时间内[yiming](https://github.com/yimingWOW)仍然会维护该项目（及时同步atomicals-electrumx的更新）
- 期待更多项目方或开发者加入，我会将本仓库提交到某个公开项目下，并移交仓库管理权限

## [atomicals-core文档目录](https://github.com/yimingWOW/atomicals-core/tree/main/doc)
- [UXTO染色币原理](https://github.com/yimingWOW/atomicals-core/tree/main/doc/1.utxoColor.md)
- [atomicals-core 架构](https://github.com/yimingWOW/atomicals-core/tree/main/doc/0.atomicalsCoreFramework.md)
- atomicals protocal 链上命令解析器+indexer检查条件
    - 部署和铸造
        - [dft 部署distributed ft](https://github.com/yimingWOW/atomicals-core/tree/main/doc/3.dft.md)
        - [dmt 铸造distributed ft](https://github.com/yimingWOW/atomicals-core/tree/main/doc/4.dmt.md)
        - [nft 铸造nft](https://github.com/yimingWOW/atomicals-core/tree/main/doc/5.nft.md)
        - [ft  铸造Direct ft](https://github.com/yimingWOW/atomicals-core/tree/main/doc/6.ft.md)
    - 转账
        - x 拆分
        - y 移动
    - payment

## How to run atomicals-core
- you need a btc node and make sure golang has been installed on ur os.
- download atomicals-core
- edit atomicals-core/main.go : replace btcsync.NewBtcSync("rpcURL", "user", "password") with your btc node url, user and password 
- cd to atomicals-core path
- run: go mod tidy
- run: go run ./
there are many unnencessary log, i will delete them when we have a stable version.

## TODO:
- atomicals协议文档编写
- atomicals optionType未完成：
    - operationType = "dat" // dat - Store data on a transaction (dat)
    - operationType = "evt" // evt - Message response/reply
    - operationType = "mod" // mod - Modify general state
    - operationType = "sl" // sl - Seal an NFT and lock it from further changes forever
- transfer测试：
    - atomicals/transferFt.go
    - atomicals/transferNft.go
- http接口
- golang命令

- atomicals optionType待测试：
    - atomicals/operationDmt.go
    - atomicals/operationDft.go
    - atomicals/operationNft.go
    - atomicals/operationFt.go
- 存储层抽象


这个项目的后续开发和测试工作还有很多，欢迎感兴趣的开发者和项目方联系我一起构建它

请各位看观有多提pr和isssus，代码中不合理之处尽管提出来，我一个人的能力有限，感兴趣的小伙伴共建才能越来越好