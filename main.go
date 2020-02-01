package main

//Import the packages we need
import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/sensu/sensu-go/types"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"github.com/spf13/cobra"
)

//Set up some variables. Most notably, warning and critical as load 1/5/15 values
var (
	warning, critical string
	stdin             *os.File
)

//Start our main function
func main() {
	rootCmd := configureRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

//Set up our flags for the command. Note that we have load defaults for warning & critical
func configureRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sensu-go-load-check",
		Short: "The Sensu Go check for system Load",
		RunE:  run,
	}

	cmd.Flags().StringVarP(&warning,
		"warning",
		"w",
		"2.75, 2.5, 2.0",
		"Load WARNING threshold, 1/5/15 min average")

	cmd.Flags().StringVarP(&critical,
		"critical",
		"c",
		"3.5, 3.25, 3.0",
		"Load CRITICAL threshold, 1/5/15 min average")

	return cmd
}

func parseArg(arg string) []float64 {
	var list []float64
	for _, x := range strings.Split(arg, ",") {
		w, err := strconv.ParseFloat(strings.TrimSpace(x), 64)
		if err != nil {
			msg := fmt.Sprintf("Failed to parse %s", err.Error())
			io.WriteString(os.Stdout, msg)
			os.Exit(3)
		}
		list = append(list, w)
	}

	return list
}

func run(cmd *cobra.Command, args []string) error {

	if len(args) != 0 {
		_ = cmd.Help()
		return fmt.Errorf("invalid argument(s) received")
	}

	if stdin == nil {
		stdin = os.Stdin
	}

	event := &types.Event{}

	return checkLoad(event)
}

//Here we start the meat of what we do.
func checkLoad(event *types.Event) error {

	const checkName = "CheckLoad"
	const metricName = "core_load"

	warn := parseArg(warning)
	crit := parseArg(critical)

	cpuLoad, err := load.Avg()
	if err != nil {
		msg := fmt.Sprintf("Failed to obtain load info %s", err.Error())
		io.WriteString(os.Stdout, msg)
		os.Exit(3)
	}

	cpuCount, err := cpu.Counts(false)
	if err != nil {
		msg := fmt.Sprintf("Failed to obtain CPU counts %s", err.Error())
		io.WriteString(os.Stdout, msg)
		os.Exit(3)
	}

	cpuLoadList := []float64{cpuLoad.Load1, cpuLoad.Load5, cpuLoad.Load15}

	// Calculate load per core
	for i, x := range cpuLoadList {
		cpuLoadList[i] = x / float64(cpuCount)
	}

	// Detect total level
	// 0=ok, 1=warn, 2=crit
	var level int = 0
	for i, v := range cpuLoadList {
		if v > crit[i] {
			if 2 > level {
				level = 2
			}
		} else if v >= warn[i] && v <= crit[i] {
			if 1 > level {
				level = 1
			}
		}
	}

	levelStrings := []string{"OK", "WARNING", "CRITICAL"}

	msg := fmt.Sprintf(
		"%[1]s %[2]s - value = %.2[4]f, %.2[5]f, %.2[6]f | %[3]s_1=%.2[4]f, %[3]s_5=%.2[5]f, %[3]s_15=%.2[6]f\n",
		checkName,
		levelStrings[level],
		metricName,
		cpuLoadList[0],
		cpuLoadList[1],
		cpuLoadList[2])
	io.WriteString(os.Stdout, msg)
	os.Exit(level)

	return nil
}
