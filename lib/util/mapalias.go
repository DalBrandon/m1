// --------------------------------------------------------------------
// mapalias.go -- Provide alias to map keys for map[string]string
//
// Created 2018-09-26 DLB
// --------------------------------------------------------------------

package util

func MapAlias(params map[string]string, names ...string) (string, bool) {
	for _, n := range names {
		v, ok := params[n]
		if ok {
			return v, true
		}
	}
	return "", false
}
