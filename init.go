package push

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
)

const (
	TypeAliyun  = "aliyun"  // 阿里云大鱼
	TypeBaidu   = "baidu"   // 百度云
	TypeQiniu   = "qiniu"   // 七牛云
	TypeTencent = "tencent" // 腾讯云
	TypeMail    = "mail"    // 邮件
)

type NewAdapter func(ctx context.Context, i interface{}) (Adapter, error)

var (
	adapterFuns = map[string]NewAdapter{
		TypeAliyun:  NewAdapterAliyun,
		TypeBaidu:   NewAdapterBaidu,
		TypeQiniu:   NewAdapterQiniu,
		TypeTencent: NewAdapterTencent,
		TypeMail:    NewAdapterMail,
	}
)

// New 实例化适配器
func New(ctx context.Context, adapterType string, cfg interface{}) (adapter Adapter, err error) {
	adapterFun, ok := adapterFuns[adapterType]
	if !ok || adapterFun == nil {
		err = gerror.Newf("适配器[%s]不存在", adapterType)
		return
	}
	adapter, err = adapterFun(ctx, cfg)
	return
}
