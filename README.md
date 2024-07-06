# atomicals-go ⚛️

#### atomicals-go 是什么
- Atomicals: 是一个使用染色方法在BTC链上发行资产的协议，作者[Arthur's twitter(X)](https://twitter.com/atomicalsxyz)，该协议包括：realm，nft和ft
- 目前Arthur 并未以文档或protocal形式披露atomicals的具体内容。但是提供了一个python版本的实现，包括：
    - atomicals索引器[atomicals-electrumx](https://github.com/atomicals/atomicals-electrumx)
    - atomicals交易发送工具[atomicals-js](https://github.com/atomicals/atomicals-js)命令工具
- Atomicals-go: 是atomicals索引器atomicals-electrumx的golang版本，并以文本方式提供了atomicals协议的详细内容（在本仓库的doc目录下）

- 在未来一段时间内[github:random](https://github.com/yimingWOW)仍然会维护该项目（及时同步atomicals-electrumx的更新）
- 如果您想加入，可以通过twitter联系我：[x:random](https://twitter.com/isyiming)
- 或者为我捐款: bc1p7uaqs0qq40mxqyljd93raxullh0ece2xvns5s5y9700v4ec0qjmsdt2q2n 接受任何类型的资产

#### 嗨，atomicals-go终于完成了，我简单说一下这个indexer的优点

- 不需要btc全节点：一个完整的btc全节点需要730GB的磁盘空间。atomicals-go不需要btc全节点，你只需要以prune mode运行btc node，在你的电脑中只保存从808080高度开始的区块即可（这些区块大概占用大概140GB）
- 占用更少的存储空间：btc链上和atomicals协议无关的信息全部被过滤，atomicals-go只将有效数据都存储在sql中，这部份数据大概占用
- 防宕机：可以随时终止运行服务，即使是因为断电或者电脑死机等原因导致服务中断，没关系，只需要重启服务。它会在之前的区块高度继续同步，并且保证继续写入的数据是正确的
- 适应btc链分叉：无需担心btc链分叉的影响，保证通过atomicals-go查到的atomicals永远是最新的正确的，并且包括最新区块
- 支持查询mempool中的交易：即使某笔交易还没有被打包，只要你运行的btc节点可以查询到mempool中的交易，你就可以通过接口查看这笔atomicals交易包含的资产详情

## Performance
- atomicals-core will spend 2.5s per block. if currentBlockHeight=834773, it will take about 20 hours to sync all btc blocks
- 同步耗时平均2～4s/block，一天左右可以同步完成

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
$ docker run --name postgres -p 5432:5432 -e POSTGRES_DB=atomicals -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin123 -d postgres:14
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
cd repo/postsql/init
go run ./

// start indexer
cd atomicals-core
go run ./  
// or run it with nohup: nohup go run ./ > log.txt 2>&1 &
``` 
// start atomicals-api service if you need
cd atomicals-api
go run ./

#### TODO:
- 我发现payment没有什么用，所以atomicals-go没有保存任何payment信息，如果有必要，希望有人来完成它
- 为api-service服务提供更多必要的http接口

