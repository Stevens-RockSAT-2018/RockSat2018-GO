export GOOS=linux
export GOARCH=arm
go build main.go
scp main debian@192.168.1.33:/home/debian/main
