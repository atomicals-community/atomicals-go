package main

import (
	"github.com/atomicals-core/atomicals"
	"github.com/atomicals-core/atomicals/common"
	"github.com/atomicals-core/pkg/btcsync"
)

func main() {
	b, err := btcsync.NewBtcSync("rpcURL", "user", "password")
	if err != nil {
		panic(err)
	}
	a := atomicals.NewAtomicals(b, common.ATOMICALS_ACTIVATION_HEIGHT-1) // common.ATOMICALS_ACTIVATION_HEIGHT
	for {
		a.TraceBlock()
	}

	// traced Specific Tx
	// tx, err := b.GetTransaction("b28f089b5a96c4803db73d51ed801aec4efec997761ee8dc914e0f934b6fcd59")
	// if err != nil {
	// 	log.Log.Warnf("GetTransaction err:%v", err)
	// }
	// a.TraceTx(*tx, height)
}
