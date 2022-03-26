package uapp

import (
	"context"
	"fmt"
	"github.com/WingGao/errors"
	"github.com/WingGao/go-utils/ucore"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/ungerik/go-dry"
	"os"
	"path/filepath"
	"strconv"
)

// 保持单进程运行
// 在文件所在目录下，创建一个 {key}.pid文件
// sameName 表示是否需要同名
// killOld 表示是否杀掉老进程
func KeepSingleProcess(key string, sameName bool, killOld bool) error {
	binPath := ucore.BinPath()
	pidPath := filepath.Join(filepath.Dir(binPath), "."+key+".pid")
	pidStr, err := dry.FileGetString(pidPath)
	if err == nil {
		//文件存在
		pid, _ := strconv.Atoi(pidStr)
		proc, err := process.NewProcess(int32(pid))
		if err == nil { //有进程在运行
			killOrReturn := func() error {
				if killOld {
					ctx := context.TODO()
					err = proc.KillWithContext(ctx)
					<-ctx.Done()
					return err
				}
				return errors.Errorf("已有进程在运行 pid=%d", pid)
			}
			if !sameName {
				return killOrReturn()
			}
			cmds, _ := proc.CmdlineSlice()
			if cmds[0] == binPath {
				return killOrReturn()
			}
		}
	}
	// 运行运行，记录pid
	return dry.FileSetString(pidPath, fmt.Sprintf("%d", os.Getpid()))
}
