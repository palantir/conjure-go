module github.com/palantir/conjure-go

go 1.13

require (
	github.com/danverbraganza/varcaser v0.0.0-20190207223536-e3fb03ee5b4c
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/julienschmidt/httprouter v1.3.0
	github.com/mholt/archiver v2.1.0+incompatible // indirect
	github.com/nmiyake/pkg/dirs v1.0.0
	github.com/nwaples/rardecode v1.0.0 // indirect
	github.com/palantir/conjure-go-runtime v0.4.0
	github.com/palantir/go-ptimports/v2 v2.9.0
	github.com/palantir/goastwriter v0.0.1
	github.com/palantir/godel-conjure-plugin v4.3.0+incompatible
	github.com/palantir/godel/v2 v2.22.0
	github.com/palantir/pkg/bearertoken v1.0.0
	github.com/palantir/pkg/binary v1.0.0
	github.com/palantir/pkg/cobracli v1.0.0
	github.com/palantir/pkg/datetime v1.0.0
	github.com/palantir/pkg/httpserver v1.0.0
	github.com/palantir/pkg/rid v1.0.0
	github.com/palantir/pkg/safejson v1.0.0
	github.com/palantir/pkg/safelong v1.0.0
	github.com/palantir/pkg/safeyaml v1.0.0
	github.com/palantir/pkg/uuid v1.0.0
	github.com/palantir/witchcraft-go-error v1.2.0
	github.com/palantir/witchcraft-go-logging v1.5.0
	github.com/palantir/witchcraft-go-server v1.13.0
	github.com/pkg/errors v0.8.1
	github.com/spf13/cobra v0.0.5
	github.com/stretchr/testify v1.4.0
	golang.org/x/tools v0.0.0-20191122232904-2a6ccf25d769
	gopkg.in/yaml.v2 v2.2.7
)

replace github.com/palantir/pkg => github.com/palantir/pkg v1.0.0

replace github.com/nmiyake/pkg => github.com/nmiyake/pkg v0.0.0
