module easyctl

go 1.15

require (
	github.com/modood/table v0.0.0-20200225102042-88de94bb9876
	github.com/pkg/sftp v1.12.0
	github.com/spf13/cobra v1.0.0
	github.com/vbauerster/mpb/v6 v6.0.3
	github.com/weiliang-ms/easyctl v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2
	golang.org/x/sys v0.0.0-20210331175145-43e1dd70ce54 // indirect
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/weiliang-ms/easyctl => ../easyctl
