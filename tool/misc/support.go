package misc

import (
	"fmt"
	"github.com/intmian/mian_go_lib/tool/xlog"
)

func GoWaitError(log *xlog.Mgr, c <-chan error, from string, s string) {
	go func() {
		err := <-c
		if err != nil {
			log.Log(xlog.EError, from, fmt.Sprintf("%s:%s", s, err.Error()))
		}
	}()
}
