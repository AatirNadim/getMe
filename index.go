package main

import (
	"fmt"
	// "io"
	"strings"
	// "go/types"
)


type Typo struct {
	name string
	age  int
}


type Abster interface {
	doSomeStuff() uint64
}



func stuff(typo Typo) {
	typo.age += 1
	typo.name = "Bob"
}


// func main() {
// 	i, j := 1, 2


	

// 	// abster.doSomeStuff()

// 	p := &i

// 	*p = *p + 1

// 	fmt.Println(i, j)


// 	typo := Typo{name: "Alice", age: 30}


// 	stuff(typo)

// 	fmt.Println("changed", typo)


// 	ptr := &typo


// 	fmt.Println(*ptr)


// 	var arr []int = []int{1, 2, 3, 4, 5}


// 	segment := arr[1:4]

// 	segment[0] = 10


// 	fmt.Println(arr)



// 	mp := &map[string]int{}

// 	(*mp)["one"] = 1
// 	(*mp)["two"] = 2
// 	(*mp)["three"] = 3

// 	fmt.Println(*mp)

// 	input := strings.NewReader("Hello, World!")

// 	fmt.Println("input received\n", input)

// 	// r := strings.NewReader("Hello, Reader!")

// 	// b := make([]byte, 8)
// 	// for {
// 	// 	n, err := r.Read(b)
// 	// 	fmt.Printf("n = %v err = %v b = %v\n", n, err, b)
// 	// 	fmt.Printf("b[:n] = %q\n", b[:n])
// 	// 	if err == io.EOF {
// 	// 		break
// 	// 	}
// 	// }





// 	fmt.Println("Reading from input:")
// 	_, err := fmt.Scanln(&input)

// 	if (err != nil) {
// 		fmt.Errorf("error occurred: %w", err)
// 	} else {
// 		fmt.Println("no error, data received:", input)
// 	}
	
	

// }