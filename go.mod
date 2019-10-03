module github.com/goreleaser/godownloader

go 1.12

require (
	github.com/apex/log v1.1.1
	github.com/client9/codegen v0.0.0-20180316044450-92480ce66a06
	github.com/goreleaser/goreleaser v0.118.2
	github.com/pkg/errors v0.8.1
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/yaml.v2 v2.2.4
)

// Fix invalid pseudo-version: revision is longer than canonical (6fd6a9bfe14e)
replace github.com/go-macaron/cors => github.com/go-macaron/cors v0.0.0-20190418220122-6fd6a9bfe14e
