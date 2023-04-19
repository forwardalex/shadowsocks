TAG=v0.0.1
Serve=socks_server
export ENV_NAME=dev
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64
export MASTER=on
go build -gcflags "all=-N -l" -o $Serve
docker build -t $Serve:$TAG .
docker tag $Serve:$TAG 10.0.8.3:5000/$Serve:$TAG
docker push 10.0.8.3:5000/$Serve:$TAG



