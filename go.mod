module ogm-startkit

go 1.16

require (
	github.com/asim/go-micro/plugins/client/grpc/v3 v3.0.0-20210630062103-c13bb07171bc
	github.com/asim/go-micro/plugins/config/encoder/yaml/v3 v3.0.0-20210721080634-e1bc7e302871
	github.com/asim/go-micro/plugins/config/source/etcd/v3 v3.0.0-20210721080634-e1bc7e302871
	github.com/asim/go-micro/plugins/logger/logrus/v3 v3.0.0-20210721080634-e1bc7e302871
	github.com/asim/go-micro/plugins/registry/etcd/v3 v3.0.0-20210721080634-e1bc7e302871
	github.com/asim/go-micro/plugins/server/grpc/v3 v3.0.0-20210721080634-e1bc7e302871
	github.com/asim/go-micro/v3 v3.5.2
	github.com/sirupsen/logrus v1.8.1
	google.golang.org/genproto v0.0.0-20210721163202-f1cecdd8b78a
	google.golang.org/protobuf v1.27.1
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)
