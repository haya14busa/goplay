# https://github.com/golang/go/wiki/HostedContinuousIntegration
# https://loads.pickle.me.uk/2015/08/22/easy-peasy-github-releases-for-go-projects-using-travis/
package = github.com/haya14busa/goplay/cmd/goplay

.PHONY: release

release:
	mkdir -p release
	GOOS=linux   GOARCH=amd64 go build -o release/goplay-linux-amd64       $(package)
	GOOS=linux   GOARCH=386   go build -o release/goplay-linux-386         $(package)
	GOOS=linux   GOARCH=arm   go build -o release/goplay-linux-arm         $(package)
	GOOS=windows GOARCH=amd64 go build -o release/goplay-windows-amd64.exe $(package)
	GOOS=windows GOARCH=386   go build -o release/goplay-windows-386.exe   $(package)
	GOOS=darwin  GOARCH=amd64 go build -o release/goplay-darwin-amd64      $(package)
	GOOS=darwin  GOARCH=386   go build -o release/goplay-darwin-386        $(package)
