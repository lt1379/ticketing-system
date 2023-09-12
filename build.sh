$env:GOOS = "linux"
$env:GOARCH= "amd64"
ENV=testing go test -v usecase/user_usecase_test.go

go test --coverprofile=cover.txt ./...
go tool cover -html=cover.txt
mockery --dir=domain/repository --name=IUser --filename=user.go --output=domain/mocks/repomocks --outpkg=repomocks
export AZURE_SERVICEBUS_HOSTNAME=Endpoint=sb://gra-sb-gateway.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=3k7vB9LCiVW5zWtpyJoGHmLTy4SrLPsk6+ASbNdCwpU=
go build main.go