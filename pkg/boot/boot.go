/*
	MIT License

Copyright (c) 2020 xzx.weiliang

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/
package boot

import (
	"fmt"
	"github.com/containerd/cgroups"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/sirupsen/logrus"
	"github.com/weiliang-ms/easyctl/pkg/util/command"
	strings2 "github.com/weiliang-ms/easyctl/pkg/util/strings"
	"gopkg.in/yaml.v2"
	"os/exec"
)

type Config struct {
	BootApp []DeSerializationAppItem `yaml:"boot-app"`
}

type DeSerializationAppItem struct {
	AppName   string `yaml:"app-name"`
	BootCmd   string `yaml:"boot-cmd"`
	Resources struct {
		Limits struct {
			Cpu    int    `yaml:"cpu"`
			Memory string `yaml:"memory"`
		} `yaml:"limits"`
	} `yaml:"resources"`
}

type Apps struct {
	Items []AppItem
}

type AppItem struct {
	AppName        string
	BootCmd        string
	Pid            int
	CpuCores       int
	LimitMem       *int64
	LimitMemOrigin string
	Logger         *logrus.Logger
}

// AppWithCGroups 以控制组限制的方式启动程序
func AppWithCGroups(item command.OperationItem) command.RunErr {
	config := Config{}
	if err := yaml.Unmarshal(item.B, &config); err != nil {
		return command.RunErr{Err: err, Msg: "反序列化失败"}
	}

	apps := Apps{}

	for _, v := range config.BootApp {
		mem, _ := strings2.GetMemoryBytes(v.Resources.Limits.Memory)
		apps.Items = append(apps.Items, AppItem{
			AppName:        v.AppName,
			BootCmd:        v.BootCmd,
			CpuCores:       v.Resources.Limits.Cpu,
			LimitMem:       &mem,
			LimitMemOrigin: v.Resources.Limits.Memory,
			Logger:         item.Logger,
		})
	}

	for _, v := range apps.Items {
		if err := v.bootProcess(); err.Err != nil {
			return err
		}
	}

	return command.RunErr{}
}

// 根据命令启动进程 -> 返回进程id
func (item *AppItem) bootProcess() command.RunErr {
	if item.AppName == "" {
		return command.RunErr{Err: fmt.Errorf("应用名: %s非法", item.AppName)}
	}

	cmd := exec.Command("/bin/sh", "-c", item.BootCmd)

	if err := cmd.Start(); err != nil {
		return command.RunErr{Err: err}
	}

	// 获取进程ID
	item.Pid = cmd.Process.Pid
	item.Logger.Infof("启动命令: %s, 进程id: %d", item.BootCmd, item.Pid)

	// 限制配置
	item.Logger.Infof("限制程序配额 -> CPU: %d核, 内存: %s", item.CpuCores, item.LimitMemOrigin)

	// 0.1s
	period := uint64(100000)

	item.Logger.Infof("创建cpu子系统: /sys/fs/cgroup/cpu/%s "+
		"memory子系统: /sys/fs/cgroup/memory/%s", item.AppName, item.AppName)

	// 初始化内存配额
	memory := &specs.LinuxMemory{}
	if *item.LimitMem == 0 {
		*memory.Limit = 9223372036854771712
	} else {
		memory.Limit = item.LimitMem
	}

	// 初始化cpu配额
	//var cpu *specs.LinuxCPU
	defaultQuota := int64(-1)
	cpu := &specs.LinuxCPU{Period: &period, Quota: &defaultQuota}

	quota := int64(100000 * item.CpuCores)
	//cpu.Period = &period
	if item.CpuCores != 0 {
		cpu.Quota = &quota
	}

	item.Logger.Debug(item.CpuCores)
	item.Logger.Debugf("Quota: %d Period: %d", quota, period)

	// 创建控制组，cpu、内存赋值
	control, err := cgroups.New(cgroups.V1, cgroups.StaticPath(fmt.Sprintf("/%s", item.AppName)), &specs.LinuxResources{
		CPU:    cpu,
		Memory: memory,
	})

	if err != nil {
		return command.RunErr{Err: err}
	}

	// 将进程添加至控制组
	if err := control.Add(cgroups.Process{Pid: item.Pid}); err != nil {
		return command.RunErr{Err: err}
	}

	return command.RunErr{}
}
