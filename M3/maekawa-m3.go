package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var colorReset = "\033[0m"

var colorRed = "\033[31m"
var colorGreen = "\033[32m"
var colorYellow = "\033[33m"
var colorBlue = "\033[34m"
var colorPurple = "\033[35m"
var colorCyan = "\033[36m"
var colorWhite = "\033[37m"

// Process represents a process in the system.
type Process struct {
	id    int
	state string
	queue []int
	data  *int
}

var processes []*Process

// NewProcess creates a new process.
func NewProcess(id int, data *int) *Process {
	return &Process{
		id:    id,
		state: "alive",
		queue: make([]int, 0),
		data:  data,
	}
}

// RequestCS requests permission to enter the critical section.
func (process *Process) RequestCS() {

	process.state = "ready"

	fmt.Println(string(colorCyan), "Process", process.id, "requests CS", string(colorReset))

	// Send request messages to all processes in the quorum.
	for _, pid := range process.queue {
		fmt.Println("Process", process.id, "sends request message to", pid)
	}

	// Wait for reply messages from a majority of processes in the quorum.
	granted := false
	for i := 0; i < len(process.queue); i++ {
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		pid := process.queue[i]
		fmt.Println("Process", process.id, "waits for reply message from", pid)

		// If a majority of processes have granted permission, enter the critical section.
		if pid != process.id && process.state != "critical" {
			fmt.Println("Process", process.id, "receives reply message from", pid)
			granted = granted || true
		}
	}

	// If permission was not granted, block.
	if !granted {
		fmt.Println(string(colorRed), "Process", process.id, "is blocked", string(colorReset))
		process.state = "blocked"
	} else {
		fmt.Println(string(colorRed), "Process", process.id, "enters CS", string(colorReset))
		process.state = "critical"
	}
}

func (process *Process) ExecuteCS() {
	fmt.Print(string(colorCyan), "Process ", process.id, " updates the value of shared_value from ", string(colorGreen), *process.data, string(colorCyan), " to", string(colorReset))

	// Update the shared value into the process ID
	*process.data = process.id * process.id

	fmt.Println(string(colorGreen), *process.data, string(colorReset))

}

// ReleaseCS releases the critical section.
func (process *Process) ReleaseCS() {
	fmt.Println(string(colorRed), "Process", process.id, "releases CS", string(colorReset))

	// Send release messages to all processes in the quorum.
	for _, pid := range process.queue {
		fmt.Println("Process", process.id, "sends release message to", pid)
	}

	// Change state to ready.
	process.state = "alive"
}

func (process *Process) StartProcess(wg *sync.WaitGroup) {
	fmt.Println(string(colorYellow), "Process", process.id, "wants to get into CS", string(colorReset))

	// requests permission to enter the critical section.
	process.RequestCS()

	if process.state == "critical" {
		process.ExecuteCS()
	}

	process.ReleaseCS()

	fmt.Printf("\n")
	wg.Done()
}

func main() {

	fmt.Printf("\n")
	fmt.Println(string(colorRed), "-------------------", string(colorReset))
	fmt.Println(string(colorRed), "Maekawa's Algorithm", string(colorReset))
	fmt.Println(string(colorRed), "-------------------", string(colorReset))
	fmt.Printf("\n")

	// Create a shared Variable.
	var shared_value int

	// Create a set of processes.
	var numProcs int
	fmt.Println(string(colorCyan), "Enter the number of processes: ", string(colorReset))
	fmt.Scanln(&numProcs)

	processes = make([]*Process, numProcs)
	for i := 0; i < numProcs; i++ {
		processes[i] = NewProcess(i, &shared_value)
	}

	// Set the quorums for each process.
	for i := 0; i < numProcs; i++ {
		processes[i].queue = []int{(i + 1) % numProcs, (i + 2) % numProcs}
	}

	var numberOfIterations int
	fmt.Println(string(colorCyan), "Enter the number of iterations: ", string(colorReset))
	fmt.Scanln(&numberOfIterations)

	// Print the quorums.
	fmt.Printf("\n")
	fmt.Println(string(colorPurple), "-------------------", string(colorReset))
	fmt.Println(string(colorPurple), "Quorums", string(colorReset))
	fmt.Println(string(colorPurple), "-------------------", string(colorReset))
	for _, process := range processes {
		fmt.Printf("Process %d: ", process.id)
		for _, pid := range process.queue {
			fmt.Printf("%d ", pid)
		}
		fmt.Printf("\n")
	}
	fmt.Println(string(colorPurple), "-------------------", string(colorReset))

	fmt.Printf("\n")

	var wg sync.WaitGroup
	wg.Add(numberOfIterations)
	// Create Multiple simulations
	for i := 0; i < numberOfIterations; i++ {
		var aliveProcesses []Process

		for j := 0; j < len(processes); j++ {
			if processes[j].state == "alive" {
				aliveProcesses = append(aliveProcesses, *processes[j])
			}
		}

		if len(aliveProcesses) > 0 {
			randomIndex := rand.Intn(len(aliveProcesses))
			randomProcess := aliveProcesses[randomIndex]
			go randomProcess.StartProcess(&wg)
		}

		time.Sleep(time.Duration(1000) * time.Millisecond)
	}
	wg.Wait()
}
