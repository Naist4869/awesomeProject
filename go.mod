module github.com/Naist4869/awesomeProject

go 1.14

require (
	github.com/Naist4869/log v0.0.0
	github.com/PuerkitoBio/goquery v1.5.1
	github.com/davecgh/go-spew v1.1.1
	github.com/pkg/errors v0.9.1
	github.com/russross/blackfriday/v2 v2.0.1
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/smartystreets/goconvey v1.6.4
	github.com/stretchr/testify v1.4.0
	go.mongodb.org/mongo-driver v1.3.2
	go.uber.org/zap v1.15.0
	golang.org/x/sys v0.0.0-20200116001909-b77594299b42
	gopkg.in/yaml.v2 v2.2.8
)

replace github.com/Naist4869/log v0.0.0 => ./container/log
