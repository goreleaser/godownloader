#!/bin/sh -ex

# use gfind on osx
if command -v gfind >/dev/null 2>&1; then
	alias find=gfind
fi

# add ./bin to PATH as well
export PATH="./bin:$PATH"

# clean up
rm -rf ./www/public
rm -rf ./www/static/github.com
rm -rf ./www/data/projects
mkdir -p ./www/data/projects

# generate the sh files
./godownloader --tree=tree ./www/static/

# backup of disabled/broken/archived projects
mkdir -p www/static/github.com/alecthomas
wget -O www/static/github.com/alecthomas/gometalinter.sh \
	https://install.goreleaser.com/github.com/alecthomas/gometalinter.sh

mkdir -p www/static/github.com/kaihendry
wget -O www/static/github.com/kaihendry/lk2.sh \
	https://install.goreleaser.com/github.com/kaihendry/lk2.sh

# lint generated files
# SC2034 is unused variable
# some generated scripts contain 1 or more variables with aren't used
# sometimes.
find ./www/static -name '*.sh' | while read -r f; do
	shellcheck -e SC2034 -s sh "$f"
	shellcheck -e SC2034 -s bash "$f"
	shellcheck -e SC2034 -s dash "$f"
	shellcheck -e SC2034 -s ksh "$f"
done

# generate the hugo data files
find tree -name '*.yaml' -printf '%P\n' | while read -r f; do
	ff="$(echo "$f" | sed -e 's/\.yaml//' -e 's/\./-/g' -e 's/\//-/g')"
	echo "path: $f" | sed 's/\.yaml//' > ./www/data/projects/"$ff.yaml"
done

# generate the site
hugo -s www -d public
