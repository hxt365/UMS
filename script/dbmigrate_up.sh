#!/bin/bash

MYSQL_URL=$1

migrate -database ${MYSQL_URL} -path db/migrations up
