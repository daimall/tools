package cmdctrl

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/prometheus/procfs"
)

// 根据进程名，杀死进程
func KillProcessByName(processName string) bool {
	procs, err := procfs.AllProcs()
	if err != nil {
		return false
	}

	killed := false
	for _, p := range procs {
		cmdline, _ := p.CmdLine()
		var name string
		if len(cmdline) >= 1 {
			name = filepath.Base(cmdline[0])
		} else {
			name, _ = p.Comm()
		}

		if name == processName {
			process, err := os.FindProcess(p.PID)
			if err == nil {
				process.Kill()
				killed = true
			}
		}
	}
	return killed
}

// 根据进程名模糊匹配，杀死进程
func KillProcessByMatchName(processNameMatch string) bool {
	procs, err := procfs.AllProcs()
	if err != nil {
		return false
	}

	killed := false
	for _, p := range procs {
		cmdline, _ := p.CmdLine()
		var name string
		if len(cmdline) >= 1 {
			name = filepath.Base(cmdline[0])
		} else {
			name, _ = p.Comm()
		}
		if ok, err := regexp.MatchString(processNameMatch, name); err == nil && ok {
			process, err := os.FindProcess(p.PID)
			if err == nil {
				process.Kill()
				killed = true
			}
		}
	}
	return killed
}
