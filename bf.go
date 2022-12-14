package main

import "fmt"
import "io/ioutil"
import "time"

type Node struct {
	kind  bool
	left  []Node
	right byte
}

type TuringMachine struct {
	cursor       int
	instructions Node
	verbose      bool
	memory       Memory
}

type Memory []int64

func allocMemory(size int) Memory {
	m := []int64{}
	for i := 0; i < size; i += 1 {
		m = append(m, int64(0))
	}
	return m
}

func extendMemoryRight(m []int64) []int64 {
	for i := 0; i < 20; i += 1 {
		m = append(m, int64(0))
	}
	return m
}

func extendMemoryLeft(mem []int64) []int64 {
	obj := allocMemory(len(mem) + 20)
	for i := 20; i < len(mem)+20; i += 1 {
		obj[i] = mem[i-20]
	}
	return obj
}

func makeTuringMachine() *TuringMachine {
	return &TuringMachine{0, Node{true, []Node{}, ' '}, true, allocMemory(50)}
}

func parseBF(expr string, idx int) Node {
	if len(expr) == 0 {
		fmt.Println("Invalid input")
	}
	token := expr[idx]
	idx += 1
	if token == '[' {
		//tree := []Node{}
		nd := Node{
			kind:  false,
			left:  []Node{},
			right: ' ',
		}
		for expr[idx] != ']' && idx < len(expr) {
			nd.left = append(nd.left, parseBF(expr, idx))
			idx += 1
		}
		idx += 1
		return nd
	}
	//fmt.Println("Inserting into the ast: ", string(token))
	return Node{true, []Node{}, token}
}
func (nd Node) inspect() string {
	str := ""
	if nd.kind {
		str += string(nd.right)
		return str
	} else {
		sep := ", "
		for i, stuff := range nd.left {
			if i == len(nd.left)-1 {
				sep = ""
			}
			str += stuff.inspect() + sep
		}
	}
	return "[" + str + "]"
}

func (tm *TuringMachine) inspect() string {
	return "< cursor: " + fmt.Sprintf("%v", tm.cursor) + " | " + fmt.Sprintf("%v", tm.memory) + " >"
}

func (arr Memory) pp(position int) string {
	str := ""
	sep := ""
	for _, val := range arr[0:position] {
		str += fmt.Sprintf("%d", val) + ", "
	}
	if position == len(arr)-1 {
		sep = ", "
	}
	str += "<" + fmt.Sprintf("%d", arr[position]) + ">" + sep
	for idx, valeur := range arr[position+1:] {
		if idx == len(arr[position:])-1 {
			sep = ""
		}
		str += fmt.Sprintf("%d", valeur) + sep
	}
	return str
}

func (tm *TuringMachine) atomiq_execution(op string) {
	//fmt.Printf("Executing: %v\n", op)
	switch op {
	case "+":
		tm.memory[tm.cursor] += 1
	case "-":
		tm.memory[tm.cursor] -= 1
	case ">":
		if tm.cursor+1 >= len(tm.memory) {
			tm.memory = extendMemoryRight(tm.memory)
		}
		tm.cursor += 1
	case "<":
		if tm.cursor <= 0 {
			tm.memory = extendMemoryLeft(tm.memory)
		}
		tm.cursor -= 1
	case ".":
		fmt.Print(fmt.Sprintf("%s", string(tm.memory[tm.cursor])))
	case ",":
		fmt.Scanf("%d", &tm.memory[tm.cursor])
	}
}

func (tm *TuringMachine) execute() {
	nd := tm.instructions
	if nd.kind {
		op := string(nd.right)
		tm.atomiq_execution(fmt.Sprintf("%v", op))
		if tm.verbose {
			//fmt.Println(tm.memory.pp(tm.cursor))
			fmt.Printf("\033[2K\r%s", tm.memory.pp(tm.cursor))
			time.Sleep(40 * time.Millisecond)
		}
	} else {
		obj := nd.left
		/*
			for _, stuff := range obj {
				old := tm.instructions
				tm.instructions = stuff
				tm.execute()
				tm.instructions = old
			}
		*/
		for tm.memory[tm.cursor] != 0 {
			for _, stuff := range obj {
				old := tm.instructions
				tm.instructions = stuff
				tm.execute()
				tm.instructions = old
			}
		}
	}
}

func (tm *TuringMachine) execute_from_file(file_name string) error {
	content, err := ioutil.ReadFile(file_name)
	tm.run(string(content))
	return err
}

func (tm *TuringMachine) run(code string) *TuringMachine {
	nd := parseBF(code, 0)
	for _, stuff := range nd.left {
		tm.instructions = stuff
		tm.execute()
	}
	return tm
}

func (tm *TuringMachine) repl() {
	var fileName string
	for {
		var input string
		fmt.Printf("~ % ")
		fmt.Scanf("%s", &input)
		input = string(input)
		if input == ":i" {
			fmt.Print("Enter file name: ")
			fmt.Scanf("%s", &fileName)
			tm.execute_from_file(fileName)
		} else if input == ":q" {
			break
		} else if input == ":r" {
			tm.memory = allocMemory(50)
		} else if input == ":v" {
			tm.verbose = !tm.verbose
		} else if input == ":p" {
			fmt.Println(tm.memory.pp(tm.cursor))
		} else {
			tm.run("[" + input + "]")
		}
	}
}

func main() {
	//hw := "[++++>--->++<<.>.]"
	makeTuringMachine().repl() //.execute_from_file("./o.bf")
}
