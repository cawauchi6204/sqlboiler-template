module todoapp

go 1.20

require (
	github.com/labstack/echo/v4 v4.9.0
	github.com/go-sql-driver/mysql v1.7.0
	github.com/volatiletech/sqlboiler v3.7.1+incompatible
	github.com/volatiletech/null v8.0.0+incompatible
	github.com/golang-migrate/migrate/v4 v4.16.2
)

// 依存関係の明示的な指定
require (
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/labstack/gommon v0.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	golang.org/x/time v0.3.0 // indirect
)
