package util

import (
	"bytes"
	"github.com/pkg/errors"
	"os/exec"
	xruntime "xcore/lib/runtime"
)

// Command 调用 linux 命令
//
//	参数:
//		args:"chmod +x /xx/xx/x.sh"
func Command(args string) (outStr string, errStr string, err error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("/bin/bash", "-c", args)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	outStr, errStr = string(stdout.Bytes()), string(stderr.Bytes())
	if err != nil {
		return outStr, errStr, errors.WithMessage(err, xruntime.Location())
	}
	return outStr, errStr, nil
}
