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

// TODO: REFACTOR
func (server *Server) Encode(path string) interface{} {
	metrics := server.store.SampleAll()
	if path == "." {
		m := make(map[string]interface{})
		EncodeAll(metrics, m)
		return m
	} else {
		return EncodeOne(path[1:len(path)], metrics)
	}
}

// TODO: REFACTOR
func EncodeAll(metrics map[string]ekg_core.Value, m map[string]interface{}) {
	for k, _ := range metrics {
		EncodeNested(k, metrics, m)
	}
}

// TODO: REFACTOR
func EncodeOne(path string, metrics map[string]ekg_core.Value) interface{} {
	if val, ok := metrics[path]; ok {
		return EncodeMetric(val)
	}
	return nil
}

// TODO: REFACTOR
func EncodeNested(path string, metrics map[string]ekg_core.Value, m map[string]interface{}) {
	paths := strings.Split(path, ".")
	s := slashes2dots(path)
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

// TODO: REFACTOR
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

func slashes2dots(s string) string {
	return strings.Map(slash2dot, s)
}

func slash2dot(ch rune) rune {
	if ch == '/' {
		return '.'
	} else {
		return ch
	}
}

func dots2slashes(s string) string {
	return strings.Map(dot2slash, s)
}

func dot2slash(ch rune) rune {
	if ch == '.' {
		return '/'
	} else {
		return ch
	}
}
