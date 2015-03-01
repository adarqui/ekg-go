// * JSON serialization
package ekg

import (
	"github.com/adarqui/ekg-core-go"
    "strings"
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

/*
func EncodeAll(metrics map[string]ekg_core.Value) map[string]value {
	m := make(map[string]value)
	for k, v := range metrics {
		m[k] = EncodeMetric(v)
	}
	return m
}
*/
func EncodeAll(metrics map[string]ekg_core.Value) map[string]interface{} {
	m := make(map[string]interface{})
	for k, _ := range metrics {
        EncodeOne(k, metrics, m)
	}
	return m
}

func EncodeOne(path string, metrics map[string]ekg_core.Value, m map[string]interface{}) {
    paths := strings.Split(path, ".")
    s := strings.Map(func(ch rune) rune { if ch == '/' { return '.' } else { return ch } }, path)
    if val, ok := metrics[s]; ok {
        nestedMaps(paths, 0, EncodeMetric(val), m)
    }
}

func EncodeMetric(metric ekg_core.Value) value {
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

// eek
func nestedMaps(paths []string, index int, v value, m map[string]interface{}) {

    if (index >= len(paths)) {
        s := paths[index-1]
        m[s] = v
        return
    }

    s := paths[index]
    if val, ok := m[s]; ok {
        val2 := val.(map[string]interface{})
        nestedMaps(paths, index+1, v, val2)
    } else {
        m2 := make(map[string]interface{})
        m[s] = m2
        if (index+1 == len(paths)) {
            nestedMaps(paths, index+1, v, m)
        } else {
            nestedMaps(paths, index+1, v, m2)
        }
    }
}
