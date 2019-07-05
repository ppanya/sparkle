


backingsvc-up:
	sh scripts/backingservice-up.sh
backingsvc-down:
	sh scripts/backingservice-down.sh

gen-proto:
	prototool generate
	protoc-go-inject-tag -XXX_skip=gorm,bson -input=./pkg/entities/v1/user.pb.go

