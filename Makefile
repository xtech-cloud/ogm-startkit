APP_NAME := ogm-startkit
BUILD_VERSION   := $(shell git tag --contains)
BUILD_TIME      := $(shell date "+%F %T")
COMMIT_SHA1     := $(shell git rev-parse HEAD )
GOOGLEAPIS_DIR  := /usr/local/include/googleapis

.PHONY: build
build: proto
	go build -ldflags \
		"\
		-X 'main.BuildVersion=${BUILD_VERSION}' \
		-X 'main.BuildTime=${BUILD_TIME}' \
		-X 'main.CommitID=${COMMIT_SHA1}' \
		"\
		-o ./bin/${APP_NAME}

.PHONY: proto
proto:
	protoc --proto_path=. --micro_out=. --go_out=. proto/startkit/healthy.proto
	#protoc -I/usr/local/include/googleapis -I./proto --include_imports --include_source_info --descriptor_set_out=./proto/startkit.pb proto/startkit/healthy.proto
	#/mnt/c/_wsl/protoc.exe --proto_path=./ --csharp_out=./proto/startkit --grpc_out=./proto/startkit --plugin=protoc-gen-grpc=c:/_wsl/grpc_csharp_plugin.exe proto/startkit/echo.proto
	#/mnt/c/_wsl/protoc.exe -I=./proto --js_out=import_style=typescript:./proto/startkit --grpc-web_out=import_style=typescript,mode=grpcwebtext:./proto proto/startkit/echo.proto
	#protoc --proto_path=./ --java_out=./proto/startkit proto/startkit/echo.proto
	#protoc --proto_path=./ --grpc-java_out=./proto/startkit --plugin=protoc-gen-grpc-java=/usr/bin/protoc-gen-grpc-java proto/startkit/echo.proto
	#mv proto/startkit/omo/msa/startkit/* proto/startkit/
	#rm -rf proto/startkit/omo

.PHONY: run
run:
	./bin/${APP_NAME}

.PHONY: run-fs
run-fs:
	MSA_CONFIG_DEFINE='{"source":"file","prefix":"/etc/ogm/","key":"startkit.yml"}' ./bin/${APP_NAME}

.PHONY: run-cs
run-cs:
	MSA_CONFIG_DEFINE='{"source":"consul","prefix":"/xtc/ogm/config","key":"startkit.yml"}' ./bin/${APP_NAME}

.PHONY: call
call:
	MICRO_REGISTRY=etcd micro call xtc.api.ogm.startkit Healthy.Echo '{"msg":"hello"}'

.PHONY: post
post:
	curl -X POST -d '{"msg":"hello"}' 127.0.0.1:18800/ogm/startkit/Healthy/Echo

.PHONY: tester
tester:
	go build -o ./bin/ ./tester

.PHONY: benchmark
benchmark:
	python3 ./benchmark.py

.PHONY: dist
dist:
	rm -rf ./dist
	mkdir ./dist
	tar -zcf dist/${APP_NAME}-${BUILD_VERSION}.tar.gz ./bin/${APP_NAME}
