package install

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

/*
	Test Error
*/

var (
	mockHost = "1.1.1.1"
	mockErr  = fmt.Errorf("mock err")
)

func Test_BootErr(t *testing.T) {
	require.Equal(t, fmt.Sprintf("[%s] 启动失败 -> %s", mockHost, mockErr), BootErr{
		Host: mockHost,
		Err:  mockErr,
	}.Error())
}

func Test_SetSystemdErr(t *testing.T) {
	require.Equal(t, fmt.Sprintf("[%s] 配置systemd失败 -> %s", mockHost, mockErr), SetSystemdErr{
		Host: mockHost,
		Err:  mockErr,
	}.Error())
}

func Test_SetConfigErr(t *testing.T) {
	require.Equal(t, fmt.Sprintf("[%s] 配置docker失败 -> %s", mockHost, mockErr), SetConfigErr{
		Host: mockHost,
		Err:  mockErr,
	}.Error())
}

func Test_SetUpRuntimeErr(t *testing.T) {
	require.Equal(t, fmt.Sprintf("[%s] 配置docker运行时失败 -> %s", mockHost, mockErr), SetUpRuntimeErr{
		Host: mockHost,
		Err:  mockErr,
	}.Error())
}

func Test_InstallErr(t *testing.T) {
	require.Equal(t, fmt.Sprintf("[%s] 安装docker失败 -> %s", mockHost, mockErr), InstallErr{
		Host: mockHost,
		Err:  mockErr,
	}.Error())
}

func Test_PruneErr(t *testing.T) {
	require.Equal(t, fmt.Sprintf("[%s] 清理docker失败 -> %s", mockHost, mockErr), PruneErr{
		Host: mockHost,
		Err:  mockErr,
	}.Error())
}

func Test_DetectErr(t *testing.T) {
	require.Equal(t, fmt.Sprintf("[%s] 检测docker安装依赖失败 -> %s", mockHost, mockErr), DetectErr{
		Host: mockHost,
		Err:  mockErr,
	}.Error())
}

func Test_ParseServerListErr_Err(t *testing.T) {
	require.Equal(t, fmt.Sprintf("反序列化主机列表失败 -> %s", mockErr), ParseServerListErr{Err: mockErr}.Error())
}

func Test_TransferPackageErr_Err(t *testing.T) {
	require.Equal(t, fmt.Sprintf("传输%s:/tmp/%s失败 -> %s",
		mockHost,
		"ddd",
		mockErr.Error()), TransferPackageErr{
		Host: mockHost,
		Path: "ddd",
		Err:  mockErr,
	}.Error())
}
