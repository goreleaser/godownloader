

samples:
	go run main.go -repo spf13/hugo > samples/godownloader-hugo.sh
	go run main.go -repo goreleaser/goreleaser > samples/godownloader-goreleaser.sh
.PHONY: samples

clean:
	rm -rf ./bin
	git gc --aggressive

