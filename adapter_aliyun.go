package push

import (
	"context"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gvalid"
)

type ConfigAliyun struct {
	AccessKey string `json:"accessKey" v:"required#AccessKey不能为空"`
	SecretKey string `json:"secretKey" v:"required#SecretKey不能为空"`
	SignName  string `json:"signName" v:"required#SignName不能为空"`
	Region    string `json:"region"`
}

func (c *ConfigAliyun) GetRegion() string {
	if c.Region == "" {
		return "cn-hangzhou"
	}
	return c.Region
}

// AliyunAdapter 阿里云
type AliyunAdapter struct {
	Config *ConfigAliyun
	Client *dysmsapi20170525.Client
}

func NewAdapterAliyun(ctx context.Context, i interface{}) (Adapter, error) {
	cfg := (*ConfigAliyun)(nil)
	if err := gconv.Scan(i, &cfg); err != nil {
		return nil, err
	}
	if verr := gvalid.New().Data(&cfg).Run(ctx); verr != nil {
		if err := verr.FirstError(); err != nil {
			return nil, err
		}
	}

	o := &AliyunAdapter{
		Config: cfg,
	}
	client, err := dysmsapi20170525.NewClient(&openapi.Config{
		AccessKeyId:     tea.String(cfg.AccessKey),
		AccessKeySecret: tea.String(cfg.SecretKey),
		RegionId:        tea.String(cfg.GetRegion()),
		Endpoint:        tea.String("dysmsapi.aliyuncs.com"),
	})

	if err != nil {
		return nil, err
	}
	o.Client = client
	return o, nil
}

func (s *AliyunAdapter) Send(ctx context.Context, accounts []string, template interface{}, templateParams map[string]string) (err error) {
	tmp := (*SmsParamTemplate)(nil)
	if err = gconv.Scan(template, &tmp); err != nil {
		return
	}
	if verr := gvalid.New().Data(&tmp).Run(ctx); verr != nil {
		if err = verr.FirstError(); err != nil {
			return
		}
	}

	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers: tea.String(garray.NewStrArrayFrom(accounts).Join(",")),
		SignName:     tea.String(s.Config.SignName),
		TemplateCode: tea.String(tmp.TemplateId),
	}

	if templateParams != nil {
		params, err := gjson.EncodeString(templateParams)
		if err != nil {
			return err
		}
		sendSmsRequest.TemplateParam = tea.String(params)
	}

	runtime := &util.RuntimeOptions{}
	err = func() (err error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				err = r
			}
		}()
		res, err := s.Client.SendSmsWithOptions(sendSmsRequest, runtime)
		if err != nil {
			if terr, ok := err.(*tea.SDKError); ok {
				err = gerror.New(terr.Error())
			}
			return err
		}
		if res != nil && tea.StringValue(res.Body.Code) == "OK" {
			return nil
		}
		err = gerror.New(tea.StringValue(res.Body.Message))
		return
	}()

	return err
}
