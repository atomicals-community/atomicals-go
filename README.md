# atomicals-go ⚛️

Atomicals: 是一个使用染色方法在BTC链上发行资产的协议，作者[Arthur's twitter(X)](https://twitter.com/atomicalsxyz)，该协议包括：realm，nft和ft

- 目前Arthur 并未以文档或protocal形式披露atomicals的具体内容。但是提供了一个python版本的实现，包括：
    - atomicals索引器[atomicals-electrumx](https://github.com/atomicals/atomicals-electrumx)
    - atomicals交易发送工具[atomicals-js](https://github.com/atomicals/atomicals-js)命令工具
- Atomicals-go: 是atomicals索引器atomicals-electrumx的golang版本，并以文本方式提供了atomicals协议的详细内容（在本仓库的doc目录下）

- 在未来一段时间内[yiming](https://twitter.com/isyiming)仍然会维护该项目（及时同步atomicals-electrumx的更新）
- 如果您想加入，可以通过twitter联系我：[yiming](https://twitter.com/isyiming)
- 或者为我捐款: bc1p7uaqs0qq40mxqyljd93raxullh0ece2xvns5s5y9700v4ec0qjmsdt2q2n 接受任何类型的资产


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
- atomicals-core will spend 2.5s per block. if currentBlockHeight=834773, it will take about 20 hours to sync all btc blocks

## TODO:
补齐如下命令处理逻辑
- operationType = "dat" // dat - Store data on a transaction (dat)
- operationType = "evt" // evt - Message response/reply
- operationType = "sl" // sl - Seal an NFT and lock it from further changes forever
- payment
为api-service服务提供必要的http接口

