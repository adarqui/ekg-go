package main

import (
	"fmt"
	"github.com/adarqui/ekg-go"
)

func main() {
	serv, _ := ekg.ForkServer("127.0.0.1:8111")
	_ = serv.GetCounter("app.loop.counter")
	_ = serv.GetStore().CreateCounter("app.loop.counter_direct")
	_ = serv.GetGauge("app.loop.gauge")
	_ = serv.GetLabel("app.loop.label")
	_ = serv.GetDistribution("app.loop.distribution")

	metrics := serv.GetStore().SampleAll()
	v := ekg.EncodeAll(metrics)

	fmt.Println(v)
}
