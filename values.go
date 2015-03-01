// * JSON serialization
package ekg

import (
	"github.com/adarqui/ekg-core-go"
)

type value struct {
	T string      `json:"type"`
	V interface{} `json:"val"`
}

/*
data MetricType =
      CounterType
    | GaugeType
    | LabelType
    | DistributionType
*/

/*
func metricType (v interface{}) string {
    return "a"
}
*/

func EncodeAll(metrics map[string]ekg_core.Value) map[string]value {
	m := make(map[string]value)
	for k, v := range metrics {
		m[k] = EncodeOne(v)
	}
	return m
}

func EncodeOne(metric ekg_core.Value) value {
	switch metric.Typ {
	case ekg_core.COUNTER:
		val := value{T: "c", V: metric.Val}
		return val
	case ekg_core.GAUGE:
		val := value{T: "g", V: metric.Val}
		return val
	case ekg_core.LABEL:
		val := value{T: "l", V: metric.Val}
		return val
	case ekg_core.DISTRIBUTION:
		val := value{T: "d", V: metric.Val}
		return val
	case ekg_core.TIMESTAMP:
		val := value{T: "t", V: metric.Val}
		return val
	default:
		val := value{T: "u", V: false}
		return val
	}
}
