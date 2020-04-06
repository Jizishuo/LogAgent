package config

// KafkaCnf 配置
type KafkaConf struct {
	Address string `ini:"address"`
	Topic   string `ini:"topic"`
	MaxChan int    `ini:"maxchan"`
}

// TaillogCof配置
type TaillogCof struct {
	FilePth string `ini:"filePath"`
}

// etcd配置
type Etcd struct {
	Address string `ini:"address"`
	Timeout int    `ini:"timeout"`
	Key     string `ini:"key"`
}

// elasticsearch配置
type ElasticSearch struct {
	Address string `ini:"address"`
	MaxChan int    `ini:"maxchan"`
	Index   string `ini:"index"`
	Type    string `ini:"type"`
}

// conf 配置
type Conf struct {
	KafkaConf     `ini:"kafka"`
	TaillogCof    `ini:"taillog"`
	Etcd          `ini:"etcd"`
	ElasticSearch `ini:"elasticsearch"`
}
