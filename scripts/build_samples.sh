#!/bin/sh -ex
./godownloader --repo gohugoio/hugo >samples/godownloader-hugo.sh
./godownloader --repo goreleaser/goreleaser >samples/godownloader-goreleaser.sh
./godownloader --repo client9/misspell >samples/godownloader-misspell.sh
./godownloader --repo tdewolff/minify >samples/godownloader-minify.sh

# good example of multi-build and wrapped directory
./godownloader --repo https://github.com/alecthomas/gometalinter >samples/godownloader-gometalinter.sh

# binary and repo name do not match
./godownloader --source raw --repo mvdan/sh --exe shfmt >samples/godownloader-shfmt.sh

# uses zip
./godownloader --repo serverless/event-gateway >samples/godownloader-event-gateway.sh

chmod a+x samples/*.sh

# https://github.com/goreleaser/godownloader/issues/49
# still available and good to test equinoxio but no longer current
#./godownloader -source equinoxio -repo tdewolff/minify >samples/godownloader-minify.sh
