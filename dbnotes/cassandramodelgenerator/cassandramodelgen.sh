#!/bin/bash
echo 'generate cassandra model'

go run cassandramodelgen.go \
-tplFile='./cassandramodel.tpl' \
-modelFolder='../modelcassandra/' \
-packageName='cassandra' \
-dbIP='127.0.0.1'  \
-dbConnection='dbhelper.DBCassandra' \
-dbName='space_for_back' \
-dbPort=9042 \
-userName='test' \
-pwd='123456' \
-genTable='num_log_10ms#num_log_100ms#num_log_2s#num_log_4s' \

echo 'done'
