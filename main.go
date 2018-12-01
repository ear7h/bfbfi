package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := os.Stdin
	prompt := ""

	if len(os.Args) == 2 {
		if os.Args[1] == "-i" {
			prompt = ":bfbfi: "
		} else {
			fi, err := os.Open(os.Args[1])
			if err != nil {
				panic(err)
			}
			reader = fi
		}
	}

	vm := VM{}
	scn := bufio.NewScanner(reader)

	fmt.Printf(prompt)
	for scn.Scan() {
		fmt.Printf(prompt)
		vm.Write(scn.Bytes())
	}
}

type VM struct {
	input chan []byte
	lock  chan struct{}
}

func (vm *VM) Write(byt []byte) (int, error) {
	if vm.input == nil {
		vm.input = make(chan []byte)
		vm.lock = make(chan struct{})
		go vm.start()
	}

	vm.input <- byt
	<-vm.lock
	return len(byt), nil
}

func (vm *VM) start() {
	fmt.Println("vm started")

	pc := 0
	p := []byte{}
	ll := make([]int, 0, 64)
	bd := 0 //break depth
	dp := 0
	d := [30000]byte{}
	e := true

	for v := range vm.input {
		p = append(p, v...)
		for ; pc < len(p); pc++ {
			c := p[pc]

			//fmt.Printf("loop(%d): %c %v\n", pc, c, d[:4])
			if !e {
				if c == '[' {
					bd++
				} else if c == ']' {
					if bd == 0 {
						e = true
					} else {
						bd--
					}
				}
				// fmt.Println("break depth", bd)

				continue
			}
			//fmt.Printf("loop(%d): %c %v\n", pc, c, d[:4])

			switch c {
			case '>':
				dp++
			case '<':
				dp--
			case '+':
				d[dp]++
			case '-':
				d[dp]--
			case '.':
				os.Stdin.Write(d[dp : dp+1])
			case ',':
				os.Stdin.Read(d[dp : dp+1])
			case '[':
				e = d[dp] != 0
				if e {
					ll = append(ll, pc)
					//bd++
				}
			case ']':
				//fmt.Println(ll)
				pc = ll[len(ll)-1] - 1
				ll = ll[:len(ll)-1] // pop
			}
		}

		vm.lock <- struct{}{}
	}
}
