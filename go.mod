module github.com/Naist4869/awesomeProject

go 1.14

require (
	github.com/Naist4869/log v0.0.0
	github.com/pkg/errors v0.9.1
	go.mongodb.org/mongo-driver v1.3.2
	go.uber.org/zap v1.15.0
	gopkg.in/yaml.v2 v2.2.8
)

replace github.com/Naist4869/log v0.0.0 => ./container/log
