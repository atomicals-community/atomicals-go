# atomicals-core ⚛️

Atomicals: 是一个使用染色方法在BTC链上发行资产的协议，作者[Arthur's twitter(X)](https://twitter.com/atomicalsxyz)，该协议包括：域名机制（realm），nft，ft（arc-20）

- 目前Arthur 并未以文档或protocal形式披露atomicals的具体内容。但是提供了一个python版本的实现，包括：
    - atomicals索引器[atomicals-electrumx](https://github.com/atomicals/atomicals-electrumx)
    - atomicals交易发送工具[atomicals-js](https://github.com/atomicals/atomicals-js)命令工具

Atomicals-core: 是atomicals索引器atomicals-electrumx的golang版本，并以文本方式提供了atomicals协议的详细内容（在本仓库的doc目录下）

- 请注意，如果Arthur想要对atomicals-electrumx进行协议升级或更新，Atomicals-core是滞后的。
- 在未来一段时间内[yiming](https://github.com/yimingWOW)仍然会维护该项目（及时同步atomicals-electrumx的更新）
- 我并没有完成atomicals-core的全部测试工作，我没有趁手的机器在本地跑btc local node，这太花费我的时间了，atomicals生态项目方比我更有动力做这件事。
- 测试工作非常简单
    - 如果您同步完全部节点后发现某个账户资产x和atomicals-electrum的结果不一致，您清理库表，重新运行atomicals-core至x被mint或deploy的交易，即可找到x在atomicals-core下没有被正确mint的原因
    - 如果是transfer后的结果不一致，同样的操作，您检查下atomicals/transferNft.go和atomicals/transferFt.go的逻辑即可


## [atomicals-core文档目录](https://github.com/yimingWOW/atomicals-core/tree/main/doc)
- [atomicals-core 架构](https://github.com/yimingWOW/atomicals-core/tree/main/doc/0.atomicalsCoreFramework.md)
- [UXTO染色币原理](https://github.com/yimingWOW/atomicals-core/tree/main/doc/1.utxoColor.md)
- atomicals protocal 链上命令解析器+indexer检查条件
    - 部署和铸造
        - [dft 部署distributed ft](https://github.com/yimingWOW/atomicals-core/tree/main/doc/3.dft.md)
        - [dmt 铸造distributed ft](https://github.com/yimingWOW/atomicals-core/tree/main/doc/4.dmt.md)
        - [nft 铸造nft](https://github.com/yimingWOW/atomicals-core/tree/main/doc/5.nft.md)
        - [ft  铸造Direct ft](https://github.com/yimingWOW/atomicals-core/tree/main/doc/6.ft.md)
    - 转账
        - nft
            - x 拆分
            - y 移动
        - ft
            - x 拆分
            - y 移动
    - payment

## How to run atomicals-core
- you need a btc node
- install golang and docker
- start a postgres sql by docker
```
$ docker run --name YourDatabaseName -p 5432:5432 -e POSTGRES_DB=postgres -e POSTGRES_USER=yourUserName -e POSTGRES_PASSWORD=yourPassword postgres
``` 
- download atomicals-core
- edit conf/config.json : update it with your btc node url, user and password, and sql_dns(you have got it )
- cd to atomicals-core path
``` 
go mod tidy
go run ./
``` 

## Performance
- atomicals-core will spend 2.5s to process one block. if currentBlockHeight=834773, it will take about 20 hours to sync all btc blocks

## TODO:
- TraceTx可以并发处理
- PayLoad.Arg中有一个parents字段，它在get_mint_info_op_factory中被使用，但是我还没有捕捉到这个字段
- PayLoad中有一个存储nft图片信息的字段：image.png，但是在golang中解析为cbor结构体时有问题，我还没有解决它
- payment的逻辑我没有写，看起来不影响atomicals资产归属
- 下面几个atomicals命令我没有写处理逻辑
    - operationType = "dat" // dat - Store data on a transaction (dat)
    - operationType = "evt" // evt - Message response/reply
    - operationType = "mod" // mod - Modify general state
    - operationType = "sl" // sl - Seal an NFT and lock it from further changes forever

- atomicals协议文档编写
    - 累了，有志者去完成它吧，您只需要结合atomicals/witness/payload.go中PayLoad字段与atomicals/operation*.go，将每种atomicals命令对应的结构体信息和后续处理逻辑写清楚即可
    - 或者说我现在有了更重要的事情去做

- http接口
    - atomicals项目方比我更加有动力去完成它。
    - 非常简单，您只需要使用sql启动atomicals-core，再写一个http服务读sql中的表即可

- golang命令
    - 这个不属于indexer的职能，开始我很有动力去完成它，但是最近的紧张的开发让我身心疲惫，我现在觉的它是在没必要，有一些golang btc钱包库做的非常好，我觉得感兴趣的人可以去看看


这个项目的后续开发和测试工作还有很多，欢迎感兴趣的开发者和项目方联系我一起构建它。请各位看观有多提pr和isssus，代码中不合理之处尽管提出来，我一个人的能力有限，感兴趣的小伙伴共建才能越来越好


## futher more
最近的atomicals-core的重构让我对btc生态有了全新的认识，并且有了更加大胆的想法！期待我的新仓库吧！