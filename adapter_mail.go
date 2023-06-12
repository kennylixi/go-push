package push

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gvalid"
	"gopkg.in/gomail.v2"
)

type ConfigMail struct {
	Host     string `json:"host" v:"required#Host不能为空"`
	Port     int    `json:"port" v:"required#Port不能为空"`
	UserName string `json:"userName" v:"required#UserName不能为空"`
	Password string `json:"password" v:"required#Password不能为空"`
}

// MailAdapter 邮件
type MailAdapter struct {
	config *ConfigMail
	dialer *gomail.Dialer
}

func NewAdapterMail(ctx context.Context, i interface{}) (Adapter, error) {
	cfg := (*ConfigMail)(nil)
	if err := gconv.Scan(i, &cfg); err != nil {
		return nil, err
	}
	if verr := gvalid.New().Data(&cfg).Run(ctx); verr != nil {
		if err := verr.FirstError(); err != nil {
			return nil, err
		}
	}

	o := &MailAdapter{
		config: cfg,
	}

	o.dialer = gomail.NewDialer(cfg.Host, cfg.Port, cfg.UserName, cfg.Password)
	return o, nil
}

func (s *MailAdapter) Send(ctx context.Context, accounts []string, template interface{}, templateParams map[string]string) (err error) {
	tmp := (*MailParamTemplate)(nil)
	if err = gconv.Scan(template, &tmp); err != nil {
		return
	}
	if verr := gvalid.New().Data(&tmp).Run(ctx); verr != nil {
		if err = verr.FirstError(); err != nil {
			return
		}
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.config.UserName)
	m.SetHeader("To", accounts...)
	m.SetHeader("Subject", parserTemplate(tmp.Title, templateParams))
	m.SetBody("text/html", parserTemplate(tmp.Content, templateParams))

	return s.dialer.DialAndSend(m)
}
