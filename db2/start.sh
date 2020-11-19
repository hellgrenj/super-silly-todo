#!/bin/bash
echo "sleeping 10 seconds and then running migrations" 
sleep 10 

echo "running migrations"
migrate -database 'mysql://root:example@(mysqldb:3306)/silly' -path migrations up
