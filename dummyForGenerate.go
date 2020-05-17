package main

//go:generate go run --tags sdkcodegen ./internal/sdkcodegen ./docs/apis.md ./usecase/workwx/apis.md.go
//go:generate go run --tags sdkcodegen ./internal/sdkcodegen ./docs/chat_info.md ./model/wxmodel/chat_info.md.go
//go:generate go run --tags sdkcodegen ./internal/sdkcodegen ./docs/dept_info.md ./model/wxmodel/dept_info.md.go
//go:generate go run --tags sdkcodegen ./internal/sdkcodegen ./docs/user_info.md ./model/wxmodel/user_info.md.go
//go:generate go run --tags sdkcodegen ./internal/sdkcodegen ./docs/rx_msg.md ./model/wxmodel/rx_msg.md.go
//go:generate go run --tags sdkcodegen ./internal/errcodegen ./errcodes/mod.go
//go:generate go run --tags sdkcodegen ./internal/sdkcodegen ./docs/rx_msg_official.md ./model/officialmodel/rx_msg_official.md.go
