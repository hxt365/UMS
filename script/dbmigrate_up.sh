#!/bin/bash

MYSQL_URL=$1

migrate -database ${MYSQL_URL} -path storage/dbmigrations up
