package main

import (
	"fmt"
	"github.com/adarqui/ekg-go"
	"strconv"
	"time"
)

func main() {
	serv, _ := ekg.ForkServer("127.0.0.1:8111")
	counter := serv.GetCounter("app.loop.counter")
	counter_direct := serv.GetStore().CreateCounter("app.loop.counter_direct")
	gauge := serv.GetGauge("app.loop.gauge")
	label := serv.GetLabel("app.loop.label")
	distr := serv.GetDistribution("app.loop.distribution")

	go func() {
		// app loop
		for {
			time.Sleep(1 * time.Second)
			i := counter.Read()
			counter.Inc()
			counter_direct.Inc()
			gauge.Set(i)
			label.Modify(func(s string) string { return fmt.Sprintf("%s%s", s, strconv.Itoa(int(i))) })
			distr.Add(i)
		}
	}()

	select {}
}
