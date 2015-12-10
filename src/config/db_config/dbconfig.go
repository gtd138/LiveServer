package dbconfig

// 配置元素
type DBConfigElem struct {
	Name string // 数据库名
	Host string // 数据库地址
}

// 数据库配置
type DBConfig struct {
	Config []DBConfigElem
}

func (this *DBConfig) GetConfig(name string) (config DBConfigElem) {
	for _, v := range this.Config {
		if name == v.Name {
			config = v
			break
		}
	}
	return
}
