package adctlpath

import "path/filepath"

func CacheIndexFile(name string) string {
	if name != "" {
		name += "-"
	}
	return name + "index.yaml"
}

func CacheChartsFile(name string) string {
	if name != "" {
		name += "-"
	}
	return name + "charts.txt"
}

func ConfigPath(elem ...string) string {
	base := configHome()
	return filepath.Join(base, "adctl", filepath.Join(elem...))
}

func CachePath(elem ...string) string {
	base := cacheHome()
	return filepath.Join(base, "adctl", filepath.Join(elem...))
}
