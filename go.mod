module github.com/Naist4869/awesomeProject

go 1.14

require (
	github.com/Naist4869/log v0.0.0
	github.com/PuerkitoBio/goquery v1.5.1
	github.com/davecgh/go-spew v1.1.1
	github.com/gin-gonic/gin v1.6.3
	github.com/go-kratos/kratos v0.5.0
	github.com/golang/protobuf v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.6.0
	github.com/russross/blackfriday/v2 v2.0.1
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/smartystreets/goconvey v1.6.4
	github.com/stretchr/testify v1.4.0
	go.mongodb.org/mongo-driver v1.3.2
	go.uber.org/zap v1.15.0
	golang.org/x/sys v0.0.0-20200420163511-1957bb5e6d1f
	golang.org/x/tools v0.0.0-20191105231337-689d0f08e67a
	gopkg.in/go-playground/validator.v9 v9.31.0
	gopkg.in/yaml.v2 v2.2.8
)

replace github.com/Naist4869/log v0.0.0 => ./container/log
