// pseudo_test.go - test cases.
// The principle test is C_source/pseudo.c#main()

package pseudo

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
)

/*
   cat dimacsMaxf.txt
   p max 6 8
   n 1 s
   n 6 t
   a 1 2 5
   a 1 3 15
   a 2 4 5
   a 2 5 5
   a 3 4 5
   a 3 5 5
   a 4 6 15
   a 5 6 5
*/
func TestReadDimacsFile(t *testing.T) {
	s := NewSession(Context{})

	fh, err := os.Open("_data/dimacsMaxf.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()

	err = s.readDimacsFile(fh)
	if err != nil {
		t.Fatal(err)
	}

	// check some allocations made by 'p' record
	if s.numNodes != uint(6) {
		fmt.Println("numNodes != 6 :", s.numNodes)
		t.Fatal()
	}
	if s.numArcs != uint(8) {
		fmt.Println("numArcs != 8 :", s.numArcs)
		t.Fatal()
	}
	if uint(len(s.adjacencyList)) != s.numNodes {
		fmt.Println("len(adjacencyList):", len(s.adjacencyList), "numNodes:", s.numNodes)
		t.Fatal()
	}
	if uint(len(s.strongRoots)) != s.numNodes {
		fmt.Println("len(strongRoots):", len(s.strongRoots), "numNodes:", s.numNodes)
		t.Fatal()
	}
	if uint(len(s.labelCount)) != s.numNodes {
		fmt.Println("len(labelCount):", len(s.labelCount), "numNodes:", s.numNodes)
		t.Fatal()
	}
	if uint(len(s.arcList)) != s.numArcs {
		fmt.Println("len(arcList):", len(s.arcList), "numNodes:", s.numNodes)
		t.Fatal()
	}

	// check values set by 'n' records
	if s.source != uint(1) {
		fmt.Println("source != 1 :", s.source)
		t.Fatal()
	}
	if s.sink != uint(6) {
		fmt.Println("sink != 6 :", s.sink)
		t.Fatal()
	}

	// check arc record parsing
	checkVals := map[string]int{ "1_2":5, "1_3":15, "2_4":5, "2_5":5, "3_4":5, "3_5":5, "4_6":15, "5_6":5}
	for k, v := range s.arcList {
		ck := strconv.Itoa(int(v.from.number))+"_"+strconv.Itoa(int(v.to.number))
		if vcap, ok := checkVals[ck]; !ok {
			fmt.Println("unknown ck:", ck)
			t.Fatal()
		} else if vcap != v.capacity {
			fmt.Println(k, "- want:", checkVals[ck], "got:", v.capacity)
			t.Fatal()
		}
	}
}

func TestRunHeader(t *testing.T) {
	s := NewSession(Context{})

	results, err := s.Run("_data/dimacsMaxf.txt", "my customer header")
	if err != nil {
		t.Fatal(err)
	}

	check := "c my customer header"
	if results[0] != check {
		fmt.Println("wanted:", check, "got:", results[0])
		t.Fatal()
	}
}

// LowestLabel == false, FifoBuckets == false
func TestRunCase1(t *testing.T) {
	s := NewSession(Context{})

	results, err := s.Run("_data/dimacsMaxf.txt")
	if err != nil {
		t.Fatal(err)
	}

	fh, _ := os.Open("_data/dimacsMaxf.txt")
	defer fh.Close()
	input, _ := ioutil.ReadAll(fh)
	fmt.Println("input:")
	fmt.Println(string(input))

	for _, v := range results {
		fmt.Println(v)
	}

	fmt.Println("\nstats:", s.StatsJSON())
	fmt.Println("timer:", s.TimerJSON())
}

// LowestLabel == true, FifoBuckets == false
func TestRunCase2(t *testing.T) {
	s := NewSession(Context{LowestLabel:true})

	results, err := s.Run("_data/dimacsMaxf.txt")
	if err != nil {
		t.Fatal(err)
	}

		for _, v := range results {
		fmt.Println(v)
	}

	fmt.Println("\nstats:", s.StatsJSON())
	fmt.Println("timer:", s.TimerJSON())
}

// LowestLabel == false, FifoBuckets == true
func TestRunCase3(t *testing.T) {
	s := NewSession(Context{FifoBuckets:true})

	results, err := s.Run("_data/dimacsMaxf.txt")
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range results {
		fmt.Println(v)
	}

	fmt.Println("\nstats:", s.StatsJSON())
	fmt.Println("timer:", s.TimerJSON())
}

// LowestLabel == true, FifoBuckets == true
func TestRunCase4(t *testing.T) {
	s := NewSession(Context{LowestLabel:true,FifoBuckets:true})

	results, err := s.Run("_data/dimacsMaxf.txt")
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range results {
		fmt.Println(v)
	}

	fmt.Println("\nstats:", s.StatsJSON())
	fmt.Println("timer:", s.TimerJSON())
}

// Report cut set rather than flows
func TestRunDisplayCut (t *testing.T) {
	s := NewSession(Context{DisplayCut:true})

	results, err := s.Run("_data/dimacsMaxf.txt")
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range results {
		fmt.Println(v)
	}
}

func TestRunJSON(t *testing.T) {
	s := NewSession(Context{})

	results, err := s.RunJSON("_data/dimacsMaxf.txt")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(results))
}

