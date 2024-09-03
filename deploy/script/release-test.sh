#!/bin/bash
need_start_service_shell=(
    # rpc
    user-rpc-test.sh

    # api
    user-api-test.sh
)

for i in ${need_start_service_shell[*]} ; do
    chmod +x $i
    ./$i
done


docker ps

docker exec -it etcd etcdctl get --prefix ""