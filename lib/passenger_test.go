package mppassenger

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// some code here to check arguments perhaps?
	fmt.Fprintf(os.Stdout, strings.Join(os.Args[3:], " "))
	os.Exit(0)
}

func TestSimpleCommand(t *testing.T) {
	p := PassengerPlugin{}
	execCommand = fakeExecCommand
	defer func() { execCommand = exec.Command }()
	expected := "passenger-status --no-header"
	out, err := getPassengerStatus(p)
	if err != nil {
		t.Errorf("Expected nil error, got %#v", err)
	}
	if string(out) != expected {
		t.Errorf("Expected %q, got %q", expected, out)
	}
}

func TestCommandWithOptions(t *testing.T) {
	p := PassengerPlugin{
		BundlePath: "/path/to/bundle",
		StatusPath: "/path/to/passenger-status",
	}
	execCommand = fakeExecCommand
	defer func() { execCommand = exec.Command }()
	expected := "/path/to/bundle exec /path/to/passenger-status --no-header"
	out, err := getPassengerStatus(p)
	if err != nil {
		t.Errorf("Expected nil error, got %#v", err)
	}
	if string(out) != expected {
		t.Errorf("Expected %q, got %q", expected, out)
	}
}

func fakeExecStatus(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess2", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=2"}
	return cmd
}

func TestHelperProcess2(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "2" {
		return
	}

	out := `----------- General information -----------
Max pool size : 9
App groups    : 2
Processes     : 3
Requests in top-level queue : 0

----------- Application groups -----------
/var/www/passenger_status_service/current/public:
  App root: /var/www/passenger_status_service/current
  Requests in queue: 1
  * PID: 18257   Sessions: 0       Processed: 66179   Uptime: 5h 53m 29s
    CPU: 3%%      Memory  : 110M    Last used: 0s ago

/var/www/phusion_blog/current/public:
  App root: /var/www/phusion_blog/current
  Requests in queue: 2
  * PID: 18334   Sessions: 0       Processed: 4595    Uptime: 5h 53m 29s
    CPU: 0%%      Memory  : 99M     Last used: 4s ago
  * PID: 18339   Sessions: 0       Processed: 2873    Uptime: 5h 53m 26s
    CPU: 0%%      Memory  : 96M     Last used: 29s ago`
	// some code here to check arguments perhaps?
	fmt.Fprintf(os.Stdout, out)
	os.Exit(0)
}

func TestFetchMetrics(t *testing.T) {
	p := PassengerPlugin{}
	execCommand = fakeExecStatus
	defer func() { execCommand = exec.Command }()
	expected := map[string]float64 {
		"processes_in_queue": 3,
		"total_processes": 3,
		"total_memory": 305,
	}
	out, err := p.FetchMetrics()
	if err != nil {
		t.Errorf("Expected nil error, got %#v", err)
	}

	if out["processes_in_queue"] != expected["processes_in_queue"] {
		t.Errorf(
			"Expected %q, got %q",
			expected["processes_in_queue"],
			out["processes_in_queue"],
		)
	}
	if out["total_processes"] != expected["total_processes"] {
		t.Errorf(
			"Expected %q, got %q",
			expected["total_processes"],
			out["total_processes"],
		)
	}
	if out["total_memory"] != expected["total_memory"] {
		t.Errorf(
			"Expected %q, got %q",
			expected["total_memory"],
			out["total_memory"],
		)
	}
}

// func TestGenerateCmdAryNoOption(t *testing.T) {
	// p := PassengerPlugin{}
	// c := generateCmdBag(p)
	// e := []string{
		// "passenger-status",
		// "--no-header",
	// }
	// if c != e {
		// t.Errorf("\nexpected : %v\nactual : %v", e, c)
	// }
// }

// func TestGenerateCmdAry(t *testing.T) {
// }
