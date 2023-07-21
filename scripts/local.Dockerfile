# syntax=docker/dockerfile:experimental

# This Dockerfile is meant to be used with the build_local_dep_image.sh script
# in order to build an image using the local version of coreth

# Changes to the minimum golang version must also be replicated in
# scripts/build_odyssey.sh
# scripts/local.Dockerfile (here)
# Dockerfile
# README.md
# go.mod
FROM golang:1.19.6-buster

RUN mkdir -p /go/src/github.com/DioneProtocol

WORKDIR $GOPATH/src/github.com/DioneProtocol
COPY odysseygo odysseygo

WORKDIR $GOPATH/src/github.com/DioneProtocol/odysseygo
RUN ./scripts/build_odyssey.sh

RUN ln -sv $GOPATH/src/github.com/DioneProtocol/odyssey-byzantine/ /odysseygo
