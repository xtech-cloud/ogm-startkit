package main

import (
	"context"
	"io"
	"time"

	"ogm-startkit/config"

	proto "ogm-startkit/proto/startkit"

	_ "github.com/asim/go-micro/plugins/registry/etcd/v3"
	"github.com/asim/go-micro/plugins/client/grpc/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/client"
	"github.com/asim/go-micro/v3/logger"
)

func main() {
	config.Setup()
	service := micro.NewService(
		micro.Client(grpc.NewClient()),
		micro.Name(config.Schema.Service.Name+".tester"),
	)
	service.Init()

	cli := service.Client()
	cli.Init(
		client.Retries(3),
		client.RequestTimeout(time.Second*1),
		client.Retry(func(_ctx context.Context, _req client.Request, _retryCount int, _err error) (bool, error) {
			if nil != _err {
				logger.Errorf("%v | [ERR] retry %d, reason is %v\n\r", time.Now().String(), _retryCount, _err)
				return true, nil
			}
			return false, nil
		}),
	)

	healthy := proto.NewHealthyService(config.Schema.Service.Name, cli)

	logger.Trace("----------------------------------------------------------")
	// Call
	{
		rsp, err := healthy.Echo(context.Background(), &proto.Request{
			Msg: time.Now().String() + " | OGM-StartKit",
		})
		if err != nil {
			logger.Error(err)
		} else {
			logger.Info(rsp.Msg)
		}
	}

	stream, err := healthy.PingPong(context.Background())
	if err != nil {
		logger.Error(err)
		return
	}
	defer stream.Close()
	stroke := int64(0)
	for range time.Tick(1 * time.Second) {
		stroke = stroke + 1
		stream.Send(&proto.Ping{Stroke: stroke})
		rsp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Error(err)
			continue
		}
		logger.Infof("Pong %v", rsp.Stroke)
	}
}
