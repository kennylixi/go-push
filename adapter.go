package push

import "context"

type Adapter interface {
	Send(ctx context.Context, accounts []string, template interface{}, templateParams map[string]string) (err error)
}
