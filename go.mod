module easyctl

go 1.15

require (
	github.com/jteeuwen/go-bindata v3.0.7+incompatible // indirect
	github.com/pkg/sftp v1.12.0
	github.com/spf13/cobra v1.0.0
	github.com/weiliang-ms/easyctl v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.0.0-20201012173705-84dcc777aaee
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/weiliang-ms/easyctl => ../easyctl
