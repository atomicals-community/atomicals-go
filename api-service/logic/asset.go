package logic

import (
	"fmt"

	"github.com/atomicals-go/pkg/log"
)

func (c *Controller) NftByAtomicalsID() {
	var req ReqAssetByAtomicalsID
	if err := c.Ctx.BindJSON(&req); err != nil {
		log.Log.Errorf("AssetByAtomicalsID BindJSON err:%v", err)
		c.ReturnError(err)
		return
	}
	entity, err := c.NftUTXOsByLocationID(req.AtomicalsID)
	if err != nil {
		log.Log.Errorf("AssetByAtomicalsID NftUTXOsByLocationID err:%v", err)
		c.ReturnError(err)
		return
	}
	resp := &RespNftByAtomicalsID{
		Assets: entity,
	}
	fmt.Println(resp)
	c.ReturnSuccess(resp)
}

func (c *Controller) FtByAtomicalsID() {
	var req ReqAssetByAtomicalsID
	if err := c.Ctx.BindJSON(&req); err != nil {
		log.Log.Errorf("AssetByAtomicalsID BindJSON err:%v", err)
		c.ReturnError(err)
		return
	}
	entity, err := c.FtUTXOsByLocationID(req.AtomicalsID)
	if err != nil {
		log.Log.Errorf("AssetByAtomicalsID NftUTXOsByLocationID err:%v", err)
		c.ReturnError(err)
		return
	}
	resp := &RespFtByAtomicalsID{
		Assets: entity,
	}
	fmt.Println(resp)
	c.ReturnSuccess(resp)
}
