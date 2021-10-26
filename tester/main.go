package main

import (
	"context"
	"io"
	"time"

	"github.com/asim/go-micro/plugins/client/grpc/v3"
	_ "github.com/asim/go-micro/plugins/registry/etcd/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/client"
	"github.com/asim/go-micro/v3/logger"
	proto "ogm-startkit/proto/startkit"
)

func main() {
	service := micro.NewService(
		micro.Client(grpc.NewClient()),
		micro.Name("tester"),
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

	healthy := proto.NewHealthyService("xtc.ogm.startkit", cli)

	logger.Trace("----------------------------------------------------------")
	// Call
	{
		rsp, err := healthy.Echo(context.Background(), &proto.Request{
			Msg: time.Now().String() + " | hello",
		})
		if err != nil {
			logger.Error(err)
			return
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
