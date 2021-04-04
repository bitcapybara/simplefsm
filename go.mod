module github.com/bitcapybara/simplefsm

go 1.15

require (
	github.com/bitcapybara/raft v0.0.1
	github.com/go-resty/resty/v2 v2.5.0
	github.com/kr/pretty v0.1.0 // indirect
	github.com/labstack/echo/v4 v4.2.1
	github.com/sirupsen/logrus v1.8.1
	github.com/vmihailenco/msgpack v4.0.4+incompatible
	golang.org/x/sys v0.0.0-20210330210617-4fbd30eecc44 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
)

replace github.com/bitcapybara/raft => ../raft
