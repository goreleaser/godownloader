#!/bin/sh -ex

./godownloader -repo gohugoio/hugo > samples/godownloader-hugo.sh
./godownloader -repo goreleaser/goreleaser > samples/godownloader-goreleaser.sh
./godownloader -repo client9/misspell > samples/godownloader-misspell.sh
./godownloader -source equinoxio -repo tdewolff/minify > samples/godownloader-minify.sh
./godownloader -source raw -repo mvdan/sh -exe shfmt > samples/godownloader-shfmt.sh
chmod a+x samples/*.sh
