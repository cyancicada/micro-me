module micro-me

go 1.12

require (
	github.com/SAP/go-hdb v0.14.1 // indirect
	github.com/StackExchange/wmi v0.0.0-20181212234831-e0a55b97c705 // indirect
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/coredns/coredns v1.4.0 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/gin-gonic/gin v1.4.0
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/go-sql-driver/mysql v1.4.1
	github.com/go-xorm/xorm v0.7.4
	github.com/golang/protobuf v1.3.2
	github.com/gorilla/websocket v1.4.0
	github.com/hashicorp/consul v1.5.2 // indirect
	github.com/hashicorp/go-gcp-common v0.5.0 // indirect
	github.com/hashicorp/go-memdb v1.0.0 // indirect
	github.com/hashicorp/go-plugin v1.0.0 // indirect
	github.com/hashicorp/hil v0.0.0-20190212132231-97b3a9cdfa93 // indirect
	github.com/hashicorp/raft-boltdb v0.0.0-20171010151810-6e5ba93211ea // indirect
	github.com/hashicorp/vault v1.1.0 // indirect
	github.com/hashicorp/vault-plugin-auth-alicloud v0.0.0-20190320211238-36e70c54375f // indirect
	github.com/hashicorp/vault-plugin-auth-azure v0.0.0-20190320211138-f34b96803f04 // indirect
	github.com/hashicorp/vault-plugin-auth-centrify v0.0.0-20190320211357-44eb061bdfd8 // indirect
	github.com/hashicorp/vault-plugin-auth-kubernetes v0.0.0-20190328163920-79931ee7aad5 // indirect
	github.com/hashicorp/vault-plugin-secrets-ad v0.0.0-20190327182327-ed2c3d4c3d95 // indirect
	github.com/hashicorp/vault-plugin-secrets-alicloud v0.0.0-20190320213517-3307bdf683cb // indirect
	github.com/hashicorp/vault-plugin-secrets-azure v0.0.0-20190320211922-2dc8a8a5e490 // indirect
	github.com/hashicorp/vault-plugin-secrets-gcp v0.0.0-20190320211452-71903323ecb4 // indirect
	github.com/hashicorp/vault-plugin-secrets-gcpkms v0.0.0-20190320213325-9e326a9e802d // indirect
	github.com/influxdata/influxdb v1.7.5 // indirect
	github.com/juju/ratelimit v1.0.1
	github.com/lyft/protoc-gen-validate v0.0.14 // indirect
	github.com/micro/go-micro v1.5.0
	github.com/micro/go-plugins v1.1.0
	github.com/nats-io/nats-server/v2 v2.0.0 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/shirou/gopsutil v2.18.12+incompatible // indirect
	github.com/sourcegraph/go-diff v0.5.1 // indirect
	github.com/ugorji/go v1.1.7 // indirect
	go.etcd.io/bbolt v1.3.2 // indirect
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4 // indirect
	golang.org/x/net v0.0.0-20190628185345-da137c7871d7 // indirect
	golang.org/x/sys v0.0.0-20190712062909-fae7ac547cb7 // indirect
	gopkg.in/go-playground/validator.v8 v8.18.2
	layeh.com/radius v0.0.0-20190322222518-890bc1058917 // indirect
)

replace (
	github.com/testcontainers/testcontainer-go => github.com/testcontainers/testcontainers-go v0.0.4
	golang.org/x/lint => github.com/golang/lint v0.0.0-20190409202823-959b441ac422
)

replace github.com/golang/lint => golang.org/x/lint v0.0.0-20190409202823-959b441ac422
