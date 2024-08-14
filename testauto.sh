#!/bin/bash

# docker compose up
# psql -U postgres -d gmart

# go build -o ./cmd/gophermart/gophermart ./cmd/gophermart/*.go
# postgresql://postgres:postgres@0.0.0.0:5432/gmart?sslmode=disable
# ./cmd/accrual/accrual_linux_amd64 -a localhost:8081 -d postgresql://postgres:postgres@0.0.0.0:5432/gmart?sslmode=disable
# ./gophermart -a localhost:8080 -d postgresql://postgres:postgres@0.0.0.0:5432/gmart?sslmode=disable -r localhost:8081

~/bin/gophermarttest \
            -test.v -test.run=^TestGophermart$ \
            -gophermart-binary-path=cmd/gophermart/gophermart \
            -gophermart-host=localhost \
            -gophermart-port=8080 \
            -gophermart-database-uri="postgresql://postgres:postgres@0.0.0.0:5432/gmart?sslmode=disable" \
            -accrual-binary-path=cmd/accrual/accrual_linux_amd64 \
            -accrual-host=localhost \
            -accrual-port=8081 \
            -accrual-database-uri="postgresql://postgres:postgres@0.0.0.0:5432/gmart?sslmode=disable"