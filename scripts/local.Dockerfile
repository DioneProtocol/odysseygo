# syntax=docker/dockerfile:experimental

# This Dockerfile is meant to be used with the build_local_dep_image.sh script
# in order to build an image using the local version of coreth

# Changes to the minimum golang version must also be replicated in
# scripts/build_dione.sh
# scripts/local.Dockerfile (here)
# Dockerfile
# README.md
# go.mod
FROM golang:1.18.5-buster

RUN mkdir -p /go/src/github.com/dioneprotocol

WORKDIR $GOPATH/src/github.com/dioneprotocol
COPY dionego dionego

WORKDIR $GOPATH/src/github.com/dioneprotocol/dionego
RUN ./scripts/build_dione.sh

RUN ln -sv $GOPATH/src/github.com/dioneprotocol/dione-byzantine/ /dionego
