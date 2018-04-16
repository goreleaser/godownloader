#!/bin/sh -ex

# clean up
rm -rf ./www/public
rm -rf www/static/github.com
rm ./www/data/projects/*.yaml

# generate the sh files
./godownloader --tree=tree www/static/

# generate the hugo data files
gfind tree -name '*.yaml' -printf '%P\n' | while read -r f; do
	ff="$(echo "$f" | sed -e 's/\.yaml//' -e 's/\./-/g' -e 's/\//-/g')"
	echo "path: $f" | sed 's/\.yaml//' > ./www/data/projects/"$ff.yaml"
done

# generate the site
hugo -s www -d public
