package main

import "flag"

type Config struct {
	ESAddr    string
	MysqlAddr string
	MysqlDSM  string
	Base      string
}

var GLOBAL_CONFIG = Config{}

func init() {
	h := flag.String("h", "127.0.0.1", "Global host")
	md := flag.String("m", "root:123456@tcp(127.0.0.1:3306)/git_backend?charset=utf8mb4&parseTime=True&loc=Local", "Elasticesearch address and port")
	GLOBAL_CONFIG.Base = *h
	GLOBAL_CONFIG.ESAddr = "http://" + *h + ":9200"
	GLOBAL_CONFIG.MysqlDSM = *md
}
