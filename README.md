# atomicals-core ⚛️

Atomicals: 是一个使用染色方法在BTC链上发行资产的协议，作者[Arthur's twitter(X)](https://twitter.com/atomicalsxyz)，该协议包括：realm，nft和ft

- 目前Arthur 并未以文档或protocal形式披露atomicals的具体内容。但是提供了一个python版本的实现，包括：
    - atomicals索引器[atomicals-electrumx](https://github.com/atomicals/atomicals-electrumx)
    - atomicals交易发送工具[atomicals-js](https://github.com/atomicals/atomicals-js)命令工具

Atomicals-core: 是atomicals索引器atomicals-electrumx的golang版本，并以文本方式提供了atomicals协议的详细内容（在本仓库的doc目录下）

- 请注意，如果Arthur想要对atomicals-electrumx进行协议升级或更新，Atomicals-core是滞后的
- 在未来一段时间内[yiming](https://twitter.com/isyiming)仍然会维护该项目（及时同步atomicals-electrumx的更新）

很高兴[Renaissance Lab](https://twitter.com/Renaissance_ARC)和[Atomicals Market](https://twitter.com/atomicalsmarket)加入到atomicals-core的开发中来。
感谢Arthur巧妙的构思，为我们带来了atomicals协议。我们需要对atomicals协议有更加清晰的认知，并在人们脑海中凝聚牢固的共识。
atomicals-core将致力于提供规范的atomicals协议文本和简洁高效的atomicals协议索引器。它是一个纯粹的开源，没有所谓的官方团队，项目方，欢迎任何对此感兴趣的人参与建设。您的参与会让atomicals更加去中心化！

## How to run atomicals-core
- you need a btc node
- install golang and docker
- start a postgres sql by docker
```
$ docker run --name postgres -p 5432:5432 -e POSTGRES_DB=postgres -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin123 postgres:14
``` 
- download atomicals-core
- edit conf/config.json : update it with your btc node url, user and password, and sql_dns("sql_dns": "host= 127.0.0.1 user=admin password=admin123 dbname=postgres port=5432 sslmode=disable")
- cd to atomicals-core path
``` 
go mod tidy
cd atomicals/DB/postsql/init/
go run ./
cd ../../../../
nohup go run ./ > log.txt 2>&1 &
``` 

## Performance
- atomicals-core will spend 2.5s to process one block. if currentBlockHeight=834773, it will take about 20 hours to sync all btc blocks


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

## TODO:
- ft 资产部署和转账 测试
- nft-realm及subrealm部署和转账 逻辑完善+测试
    - payment的逻辑
- nft-container及item部署和转账 逻辑完善+测试
    - payment的逻辑
- nft 部署和转账 逻辑完善+测试
    - PayLoad中有一个存储nft图片信息的字段：image.png，但是在golang中解析为cbor结构体时有问题，我还没有解决它
    - PayLoad.Arg中有一个parents字段，它在get_mint_info_op_factory中被使用，但是我还没有捕捉到这个字段
- 补齐如下命令处理逻辑
    - operationType = "dat" // dat - Store data on a transaction (dat)
    - operationType = "evt" // evt - Message response/reply
    - operationType = "mod" // mod - Modify general state
    - operationType = "sl" // sl - Seal an NFT and lock it from further changes forever
- TraceTx可以并发处理
- atomicals协议文档编写
- 一个api服务，提供必要的http接口
- golang命令工具

