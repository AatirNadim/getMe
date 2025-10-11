package main

import (
	"fmt"
	// "os/exec"
	"runtime"
	// "time"
)



func main() {

	// currtime := time.Now()

	// cmd := exec.Command("./testing-major.sh")
	// err := cmd.Run()


	// if err != nil {
	// 	fmt.Println("Error executing command:", err)
	// 	return
	// }

	// fmt.Println("time taken:", time.Since(currtime))


	fmt.Println("logical cores: ", runtime.NumCPU())

}


