module github.com/Naist4869/awesomeProject

go 1.14

require (
	github.com/Naist4869/log v0.0.0
	github.com/PuerkitoBio/goquery v1.5.1
	github.com/davecgh/go-spew v1.1.1
	github.com/go-kratos/kratos v0.5.0
	github.com/go-playground/validator/v10 v10.3.0
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.6.0
	github.com/russross/blackfriday/v2 v2.0.1
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/smartystreets/goconvey v1.6.4
	github.com/stretchr/testify v1.4.0
	go.mongodb.org/mongo-driver v1.3.2
	go.uber.org/zap v1.15.0
	google.golang.org/genproto v0.0.0-20200402124713-8ff61da6d932
	google.golang.org/grpc v1.28.1
	gopkg.in/go-playground/validator.v9 v9.31.0
	gopkg.in/yaml.v2 v2.2.8
)

replace github.com/Naist4869/log v0.0.0 => ./container/log
