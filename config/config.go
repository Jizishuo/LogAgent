package config

// KafkaCnf 配置
type KafkaConf struct {
	Address string `ini:"address"`
	Topic string `ini:"topic"`
	MaxChan int `ini:"maxchan"`
}




