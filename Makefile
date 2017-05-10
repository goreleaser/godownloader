

samples:
	go run main.go -repo spf13/hugo > samples/godownload-hugo.sh
	go run main.go -repo goreleaser/goreleaser > samples/godownload-goreleaser.sh
.PHONY: samples

clean:
	rm -rf ./bin
	git gc --aggressive

