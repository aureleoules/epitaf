#!/bin/sh

sleep 2;
export PGPASSWORD=root; 
createdb -h db_test -U root core;
psql --single-transaction -h db_test -U root -d core -f /app/db.sql
go test -v ./models && go test -v ./api/... && go test -v ./utils

