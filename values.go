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

func (server *Server) Encode(path string) map[string]interface{} {
	metrics := server.store.SampleAll()
	m := make(map[string]interface{})
	if path == "/" {
		EncodeAll(metrics, m)
	} else {
		EncodeOne(path[1:len(path)], metrics, m)
	}
	return m
}

func EncodeAll(metrics map[string]ekg_core.Value, m map[string]interface{}) {
	for k, _ := range metrics {
		EncodeOne(k, metrics, m)
	}
}

func EncodeOne(path string, metrics map[string]ekg_core.Value, m map[string]interface{}) {
	paths := strings.Split(path, ".")
	s := strings.Map(slashes2dots, path)
	if val, ok := metrics[s]; ok {
		nestedMaps(paths, 0, EncodeMetric(val), m)
	}
}

func EncodeMetric(metric ekg_core.Value) value {
	switch metric.Typ {
	case ekg_core.COUNTER:
		return value{T: "c", V: metric.Val}
	case ekg_core.GAUGE:
		return value{T: "g", V: metric.Val}
	case ekg_core.LABEL:
		return value{T: "l", V: metric.Val}
	case ekg_core.DISTRIBUTION:
		return value{T: "d", V: metric.Val}
	case ekg_core.TIMESTAMP:
		return value{T: "t", V: metric.Val}
	case ekg_core.BOOL:
		return value{T: "b", V: metric.Val}
	default:
		return value{T: "u", V: false}
	}
}

// eek
func nestedMaps(paths []string, index int, v value, m map[string]interface{}) {

	if index >= len(paths) {
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
		if index+1 == len(paths) {
			nestedMaps(paths, index+1, v, m)
		} else {
			nestedMaps(paths, index+1, v, m2)
		}
	}
}

func slashes2dots(ch rune) rune {
	if ch == '/' {
		return '.'
	} else {
		return ch
	}
}
