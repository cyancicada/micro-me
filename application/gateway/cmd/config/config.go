package config

import "micro-me/application/common/config"

type (
	ApiConfig struct {
		Version string
		Port    string
		Server  struct {
			Name      string
			RateLimit int64
		}
		Etcd struct {
			Address  []string
			UserName string
			Password string
		}
		Engine struct {
			Name       string
			DataSource string
		}

		UserRpcServer *config.UserRpcServer
		ImRpcServer   struct {
			ServerName   string
			ImServerList []*config.ImRpc
		}
	}
)
