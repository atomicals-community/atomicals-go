### A solution to swap atomicals ft based on btc chain

#### creat pool
``` 
// operation and args in witness
"cpool"
type Args struct {
	LiquidityATickerName string 
	LiquidityAAmount int64
	LiquidityBTickerName string
	LiquidityBAmount int64
	LpAmount int64
	FeeRate  float
}
``` 
Craft a Bitcoin transaction containing a "cpool" operation with accompanying arguments. Subsequently, we'll establish a liquidity pool. Atomicals indexer will then generate a structure containing essential pool information.
``` 
// pool data updated by indexer
type PoolInfo struct {
	AtomicalsID string
	LiquidityATickerName string 
	LiquidityAAmount float64
	LiquidityBTickerName string
	LiquidityBAmount float64
	LpAmount float64
	FeeRate  float
}
``` 
#### add liquidity
``` 
"al" 
type Args struct {
	PoolAtomicalsID string 
}
``` 
When the indexer retrieves an "al" operation, it will examine the transaction inputs (Vins), remove the corresponding UTXOFtInfo entries from the Atomicals database(atomicalsUTXOFt.go), and update the PoolInfo. Additionally, it will designate the transaction outputs (Vouts) as pool tokens by coloring them.
``` 
type PoolToken struct {
	AtomicalsID string
	UserPk string
	PoolAtomicalsID string // pool's AtomicalsID
	LpAmount float64
}
``` 
#### remove liquidity
``` 
"rl" 
type Args struct {
	PoolAtomicalsID string
}
``` 
When the indexer retrieves an "rl" operation, it will examine the transaction inputs (Vins), remove the corresponding PoolToken entries, update the PoolInfo accordingly, and designate the transaction outputs (Vouts) as LiquidityA and LiquidityB by coloring them.

#### swap
``` 
"swap" 
type Args struct {

	PoolAtomicalsID string
}
``` 
When the indexer encounters a "swap" operation, it verifies the transaction inputs (Vins), removes the corresponding UTXOFtInfo entries, updates the PoolInfo accordingly, and designates the transaction outputs (Vouts) as LiquidityA or LiquidityB by coloring them.

#### some issue
There is an issue with the current system. When a user sends a swap transaction, they prepare specific outputs (Vouts), but the PoolInfo is constantly updated. According to the current rule, when someone swaps 10 units of FtA for 20 units of FtB, if their output value exceeds 20, all of their FtB will be burned. Conversely, if their output value is less than 20, only the excess FtB above 20 will be burned.

A viable solution is to prioritize UTXOs that haven't been fully colored (i.e., where the value exceeds the amount of tokens). This approach enables the implementation of non-integer splitting of Atomicals Ff. And implementing this approach is straightforward and requires minimal effort.





