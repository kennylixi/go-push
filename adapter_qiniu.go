package push

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gvalid"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/sms"
)

type ConfigQiniu struct {
	AccessKey   string `json:"accessKey" v:"required#AccessKey不能为空"`
	SecretKey   string `json:"secretKey" v:"required#SecretKey不能为空"`
	SignatureId string `json:"signatureId" v:"required#SignatureId不能为空"` //短信签名ID
}

// QiniuAdapter 七牛云
type QiniuAdapter struct {
	Config  *ConfigQiniu
	Manager *sms.Manager
}

func NewAdapterQiniu(ctx context.Context, i interface{}) (Adapter, error) {
	cfg := (*ConfigQiniu)(nil)
	if err := gconv.Scan(i, &cfg); err != nil {
		return nil, err
	}
	if verr := gvalid.New().Data(&cfg).Run(ctx); verr != nil {
		if err := verr.FirstError(); err != nil {
			return nil, err
		}
	}

	q := &QiniuAdapter{
		Config: cfg,
	}
	q.Manager = sms.NewManager(auth.New(cfg.AccessKey, cfg.SecretKey))
	return q, nil
}

func (b *QiniuAdapter) Send(ctx context.Context, accounts []string, template interface{}, templateParams map[string]string) (err error) {
	tmp := (*SmsParamTemplate)(nil)
	if err = gconv.Scan(template, &tmp); err != nil {
		return
	}
	if verr := gvalid.New().Data(&tmp).Run(ctx); verr != nil {
		if err = verr.FirstError(); err != nil {
			return
		}
	}

	req := sms.MessagesRequest{}
	req.SignatureID = b.Config.SignatureId
	req.TemplateID = tmp.TemplateId
	req.Mobiles = accounts
	if templateParams != nil {
		req.Parameters = gconv.Map(templateParams)
	}
	resp, err := b.Manager.SendMessage(req)
	if err != nil {
		return err
	}
	if len(resp.JobID) == 0 {
		return gerror.New("发送短信异常")
	}
	return
}
