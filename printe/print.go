package printe

import (
	"fmt"
	"github.com/weiliang-ms/easyctl/constant"
	"github.com/weiliang-ms/easyctl/util"
	"os"
)

func PackageDetectionPass(banner string) {
	fmt.Printf("%s 依赖检测通过...\n", util.PrintGreen(banner))
}
func OriginPackageDetectionPass(banner string, instance util.SSHInstance) {
	fmt.Printf("%s 远程依赖检测通过...\n", util.PrintGreenMulti([]string{banner, constant.Origin, instance.Host}))
}

func PackageInstall() {
	fmt.Printf("%s 尝试安装包...\n", util.PrintGreen(constant.Install))
}

func InstallPackageFatal() {
	fmt.Printf("%s 安装包失败，请检测yum...\n", util.PrintRed(constant.Error))
	os.Exit(1)
}

func PackageOriginInstall(instance util.SSHInstance) {
	fmt.Printf("%s 尝试远程主机：%s安装包...\n", util.PrintGreen(constant.Origin), instance.Host)
}
