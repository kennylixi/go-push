package push

// SmsParamTemplate 短信参数模板
type SmsParamTemplate struct {
	TemplateId      string `json:"templateId" v:"required#模板编码不能为空"`      // 模板ID
	TemplateContent string `json:"templateContent" v:"required#模板内容不能为空"` // 内容模板（需要和短信平台配置模板一致）
}

// MailParamTemplate 邮件参数模板
type MailParamTemplate struct {
	Title   string `json:"title" v:"required#邮件标题不能为空"`   // 邮件标题模板
	Content string `json:"content" v:"required#邮件内容不能为空"` // 邮件内容模板
}
