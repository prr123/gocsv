// library to parse csv files

package csvLib

import (
	"fmt"
	"strings"
//	"encoding/csv"
)

type TRow []interface{}

type Table struct {
	Name string
	Count int
	Header []string
	ColTyp []string
	Rows []TRow
}


var ColTyp = [...]string {"int","float","bool", "string", "text" }

func Table2Md(tbl Table) ([]byte, error) {

	md:= make([]byte, 1024)
	ipos := 0

	numCol := len(tbl.Header)
	if numCol == 0 {return nil, fmt.Errorf("no columns")}
	// header
	colStr:="|"
	for i:=0; i< numCol-1; i++ {
		colStr=colStr + tbl.Header[i]+"|"
	}
	colStr +="\n"
//	fmt.Printf("dbg -- colStr->%s<-",colStr)
	copy(md[ipos:len(colStr)+1], []byte(colStr))
	ipos += len(colStr)

	fmtStr:="+"
	for i:=0; i< numCol-1; i++ {
		fmtStr=fmtStr + "---+"
	}
	fmtStr +="\n"
//	fmt.Printf("dbg -- fmtStr[%d]->%s",len(fmtStr),fmtStr)
//	fmt.Printf("len %d ipos %d\n", len(md), ipos)
	copy(md[ipos:ipos+len(fmtStr)+1], []byte(fmtStr))
	ipos += len(fmtStr)

	//table data
	for irow:=0; irow< len(tbl.Rows); irow++ {
		trow := tbl.Rows[irow]
		trStr := "|"
		for icol:=0; icol< numCol-1; icol++ {
			trStr = trStr + fmt.Sprintf("%v",trow[icol]) + "|"
		}
		trStr += "\n"
		copy(md[ipos:ipos+len(trStr)+1], []byte(trStr))
		ipos += len(trStr)
	}

	//print name
	nam:=fmt.Sprintf("Table %d: %s  \n", tbl.Count+1, tbl.Name)
	copy(md[ipos:ipos+len(nam)+1], []byte(nam))
	ipos += len(nam)
	return md, nil
}


func GetTables(lines []string)([]Table, error) {

	var TblList []Table

	Irow :=0
	istate :=0
	tbl := Table{}
	trow := TRow{}
	cells:= strings.Split(lines[0], ",")
	numCol := len(cells)
	for ilin:=0; ilin< len(lines); ilin++ {
		cells:= strings.Split(lines[ilin], ",")
//		if len(cells) != numCols {return nil, fmt.Errorf("number of columns changed in line: %d", ilin)
//fmt.Printf("dbg -- %v\n", cells)
		switch istate {
		case 0:
			if len(cells[0]) > 0 {
				isName := true
				if len(cells) > 1 {
					for j:= 1; j<numCol; j++ {
						if len(cells[j])>0 {
							isName = false
							break
						}
					}
					if isName {
						tbl.Name = cells[0]
						istate = 1
						break
					}

				}
			} else {break}
			fallthrough

		// headings
		case 1:
			if cells[0] == "header" {
				tbl.Header = make([]string, numCol)
				for j:= 1; j< numCol; j++ {
					if len(cells[j]) == 0 {
						numCol = j-1
						break
					}
					tbl.Header[j-1] = cells[j]
				}
				tbl.Header = tbl.Header[:numCol]
				istate = 2
				break
			}
			fallthrough

		// types
		case 2:
			if cells[0] == "type" {
				tbl.ColTyp = make([]string, numCol)
				for j:= 1; j< numCol; j++ {
					tbl.ColTyp[j-1] = cells[j]
				}
				istate = 3
				break
			}
			fallthrough



		// table call values
		case 3:
			if len(cells[0]) > 0 {
				return nil, fmt.Errorf("non-empty first cell in line: %d", ilin)
			}
			trow = make([]interface{}, numCol-1)
			isEmpty := true
			for j:= 1; j< numCol; j++ {
				if len(cells[j]) > 0 {isEmpty = false}
				trow[j-1]=cells[j]
			}
//fmt.Printf("dbg -- trow2: %v\n", trow[:len(cells)-1])
			if !isEmpty {
				tbl.Rows = append(tbl.Rows, trow)
//fmt.Printf("dbg -- Row[%d]: %v\n",Irow, tbl.Rows[Irow])
//fmt.Printf("dbg -- %v\n", tbl)
				Irow++
				break
			}
			fallthrough

		// fini
		case 4:
			istate = 0
			TblList = append(TblList, tbl)
//fmt.Printf("dbg -- %v\n", TblList[0])
			tbl = Table{}
		default:
			return nil, fmt.Errorf("invalid state: %d",istate)
		}
	}
	return TblList, nil
}

func ProcTable(src []byte)([]string, error) {

    var lines []string

    lineSt :=0
    for i:=0; i<len(src); i++ {
        if src[i] == '\n' {
            lines = append(lines, string(src[lineSt:i]))
            lineSt = i+1
            i += 1
        }
    }
    return lines, nil
}


func PrintLines(lines []string) {

    fmt.Println("*************** Lines ******************")
    for i:=0; i< len(lines); i++ {
        fmt.Printf("--%d: %s\n",i, lines[i])
    }
    fmt.Println("************* End Lines ****************")
}

func PrintTable(tbl Table) {
	fmt.Printf("****************** Table: %s *************\n", tbl.Name)
	Cols:=len(tbl.Rows[0])
	fmt.Printf("Table Name: %s [%d:%d]\n", tbl.Name,len(tbl.Rows), Cols)
	if tbl.Header != nil {
		fmt.Printf("Headings: |")
		for i:=1; i< len(tbl.Header); i++ {fmt.Printf("%s |", tbl.Header[i])}
		fmt.Printf("\n")
	}
	if tbl.ColTyp != nil {
		fmt.Printf("Cols: |")
		for i:=0; i< Cols; i++ {fmt.Printf("%s |", tbl.ColTyp[i])}
		fmt.Printf("\n")
	}
	for irow:=0; irow< len(tbl.Rows); irow++ {
		fmt.Printf("--%d: |", irow)
		trow := tbl.Rows[irow]
		for icol:=0; icol< Cols; icol++ {
			fmt.Printf("%v |", trow[icol])
		}
		fmt.Printf("\n")
//fmt.Printf("dbg -- trow [%d]: %v\n", irow, trow)
	}
	fmt.Printf("****************** End Table *************\n")
}

