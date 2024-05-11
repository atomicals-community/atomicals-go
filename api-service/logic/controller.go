package logic

import (
	"github.com/atomicals-go/pkg/conf"
	"github.com/atomicals-go/pkg/errors"
	"github.com/atomicals-go/repo"
	"github.com/beego/beego/v2/server/web"
)

var DB repo.DB

type Controller struct {
	web.Controller
	repo.DB
}

func (c *Controller) Prepare() {
	if DB == nil {
		panic("no init connection")
	}
	c.DB = DB
}

func InitController(conf *conf.Config) {
	DB = repo.NewSqlDB(conf.SqlDNS)
}

type Resp struct {
	Error errors.Error `json:"error"`
	Data  interface{}  `json:"data"`
}

func (c *Controller) ReturnError(err errors.Error) {
	if err := c.Ctx.JSONResp(Resp{
		Error: err,
		Data:  nil,
	}); err != nil {
		return
	}
}

func (c *Controller) ReturnSuccess(data interface{}) {
	if err := c.Ctx.JSONResp(Resp{
		Error: nil,
		Data:  data,
	}); err != nil {
		return
	}
}
