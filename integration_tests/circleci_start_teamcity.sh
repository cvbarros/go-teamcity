#!/bin/bash -eo pipefail
docker create -v /data/teamcity_server/datadir --name data alpine:3.4 /bin/true
test -d ${TEAMCITY_DATA_DIR} || tar xfz ${INTEGRATION_TEST_DIR}/teamcity_data.tar.gz -C ${INTEGRATION_TEST_DIR}
docker cp ${INTEGRATION_TEST_DIR}/data_dir/. data:/data/teamcity_server/datadir
docker run  --rm -d -p 8112:8111 \
            --volumes-from data \
            jetbrains/teamcity-server:${TEAMCITY_VERSION}
echo -n "Teamcity server is booting (this may take a while)..."
until $(curl -o /dev/null -sfI $TEAMCITY_HOST/login.html);do echo -n ".";sleep 5;done