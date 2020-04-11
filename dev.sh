#!/usr/bin/env sh

case $1 in
  
  build)
    echo "Run build all project"
    go build -o ./build/configuration ./cmd/configuration/configuration.go
    go build -o ./build/p2sub ./cmd/p2sub/p2sub.go
    ;;

  config)
    echo "Run configuration, generate new config file in conf.d"
    go run ./cmd/configuration/configuration.go
    ;;

  master1)
    echo "Run master 1"
    go run ./cmd/p2sub/p2sub.go --config ./conf.d/master1.json
    ;;

  master2)
    echo "Run master 2"
    go run ./cmd/p2sub/p2sub.go --config ./conf.d/master2.json
    ;;

  master3)
    echo "Run master 3"
    go run ./cmd/p2sub/p2sub.go --config ./conf.d/master3.json
    ;;

  notary)
    echo "Run notary"
    go run ./cmd/p2sub/p2sub.go --config ./conf.d/notary.json
    ;;

  *)
    echo "This wasn't defined";
    exit 1
    ;;
esac