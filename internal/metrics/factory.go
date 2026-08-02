package metrics

import (
	"fmt"
	"runtime"
)

func NewCollector() (Collector, error) {
	switch runtime.GOOS {
	case "linux":
		return NewLinuxCollector(), nil
	case "windows":
		return NewWindowsCollector(), nil
	default:
		return nil, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

}
