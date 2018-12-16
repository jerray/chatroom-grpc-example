NAME = chatroom
PACKAGE = github.com/jerray/chatroom
MAIN = $(PACKAGE)/main

DEFAULT_TAG = omnipay:latest
DEFAULT_BUILD_TAG = 1.10-alpine
REMOTE_IMAGE = ccr.ccs.tencentyun.com/fpay/omnipay

BUILD_FLAGS= -mod vendor -v -o $(NAME) main/main.go

CL_RED  = "\033[0;31m"
CL_BLUE = "\033[0;34m"
CL_GREEN = "\033[0;32m"
CL_ORANGE = "\033[0;33m"
CL_NONE = "\033[0m"

define color_out
	@echo $(1)$(2)$(CL_NONE)
endef

build:
	@go mod vendor
	@go build $(BUILD_FLAGS)

proto:
	# If build proto failed, make sure you have protoc installed and:
	# go get -u github.com/google/protobuf
	# go get -u github.com/golang/protobuf/protoc-gen-go
	# go install github.com/mwitkow/go-proto-validators/protoc-gen-govalidators
	# mkdir -p $GOPATH/src/github.com/googleapis && git clone git@github.com:googleapis/googleapis.git $GOPATH/src/github.com/googleapis/
	@mkdir -p pb
	@protoc \
		--proto_path=${GOPATH}/src \
		--proto_path=${GOPATH}/src/github.com/google/protobuf/src \
		--proto_path=${GOPATH}/src/github.com/googleapis/googleapis \
		--proto_path=. \
		--go_out=plugins=grpc:$(PWD)/pb \
		--govalidators_out=$(PWD)/pb \
 		chatroom.proto
	$(call color_out,$(CL_ORANGE),"Done")

.PHONY: all
all:
	build
