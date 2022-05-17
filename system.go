package golog

import (
	"os"
)

type System struct {
	Variable []string // 环境变量
	Hostname string   // 主机名
	Twd      string   // 当前目录
	Uid      int      // 用户ID
	EUid     int      // 有效用户ID
	Gid      int      // 组ID
	EGid     int      // 有效组ID
	Pid      int      // 进程ID
	PPid     int      // 父进程ID
}

func (s *System) Init() *System {
	s.Variable = os.Environ()
	s.Hostname, _ = os.Hostname()
	s.Twd, _ = os.Getwd()
	s.Uid = os.Getuid()
	s.EUid = os.Geteuid()
	s.Gid = os.Getgid()
	s.EGid = os.Getegid()
	s.Pid = os.Getpid()
	s.PPid = os.Getppid()
	return s
}
