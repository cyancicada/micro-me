package config

type (
	UserRpcServer struct {
		ClientName string
		ServerName string
	}
	RabbitMq struct {
		Address []string
		Topic   string
	}

	// im_address,topic,server_name
	ImRpc struct {
		Address     string // 这是一个真正 的ip地址
		AmqbAddress []string
		Topic       string
		ServerName  string
	}
)
