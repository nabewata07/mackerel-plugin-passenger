package mppassenger

import (
	"bufio"
	"flag"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin"
)

type PassengerPlugin struct {
	Prefix     string
	WorkDir    string
	BundlePath string
	StatusPath string
}

func (p PassengerPlugin) MetricKeyPrefix() string {
	if p.Prefix == "" {
		p.Prefix = "passenger"
	}
	return p.Prefix
}

func (p PassengerPlugin) GraphDefinition() map[string]mp.Graphs {
	labelPrefix := strings.Title(p.Prefix)
	return map[string]mp.Graphs{
		"processes": mp.Graphs{
			Label: (labelPrefix + " Processes"),
			Unit:  mp.UnitInteger,
			Metrics: []mp.Metrics{
				mp.Metrics{Name: "processes_in_queue", Label: "ProcessesInQueue"},
				mp.Metrics{Name: "total_processes", Label: "TotalProcesses"},
			},
		},
		"memory": mp.Graphs{
			Label: (labelPrefix + " MemoryInUse"),
			Unit:  mp.UnitBytes,
			Metrics: []mp.Metrics{
				mp.Metrics{
					Name:  "total_memory",
					Label: "TotalMemory",
					Scale: 1024 * 1024 * 1.0,
				},
			},
		},
	}
}

func (p PassengerPlugin) FetchMetrics() (map[string]float64, error) {
	res, err := getPassengerStatus(p)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch passenger-status: %s", err)
	}

	stat := make(map[string]float64)
	stat["processes_in_queue"] = 0
	stat["total_processes"] = 0
	stat["total_memory"] = 0

	r := regexp.MustCompile(`Memory\s+: (\d+)M`)

	scanner := bufio.NewScanner(strings.NewReader(res))
	for scanner.Scan() {
		tmp := scanner.Text()
		switch {
		case strings.Contains(tmp, "Requests in queue:"):
			arr := strings.Split(tmp, " ")
			pNum, err := strconv.ParseFloat(arr[5], 32)
			if err != nil {
				return stat, err
			}
			stat["processes_in_queue"] += pNum
		case strings.Contains(tmp, "* PID"):
			stat["total_processes"] += 1
		case r.MatchString(tmp):
			match := r.FindStringSubmatch(tmp)
			m, err := strconv.ParseFloat(match[1], 32)
			if err != nil {
				return stat, err
			}
			stat["total_memory"] += m
		}
	}
	return stat, nil
}

type PassengerStatusError struct {
	Stdout string
	Err    error
}

func (e *PassengerStatusError) Error() string {
	return fmt.Sprintf(
		"output of passenger-status : %s\nerror output : %s",
		e.Stdout, e.Err.Error(),
	)
}

func getPassengerStatus(p PassengerPlugin) (string, error) {
	cmdAry := [4]string{}
	statusIndex := 0
	cmd := exec.Command("passenger-status", "--no-header")

	if p.BundlePath != "" {
		statusIndex = 2
		cmdAry[0] = p.BundlePath
		cmdAry[1] = "exec"
		cmdAry[statusIndex] = "passenger-status"
	}

	if p.StatusPath != "" {
		cmdAry[statusIndex] = p.StatusPath
	}

	if len(cmdAry) > 0 {
		cmdAry[statusIndex+1] = "--no-header"
		cmd = exec.Command(cmdAry[0], cmdAry[1:]...)
	}

	if p.WorkDir != "" {
		cmd.Dir = p.WorkDir
	}

	res, err := cmd.Output()
	if err != nil {
		pErr := &PassengerStatusError{
			Stdout: string(res),
			Err:    err,
		}
		return "", pErr
	}

	return string(res), nil
}

func Do() {
	optTempfile := flag.String("tempfile", "", "Tempfile name")
	optWorkDir := flag.String("work-dir", "", "work directory")
	bundlePath := flag.String("bundle-path", "", "path of bundle command")
	statusPath := flag.String("status-path", "", "path of passenger-status command")
	flag.Parse()

	var p PassengerPlugin
	p.WorkDir = fmt.Sprintf("%s", *optWorkDir)
	p.BundlePath = fmt.Sprintf("%s", *bundlePath)
	p.StatusPath = fmt.Sprintf("%s", *statusPath)

	helper := mp.NewMackerelPlugin(p)
	helper.Tempfile = *optTempfile
	helper.Run()
}
