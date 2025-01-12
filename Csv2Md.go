package main

import (
	"fmt"
	"os"
    "log"

//    "encoding/csv"

	csvLib "goDemo/goCsv/csvLib"
	util "github.com/prr123/utility/utilLib"
)

func main() {

    numarg := len(os.Args)
    flags:=[]string{"dbg", "in", "out"}

    useStr := " /in=infile /out=outfile [/dbg]"
    helpStr := "markdown to html conversion program"

    if numarg > len(flags) +1 {
        fmt.Println("too many arguments in cl!")
        fmt.Println("usage: %s %s\n", os.Args[0], useStr)
        os.Exit(-1)
    }

    if numarg == 1 || (numarg > 1 && os.Args[1] == "help") {
        fmt.Printf("help: %s\n", helpStr)
        fmt.Printf("usage is: %s %s\n", os.Args[0], useStr)
        os.Exit(1)
    }

    flagMap, err := util.ParseFlags(os.Args, flags)
    if err != nil {log.Fatalf("util.ParseFlags: %v\n", err)}


    dbg:= false
    _, ok := flagMap["dbg"]
    if ok {dbg = true}

    inFil := ""
    inval, ok := flagMap["in"]
    if !ok {
        log.Fatalf("error -- no in flag provided!\n")
    } else {
        if inval.(string) == "none" {log.Fatalf("error -- no input file name provided!\n")}
        inFil = inval.(string)
    }

    outFil := ""
    outval, ok := flagMap["out"]
    if !ok {
        outFil = inFil
    } else {
        if outval.(string) == "none" {
            outFil = inFil
        } else {
            outFil = outval.(string)
        }
    }

    inFilnam := "csv/" + inFil + ".csv"
    outFilnam := "data/" + outFil + ".dat"

    if dbg {
        fmt.Printf("input:  %s\n", inFilnam)
        fmt.Printf("output: %s\n", outFilnam)
    }


	inData, err := os.ReadFile(inFilnam)
	if err != nil {log.Fatalf("error -- read file: %v\n", err)}

	if dbg {
		fmt.Println("*****************************")
		fmt.Printf("%s\n", inData)
		fmt.Println("*****************************")
	}

	lines, err := csvLib.ProcTable(inData)
	if err != nil {log.Fatalf("error -- ProcTable: %v\n", err)}

	csvLib.PrintLines(lines)

	tables, err := csvLib.GetTables(lines)
	if err != nil {log.Fatalf("error -- ProcTable: %v\n", err)}

	fmt.Printf("tables: %d\n", len(tables))
	for i:=0; i< len(tables); i++ {
		csvLib.PrintTable(tables[i])
	}

	mdData, err := csvLib.Table2Md(tables[0])
	if err != nil {log.Fatalf("error -- Table2Md: %v\n", err)}

	fmt.Printf("*** md ****\n%s\n",mdData)
	fmt.Println("success")
}
