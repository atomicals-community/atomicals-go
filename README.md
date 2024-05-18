# atomicals-go ⚛️

Atomicals: 是一个使用染色方法在BTC链上发行资产的协议，作者[Arthur's twitter(X)](https://twitter.com/atomicalsxyz)，该协议包括：realm，nft和ft

- 目前Arthur 并未以文档或protocal形式披露atomicals的具体内容。但是提供了一个python版本的实现，包括：
    - atomicals索引器[atomicals-electrumx](https://github.com/atomicals/atomicals-electrumx)
    - atomicals交易发送工具[atomicals-js](https://github.com/atomicals/atomicals-js)命令工具
- Atomicals-go: 是atomicals索引器atomicals-electrumx的golang版本，并以文本方式提供了atomicals协议的详细内容（在本仓库的doc目录下）

- 请注意，如果Arthur想要对atomicals-electrumx进行协议升级或更新，Atomicals-go 是滞后的
- 在未来一段时间内[yiming](https://twitter.com/isyiming)仍然会维护该项目（及时同步atomicals-electrumx的更新）
- 如果您想加入，可以通过twitter联系我：[yiming](https://twitter.com/isyiming)

我需要一个aws服务器运行btc全节点和atomicals-electrumx索引器，用来对比atomicals-go和atomicals-electrumx的结果是否有差异，有人愿意帮助我吗？麻烦通过推特和我联系。

或者为我捐款：bc1p2ty4uj7g9l7w4dmu0qrmm2z35jm22t2r8qa5uwfde5vz3r3mhtgsdhcs4u 接受任何类型的资产
## How to run atomicals-go
1. run a local btc node
```
// cd to a path u want to save btc node file 
mkdir btc

wget https://bitcoincore.org/bin/bitcoin-core-26.0/bitcoin-26.0-arm64-apple-darwin.tar.gz

tar -xzvf bitcoin-26.0-x86_64-linux-gnu.tar.gz

mv bitcoin-26.0 bitcoin

vim ./bitcoin/bitcoin.conf

```
```
Edit bitcoin.conf, add these params for main net. we run btc node with prune mode and set assumevalid=0000000000000000000211eb82135b8f5d8be921debf8eff1d6b38b73bc03834.
Atomicals protocal start from blockHeight=808080, we don't need all blockInfo.

# Options for mainnet
[main]
dbcache=1024
server=1
rest=1
daemon=1
rpcbind=0.0.0.0:8332 
rpcallowip=0.0.0.0/0 
rpcuser=btc
rpcpassword=btc2012
prune=240000
assumevalid=0000000000000000000211eb82135b8f5d8be921debf8eff1d6b38b73bc03834
```

2. install golang and docker
3. start a postgres sql by docker
```
$ docker run --name postgres -p 5432:5432 -e POSTGRES_DB=postgres -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin123 postgres:14
``` 
4. run atomicals indexer
download atomicals-core
- edit conf/config.json update it with your btc node url, user and password, and sql_dns:
```
{
    "btc_rpc_url" : "0.0.0.0:8332",
    "btc_rpc_user": "btc" ,
    "btc_rpc_password": "btc2012",
    "sql_dns": "host=127.0.0.1 user=admin password=admin123 dbname=atomicals port=5432 sslmode=disable"
}
```
``` 
// cd to atomicals-core path
go mod tidy

// init sql table
cd repo/postsql/init/
go run ./

// start indexer
cd ../../../
go run ./  

// or run it with nohup
nohup go run ./ > log.txt 2>&1 &
``` 

## Performance
- atomicals-core will spend 2.5s to process one block. if currentBlockHeight=834773, it will take about 20 hours to sync all btc blocks

## TODO:

1. ft 资产部署和转账 测试 (yimingWoW 正在做测试；我已经跑到blockHeight=838408，实际上这个工作早就应该完成了，但是我还有我自己的工作，只能在闲暇时间搞这些，时间比较琐碎,我预计下周完成它)
2. nft-realm及subrealm部署 逻辑完善+测试
    - merkleverify部分比较复杂，我写了个大概，还没有写完：atomicals-core/operation/merkleVerify.go； 也可能是最近动力不足了，这部份重构让我觉得很痛苦。如果有人愿意把这个逻辑完善，我非常乐意
3. nft-container及item部署 逻辑完善+测试
4. nft 部署 逻辑完善+测试
    - PayLoad中有一个存储nft图片信息的字段：image.png，但是在golang中解析为cbor结构体时有问题，我还没有解决它
    - PayLoad.Arg中有一个parents字段，它在get_mint_info_op_factory中被使用，但是我还没有捕捉到这个字段
5. 补齐如下命令处理逻辑
    - operationType = "dat" // dat - Store data on a transaction (dat)
    - operationType = "evt" // evt - Message response/reply
    - operationType = "mod" // mod - Modify general state
    - operationType = "sl" // sl - Seal an NFT and lock it from further changes forever
6. TraceTx改为并发处理
7. atomicals协议文档编写
8. 一个api服务，提供必要的http接口(yimingWOW 提供了一个示例，其他人可以参考api-service/logic/asset.go 添加你认为必要的http接口。另外请您暂时不要写缓存逻辑，缓存应该在repo层，等我们把indexer逻辑和全部http接口确定后，再在repo写好缓存更新逻辑)
9. golang命令工具

