#/bin/bash

ProtoRootPath=$(cd $(dirname $0); pwd)

pushd $ProtoRootPath
protoc ./*.proto  --proto_path=.  --proto_path=${ProtoRootPath} --gogofast_out=plugins=grpc:.
popd
