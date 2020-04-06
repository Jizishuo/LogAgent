package config

// KafkaCnf 配置
type KafkaConf struct {
	Address string `ini:"address"`
	Topic   string `ini:"topic"`
	MaxChan int    `ini:"maxchan"`
}

type TaillogCof struct {
	FilePth string `ini:"filePath"`
}

type Etcd struct {
	Address string `ini:"address"`
	Timeout int    `ini:"timeout"`
	Key     string `ini:"key"`
}

type ElasticSearch struct {
	Address string `ini:"address"`
	MaxChan int    `ini:"maxchan"`
	Index   string `ini:"index"`
	Type    string `ini:"type"`
}

type Conf struct {
	KafkaConf     `ini:"kafka"`
	TaillogCof    `ini:"taillog"`
	Etcd          `ini:"etcd"`
	ElasticSearch `ini:"elasticsearch"`
}
