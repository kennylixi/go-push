package push

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gvalid"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tsms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type ConfigTencent struct {
	AccessKey string `json:"accessKey" v:"required#AccessKey不能为空"`
	SecretKey string `json:"secretKey" v:"required#SecretKey不能为空"`
	Region    string `json:"region" v:"required#Region不能为空"` // 地域信息
	AppID     string `json:"appID" v:"required#AppID不能为空"`   // 短信应用ID: 短信SdkAppId在 [短信控制台] 添加应用后生成的实际SdkAppId，示例如1400006666
	Endpoint  string // SDK会自动指定域名。通常是不需要特地指定域名的，但是如果你访问的是金融区的服务 则必须手动指定域名，例如sms的上海金融区域名： sms.ap-shanghai-fsi.tencentcloudapi.com
	SignName  string `json:"signName" v:"required#SignName不能为空"` // 短信签名
}

// TencentAdapter 腾讯云
type TencentAdapter struct {
	Config *ConfigTencent
	Client *tsms.Client
}

func NewAdapterTencent(ctx context.Context, i interface{}) (Adapter, error) {
	cfg := (*ConfigTencent)(nil)
	if err := gconv.Scan(i, &cfg); err != nil {
		return nil, err
	}
	if verr := gvalid.New().Data(&cfg).Run(ctx); verr != nil {
		if err := verr.FirstError(); err != nil {
			return nil, err
		}
	}

	t := &TencentAdapter{
		Config: cfg,
	}
	credential := common.NewCredential(cfg.AccessKey, cfg.SecretKey)
	cpf := profile.NewClientProfile()
	if cfg.Endpoint == "" {
		cpf.HttpProfile.Endpoint = cfg.Endpoint
	}
	client, err := tsms.NewClient(credential, cfg.Region, cpf)
	if err != nil {
		return nil, err
	}
	t.Client = client
	return t, nil
}

func (c *TencentAdapter) Send(ctx context.Context, accounts []string, template interface{}, templateParams map[string]string) (err error) {
	tmp := (*SmsParamTemplate)(nil)
	if err = gconv.Scan(template, &tmp); err != nil {
		return
	}
	if verr := gvalid.New().Data(&tmp).Run(ctx); verr != nil {
		if err = verr.FirstError(); err != nil {
			return
		}
	}

	request := tsms.NewSendSmsRequest()
	request.SmsSdkAppId = common.StringPtr(c.Config.AppID)
	request.SignName = common.StringPtr(c.Config.SignName)
	if templateParams != nil {
		params := gconv.SliceStr(templateParams)
		request.TemplateParamSet = common.StringPtrs(params)
	}
	request.TemplateId = common.StringPtr(tmp.TemplateId)
	request.PhoneNumberSet = common.StringPtrs(accounts)
	res, err := c.Client.SendSms(request)
	if err != nil {
		return err
	}
	if res.Response != nil {
		for _, status := range res.Response.SendStatusSet {
			if *status.Code != "Ok" {
				err = gerror.New(*status.Message)
				return
			}
		}
	}

	return nil
}
