#!/bin/bash
echo 'generate db model'


go run modelgenerator.go \
-tplFile='./model.tpl' \
-modelFolder='../model/' \
-packageName='model' \
-dbIP='127.0.0.1'  \
-dbPort=3306 \
-dbConnection='dbhelper.DB' \
-dbName='dbnote' \
-userName='root' \
-pwd='123456' \
-genTable='mail#msg#notice' \

echo 'done'

