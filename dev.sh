#!/usr/bin/env sh
P2SUB_SOURCE=$(pwd)
P2SUB_GOPATH=$GOPATH/src/github.com/p2sub
go get github.com/btcsuite/btcutil/base58
mkdir -p $P2SUB_GOPATH
P2SUB_SYMLINK=$P2SUB_GOPATH/p2sub
if [[ -L "$P2SUB_SYMLINK" ]]; then
  rm "$P2SUB_SYMLINK"
fi
if [[ ! -L "$P2SUB_SYMLINK" ]]; then
  echo "Linked $P2SUB_SYMLINK -> $P2SUB_SOURCE"
  ln -s $(pwd) $P2SUB_SYMLINK
fi
go run ./cmd/p2sub/p2sub.go