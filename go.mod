module wstrade

go 1.22.0

require (
	github.com/dictxwang/go-binance v0.1.13
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	go.uber.org/zap v1.27.0
	golang.org/x/time v0.5.0
)

require (
	github.com/bitly/go-simplejson v0.5.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/jonboulle/clockwork v0.4.0 // indirect
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rogpeppe/go-internal v1.9.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
)

replace github.com/dictxwang/go-binance v0.1.13 => ./go-binance
