#!/bin/bash
CURRENT_DIR=$1
for x in $(find ${CURRENT_DIR}/invan_proto/* -type d); do
  protoc -I=${x} -I=${CURRENT_DIR}/invan_proto -I /usr/local/include \
	 --go_out=${CURRENT_DIR} \
   --go-grpc_opt=require_unimplemented_servers=false \
   --go-grpc_out=${CURRENT_DIR} ${x}/*.proto
done