module weqiniu

go 1.12

require (
	github.com/eclipse/paho.mqtt.golang v1.2.0
	github.com/ikuiki/wwdk v2.2.0+incompatible
	github.com/kataras/golog v0.0.0-20180321173939-03be10146386
	github.com/pkg/errors v0.8.1
	github.com/qiniu/api.v7 v7.2.5+incompatible
	github.com/qiniu/x v7.0.8+incompatible // indirect
	github.com/stretchr/testify v1.3.0 // indirect
	golang.org/x/net v0.0.0-20190503192946-f4e77d36d62c // indirect
	qiniupkg.com/x v7.0.8+incompatible // indirect
	wegate v0.0.0-00010101000000-000000000000
)

replace wegate => github.com/ikuiki/wegate v0.0.0-20190507071455-d5b73ce5cb06

replace github.com/liangdas/mqant => github.com/ikuiki/mqant v1.8.1-0.20190427142930-7dabfa32d064
