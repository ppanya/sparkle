

staging-deploy:
	sh scripts/staging-deploy.sh

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./dist/sparkle-linux-x64 ./cmd/gcp_compute/main.go

backingsvc-up:
	sh scripts/backingservice-up.sh
backingsvc-down:
	sh scripts/backingservice-down.sh

gen-proto:
	prototool generate
	protoc-go-inject-tag -XXX_skip=gorm,bson -input=./pkg/entities/v1/user.pb.go
	protoc-go-inject-tag -XXX_skip=gorm,bson -input=./pkg/entities/v1/identity.pb.go
	protoc-go-inject-tag -XXX_skip=gorm,bson -input=./pkg/entities/v1/session.pb.go

gen-statik:
	@statik -src=./resources -dest=./cmd/statik -f
