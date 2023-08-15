package config

type HarborConnectInfoInterface interface {
	GetConnectInfo(string) (*HarborConnectInfo, error) //获取连接信息
	SetConnectInfo(*HarborConnectInfo)                 //设置连接信息
	DelConnectInfo(string)                             //删除连接信息
	SetConnectInfoContext(string)                      //设置当前上下文
	SetConnectInfoAlias(string, string) error
	GetDefaultConnectInfo() *HarborConnectInfo
	View()
	PrintCurrentContext(bool)
	List(string, bool, bool)
}
