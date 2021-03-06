package nagiosfoundation

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/ncr-devops-platform/nagiosfoundation/lib/pkg/memory"
	"github.com/ncr-devops-platform/nagiosfoundation/lib/pkg/nagiosformatters"
)

// CheckAvailableMemoryWithHandler determines the percentage of free memory
// remaining then emits a critical response if it's below
// flag -critical, a warning response if below flag -warning,
// and good response otherwise.
func CheckAvailableMemoryWithHandler(memoryHandler func() uint64) (string, int) {
	const checkName = "CheckAvailableMemory"
	var warning = flag.Float64("warning", 85, "the memory threshold to issue a warning alert")
	var critical = flag.Float64("critical", 95, "the memory threshold to issue a critical alert")
	var metricName = flag.String("metric_name", "memory_used_percentage", "the name of the metric generated by this check")
	flag.Parse()
	SetDefaultGlogStderr()

	var msg string
	var retcode int
	var usedMemoryPercentage uint64
	var err error

	if memoryHandler == nil {
		err = errors.New("No used memory percentage service")
	} else {
		usedMemoryPercentage = memoryHandler()
	}

	if err != nil || usedMemoryPercentage == 0 {
		msg = fmt.Sprintf("%s CRITICAL - %s", checkName, err)
		retcode = 2
	} else {
		msg, retcode = nagiosformatters.GreaterFormatNagiosCheck(checkName, float64(usedMemoryPercentage), *warning, *critical, *metricName)
	}

	return msg, retcode
}

// CheckAvailableMemory executes CheckAvailableMemoryWithHandler(),
// passing it the OS constranted GetFreeMemory() function, prints
// the returned message and exits with the returned exit code.
//
// Returns are those of CheckAvailableMemoryWithHandler()
func CheckAvailableMemory() {
	msg, retval := CheckAvailableMemoryWithHandler(memory.GetUsedMemoryPercentage)

	fmt.Println(msg)
	os.Exit(retval)
}
