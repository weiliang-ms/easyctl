package cmd

var serverTmpl []byte

//
//func init() {
//	InitTmplCmd.AddCommand(initServerTmplCmd)
//	InitTmplCmd.AddCommand(initKeepaliveTmplCmd)
//	InitTmplCmd.AddCommand(initHATmplCmd)
//	InitTmplCmd.AddCommand(initDockerTmplCmd)
//}
//
//// InitTmplCmd 初始化模板命令
//var InitTmplCmd = &cobra.Command{
//	Use:     "init-tmpl [OPTIONS] [flags]",
//	Short:   "初始化配置模板指令集",
//	Example: "\neasyctl init-tmpl server\n",
//	Args:    cobra.MinimumNArgs(1),
//}
//
//// init servers template
//var initServerTmplCmd = &cobra.Command{
//	Use: "server [flags]",
//	Run: func(cmd *cobra.Command, args []string) {
//		tmpl, _ := asset.Asset("static/tmpl/server.yaml")
//		ioutil.WriteFile(fmt.Sprintf("%s/server.yaml", currentPath()), tmpl, 0644)
//	},
//	Args: cobra.NoArgs,
//}
//
//// init keepalive servers template
//var initKeepaliveTmplCmd = &cobra.Command{
//	Use: "keepalived [flags]",
//	Run: func(cmd *cobra.Command, args []string) {
//		tmpl, _ := asset.Asset("static/tmpl/keepalived.yaml")
//		ioutil.WriteFile(fmt.Sprintf("%s/keepalived.yaml", currentPath()), tmpl, 0644)
//	},
//	Args: cobra.NoArgs,
//}
//
//// init keepalive servers template
//var initHATmplCmd = &cobra.Command{
//	Use: "haproxy [flags]",
//	Run: func(cmd *cobra.Command, args []string) {
//		tmpl, _ := asset.Asset("static/tmpl/haproxy.yaml")
//		ioutil.WriteFile(fmt.Sprintf("%s/haproxy.yaml", currentPath()), tmpl, 0644)
//	},
//	Args: cobra.NoArgs,
//}
//
//// init docker servers template
//var initDockerTmplCmd = &cobra.Command{
//	Use: "docker [flags]",
//	Run: func(cmd *cobra.Command, args []string) {
//		tmpl, _ := asset.Asset("static/tmpl/docker.yaml")
//		ioutil.WriteFile(fmt.Sprintf("%s/docker.yaml", currentPath()), tmpl, 0644)
//	},
//	Args: cobra.NoArgs,
//}
//
//func currentPath() string {
//	re := runner.Shell("pwd")
//	if re.ExitCode != 0 {
//		log.Fatal(errors.New(fmt.Sprintf("获取当前路径失败：%s", re.StdErr)))
//	}
//
//	return re.StdOut
//}
