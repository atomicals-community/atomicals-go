# atomicals-core

## 喵🐱～
atomicals协议是一个构建于BTC上的染色币协议，但目前atomicals的具体内容并未以文档或protocal代码形式提供。目前atomicals作者只提供了一个python版本的索引器[atomicals-electrumx](https://github.com/atomicals/atomicals-electrumx)和一个[atomicals-js](https://github.com/atomicals/atomicals-js)命令工具

我想了解atomicals协议的具体格式，大约两周前我开始用golang重构该索引器atomicals-core，目前的构想是做到以下几点:

- 整理出协议规范
- 将协议本身和服务接口与存储剥离
- 提供高性能防宕机的indexer
- 提供golang命令行工具

## [atomicals-core文档目录](https://github.com/yimingWOW/atomicals-core/tree/main/doc)

- [UXTO染色币原理](https://github.com/yimingWOW/atomicals-core/tree/main/doc/1.utxoColor.md)
- atomicals operationType
    - ft部署和铸造
    - nft铸造
- 转账，拆分和合并
    - ft 转账，拆分和合并
    - nft 转账
- atomicals-core架构
- 存储层接入条件
    - 我会分别提供sql和redis的防宕机方案
    - 有其他db需求就按文档说明接入吧

## TODO:
- atomicals协议文档编写

- atomicals optionType待测试：

    - atomicals/operationDmt.go
    - atomicals/operationDft.go
    - atomicals/operationNft.go
    - atomicals/operationFt.go

- atomicals optionType未完成：

    - operationType = "dat" // dat - Store data on a transaction (dat)
    - operationType = "evt" // evt - Message response/reply
    - operationType = "mod" // mod - Modify general state
    - operationType = "sl" // sl - Seal an NFT and lock it from further changes forever

- transfer测试：
    - atomicals/transferFt.go
    - atomicals/transferNft.go

- 存储层抽象

- http接口

- golang命令


这个项目的后续开发和测试工作还有很多，欢迎感兴趣的开发者和项目方联系我一起构建它

请各位看观有多提pr和isssus，代码中不合理之处尽管提出来，我一个人的能力有限，感兴趣的小伙伴共建才能越来越好