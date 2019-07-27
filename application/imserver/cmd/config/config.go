package config

import "micro-me/application/common/config"

type (
	ImConfig struct {
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
		RabbitMq *config.RabbitMq
	}

	ImRpcConfig struct {
		Version string
		Topic   string
		Server  struct {
			Name      string
			RateLimit int64
		}
		Etcd struct {
			Address  []string
			UserName string
			Password string
		}

		ImServerList []*config.ImRpc
	}
)
