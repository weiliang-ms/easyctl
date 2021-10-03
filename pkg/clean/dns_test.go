package clean

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"strings"
	"testing"
)

const (
	pruneAll       = `echo "" > /etc/resolv.conf`
	pruneAllConfig = `
clean-dns:
  address-list:        # 地址列表，为空表示清除所有
  excludes:             # 排除哪些dns地址不被清理
`
	pruneSliceConfig = `
clean-dns:
  address-list:        # 地址列表，为空表示清除所有
    - 1.1.1.1
    - 2.2.2.2
  excludes:             # 排除哪些dns地址不被清理
`
	pruneSlice = `sed -i "/1.1.1.1/d" /etc/resolv.conf
sed -i "/2.2.2.2/d" /etc/resolv.conf`

	pruneFilterSliceConfig = `
clean-dns:
  address-list:        # 地址列表，为空表示清除所有
    - 1.1.1.1
    - 2.2.2.2
  excludes:             # 排除哪些dns地址不被清理
    - 1.1.1.1
`
	pruneFilterSlice = `sed -i "/2.2.2.2/d" /etc/resolv.conf`

	pruneAllExpectSliceConfig = `
clean-dns:
  address-list:        # 地址列表，为空表示清除所有
  excludes:             # 排除哪些dns地址不被清理
    - 1.1.1.1
`
	pruneAllExpectSlice = `while read line
do
   if [[ "$line" =~ 1.1.1.1 ]];then
 echo 过滤...
   else
 sed -i "/$line/d" /etc/resolv.conf
   fi
done < /etc/resolv.conf`
)

func TestPruneDnsScript(t *testing.T) {
	// 1.均为空
	config, err := PruneDnsScript([]byte(pruneAllConfig), pruneDnsShellTmpl)
	// 删除空白行
	re := regexp.MustCompile(`(?m)^\s*$[\r\n]*|[\r\n]+\s+\z`)
	assert.Nil(t, err)
	assert.Equal(t, pruneAll, strings.ReplaceAll(re.ReplaceAllString(config, ""), "\t", ""))

	// 2.exludes为空, address-list不为空
	configSlice, err := PruneDnsScript([]byte(pruneSliceConfig), pruneDnsShellTmpl)
	assert.Nil(t, err)
	assert.Equal(t, pruneSlice, re.ReplaceAllString(configSlice, ""))

	// 3.exludes不为空, address-list为空
	expect, err := PruneDnsScript([]byte(pruneAllExpectSliceConfig), pruneDnsShellTmpl)
	assert.Nil(t, err)
	assert.Equal(t, pruneAllExpectSlice, strings.ReplaceAll(re.ReplaceAllString(expect, ""), "\t", ""))

	// 4.均不为空
	filter, err := PruneDnsScript([]byte(pruneFilterSliceConfig), pruneDnsShellTmpl)
	assert.Nil(t, err)
	assert.Equal(t, pruneFilterSlice, strings.ReplaceAll(re.ReplaceAllString(filter, ""), "\t", ""))
}

func TestParseDnsConfig(t *testing.T) {
	const nilDnsConfig = `
clean-dns:
  address-list:
  excludes:
`
	const nilExcludeDnsConfig = `
clean-dns:
  address-list:
    - 114.114.114.114
    - 8.8.8.8
  excludes:
`
	const commonConfig = `
clean-dns:
  address-list:
    - 8.8.8.8
  excludes:
    - 114.114.114.114
`
	const badConfig = `
clean-dns:
  address-list: 8.8.8.8
  excludes:
    - 114.114.114.114
`
	var expectConfig, actureConfig DnsCleanerConfig
	actureConfig, err := ParseDnsConfig([]byte(nilDnsConfig))
	assert.Nil(t, err)
	assert.Equal(t, DnsCleanerConfig{}, actureConfig)

	actureConfig, err = ParseDnsConfig([]byte(nilExcludeDnsConfig))
	assert.Nil(t, err)
	expectConfig.CleanDns.AddressList = []string{"114.114.114.114", "8.8.8.8"}
	assert.Equal(t, expectConfig, actureConfig)

	actureConfig, err = ParseDnsConfig([]byte(commonConfig))
	assert.Nil(t, err)
	expectConfig.CleanDns.AddressList = []string{"8.8.8.8"}
	expectConfig.CleanDns.Excludes = []string{"114.114.114.114"}
	assert.Equal(t, expectConfig, actureConfig)

	actureConfig, err = ParseDnsConfig([]byte(badConfig))
	assert.NotNil(t, err)
}
