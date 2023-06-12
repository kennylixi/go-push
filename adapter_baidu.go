package push

import (
	"context"
	"github.com/baidubce/bce-sdk-go/services/sms"
	smsapi "github.com/baidubce/bce-sdk-go/services/sms/api"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gvalid"
)

type ConfigBaidu struct {
	AccessKey   string `json:"accessKey" v:"required#AccessKey不能为空"`
	SecretKey   string `json:"secretKey" v:"required#SecretKey不能为空"`
	Endpoint    string `json:"endpoint" v:"required#Endpoint不能为空"`
	SignatureId string `json:"signatureId" v:"required#SignatureId不能为空"` //短信签名ID
}

// BaiduAdapter 百度云
type BaiduAdapter struct {
	Config *ConfigBaidu
	Client *sms.Client
}

func NewAdapterBaidu(ctx context.Context, i interface{}) (Adapter, error) {
	cfg := (*ConfigBaidu)(nil)
	if err := gconv.Scan(i, &cfg); err != nil {
		return nil, err
	}
	if verr := gvalid.New().Data(&cfg).Run(ctx); verr != nil {
		if err := verr.FirstError(); err != nil {
			return nil, err
		}
	}

	b := &BaiduAdapter{
		Config: cfg,
	}

	client, err := sms.NewClient(cfg.AccessKey, cfg.SecretKey, cfg.Endpoint)
	if err != nil {
		return nil, err
	}
	b.Client = client
	return b, nil
}

func (b *BaiduAdapter) Send(ctx context.Context, accounts []string, template interface{}, templateParams map[string]string) (err error) {
	tmp := (*SmsParamTemplate)(nil)
	if err = gconv.Scan(template, &tmp); err != nil {
		return
	}
	if verr := gvalid.New().Data(&tmp).Run(ctx); verr != nil {
		if err = verr.FirstError(); err != nil {
			return
		}
	}

	req := &smsapi.SendSmsArgs{}
	req.SignatureId = b.Config.SignatureId
	req.Template = tmp.TemplateId
	req.Mobile = garray.NewStrArrayFrom(accounts).Join(",")
	if templateParams != nil {
		req.ContentVar = gconv.Map(templateParams)
	}
	resp, err := b.Client.SendSms(req)

	if resp != nil && resp.Code != "1000" {
		return gerror.New(resp.Message)
	}
	return
}
