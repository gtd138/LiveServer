package server_conf

// 单件
var instance *ServerConfig
var bReadConfig bool

// 服务器配置元素
type ConfigElem struct {
	Id         int    // 服务器id
	Type       string // 类型
	Host       string // 地址
	Port       string // rpc端口
	ClientPort string // 客户端端口
	Fronted    bool   // 是否为前端
}

// 服务器配置
type ServerConfig struct {
	// 前端服务器
	Gate      []ConfigElem // 网关服
	Connector []ConfigElem // 连接服

	// 后端服务器
	Master   []ConfigElem // 中心管理服
	Game     []ConfigElem // 游戏逻辑服
	DataBase []ConfigElem // 数据库缓存服
	Society  []ConfigElem // 社会服
	Lobby    []ConfigElem // 大厅服
	Auth     []ConfigElem // 第三方验证服

	ConfigMap map[string][]ConfigElem // 配置map
}

func (this *ServerConfig) ConvertToMap() {
	if this.ConfigMap != nil {
		return
	}
	this.ConfigMap = make(map[string][]ConfigElem)
	this.ConfigMap["gate"] = this.Gate
	this.ConfigMap["connector"] = this.Connector
	this.ConfigMap["master"] = this.Master
	this.ConfigMap["game"] = this.Game
	this.ConfigMap["database"] = this.DataBase
	this.ConfigMap["society"] = this.Society
	this.ConfigMap["lobby"] = this.Lobby
	this.ConfigMap["auth"] = this.Auth
}

// 获取服务器配置
func (this *ServerConfig) GetConfig(server_type string, server_id int) (conf *ConfigElem) {
	conf_list, bok := this.ConfigMap[server_type]
	if !bok {
		return
	}
	for i, v := range conf_list {
		if v.Id == server_id {
			conf = &(conf_list[i])
			break
		}
	}
	return
}

func GetSingleton() *ServerConfig {
	if instance == nil {
		instance = new(ServerConfig)
	}
	return instance
}

func IsReadConfig() bool {
	return bReadConfig
}

func ReadConfig(bread bool) {
	bReadConfig = bread
}
