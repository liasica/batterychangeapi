package packed

import "github.com/gogf/gf/os/gres"

func init() {
	if err := gres.Add("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"); err != nil {
		panic("add binary content to resource manager failed: " + err.Error())
	}
}
