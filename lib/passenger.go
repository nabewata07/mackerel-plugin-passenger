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
	Prefix   string
	WorkDir  string
	IsBundle bool
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
			Label: (labelPrefix + " ProccessesInQueue"),
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

	r := regexp.MustCompile(`Memory\s+: (\d)+M`)

	stat := make(map[string]float64)
	stat["processes_in_queue"] = 0
	stat["total_processes"] = 0
	stat["total_memory"] = 0

	scanner := bufio.NewScanner(strings.NewReader(res))
	for scanner.Scan() {
		tmp := scanner.Text()
		if strings.Contains(tmp, "Requests in queue:") {
			arr := strings.Split(tmp, " ")
			pNum, err := strconv.ParseFloat(arr[5], 32)
			if err != nil {
				return stat, err
			}
			stat["processes_in_queue"] += pNum
		} else if strings.Contains(tmp, "* PID") {
			stat["total_processes"] += 1
		} else if r.MatchString(tmp) {
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

func getPassengerStatus(p PassengerPlugin) (string, error) {
	cmd := exec.Command("passenger-status", "--no-header")
	if p.IsBundle {
		cmd = exec.Command("bundle", "exec", "passenger-status", "--no-header")
		cmd.Dir = p.WorkDir
	}

	res, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func Do() {
	optTempfile := flag.String("tempfile", "", "Tempfile name")
	optWorkDir := flag.String("work-dir", "", "work directory")
	optIsBundle := flag.Bool("is-bundle", false, "is using bundler")
	flag.Parse()

	var p PassengerPlugin
	p.WorkDir = fmt.Sprintf("%s", *optWorkDir)
	p.IsBundle = *optIsBundle

	helper := mp.NewMackerelPlugin(p)
	helper.Tempfile = *optTempfile
	helper.Run()
}
