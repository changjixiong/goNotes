#!/bin/bash
echo 'generate cassandra model'

go run cassandramodelgen.go \
-tplFile='./cassandramodel.tpl' \
-modelFolder='../modelcassandra/' \
-packageName='cassandra' \
-dbIP='192.168.199.157'  \
-dbConnection='dbhelper.DBCassandra' \
-dbName='contra_stats_s999_2' \
-userName='test' \
-pwd='123456' \
-genTable='player_hero_use_log' \

echo 'done'
