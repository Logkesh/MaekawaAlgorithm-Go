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
	quorum []int
	inreq chan int
	inrel chan int
	response chan int
	queue []int
	data  *int
}

var processes []*Process
var coterie [][]int

// NewProcess creates a new process.
func NewProcess(id int, data *int) *Process {
	return &Process{
		id:    id,
		state: "alive",
		quorum: []int{},
		inreq: make(chan int,len(processes)),
		inrel: make(chan int,len(processes)) ,
		response: make(chan int,len(processes)),
		queue: []int{},
		data:  data,
	}
}

func (process *Process) ManageRequest() {
	for {
		select {
			case pid := <- process.inreq:
				if len(process.queue) == 0 {
					process.queue = append(process.queue, pid)
					processes[pid].response <- process.id
				} else {
					process.queue = append(process.queue, pid)
				}
			case <- process.inrel:
				process.queue = process.queue[1:]
				if len(process.queue) > 0 { processes[process.queue[0]].response <- process.id}
		}
	}
}

// RequestCS requests permission to enter the critical section.
func (process *Process) RequestCS() {

	fmt.Println(string(colorCyan), "Process", process.id, "requests CS", string(colorReset))
	
	for _, pid := range process.quorum {
		fmt.Println("Process", process.id, "sends request message to", pid)
		processes[pid].inreq <- process.id
	}

	// Wait for reply messages from a majority of processes in the quorum.
	fmt.Println("Process", process.id, "waits for reply message from all processes")

	granted := false
	for i := 0; i < len(process.quorum); i++ {
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		pid := <- process.response
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
	
	time.Sleep(time.Duration(1500) * time.Millisecond)

}

// ReleaseCS releases the critical section.
func (process *Process) ReleaseCS() {
	fmt.Println(string(colorRed), "Process", process.id, "releases CS", string(colorReset))
	// Send release messages to all processes in the quorum.
	for _, pid := range process.quorum {
		fmt.Println("Process", process.id, "sends release message to", pid)
		processes[pid].inrel <- process.id
	}
	// Change state to alive.
	process.state = "alive"
}

func (process *Process) StartProcess(wg *sync.WaitGroup) {
	process.state = "ready"
	fmt.Println(string(colorYellow), "Process", process.id, "wants to get into CS", string(colorReset))

	// requests permission to enter the critical section.
	process.RequestCS()

	if process.state == "critical" {
		process.ExecuteCS()
	} else {
		fmt.Println(string(colorRed), "Process", process.id, " cannot get into CS", string(colorReset))
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
	var shared_value int = -1

	// Create a set of processes.
	var numProcs int
	fmt.Println(string(colorCyan), "Enter the number of processes: ", string(colorReset))
	// fmt.Scanln(&numProcs)
	numProcs = 5

	processes = make([]*Process, numProcs)
	coterie = [][]int{
		{2,3},
		{2,4},
		{1,4},
		{1,2,4},
		{1,2,3},
	}
	for i := 0; i < numProcs; i++ {
		processes[i] = NewProcess(i, &shared_value)
		go processes[i].ManageRequest()
	}
	
	// Set the quorums for each process.
	for i := 0; i < numProcs; i++ {
		processes[i].quorum = coterie[i]
	}

	var numberOfCSaccess int = 5
	fmt.Println(string(colorCyan), "Enter the number of iterations: ", string(colorReset))
	// fmt.Scanln(&numberOfCSaccess)

	// Print the quorums.
	fmt.Printf("\n")
	fmt.Println(string(colorPurple), "-------------------", string(colorReset))
	fmt.Println(string(colorPurple), "Quorums", string(colorReset))
	fmt.Println(string(colorPurple), "-------------------", string(colorReset))
	for _, process := range processes {
		fmt.Printf("Process %d: ", process.id)
		for _, pid := range process.quorum {
			fmt.Printf("%d ", pid)
		}
		fmt.Printf("\n")
	}
	fmt.Println(string(colorPurple), "-------------------", string(colorReset))

	fmt.Printf("\n")

	var wg sync.WaitGroup
	wg.Add(numberOfCSaccess)

	// Create Multiple simulations
	for i := 0; i < numberOfCSaccess; i++ {
		var aliveProcesses = []*Process{}

		for j := 0; j < len(processes); j++ {
			if processes[j].state == "alive" {
				aliveProcesses = append(aliveProcesses, processes[j])
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

