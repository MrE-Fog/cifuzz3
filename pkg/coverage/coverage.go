package coverage

import (
	"fmt"
	"io"

	"github.com/pterm/pterm"

	"code-intelligence.com/cifuzz/pkg/log"
	"code-intelligence.com/cifuzz/util/fileutil"
)

type CoverageSummary struct {
	Total *Coverage
	Files []*FileCoverage
}

type FileCoverage struct {
	Filename string
	Coverage *Coverage
}

type Coverage struct {
	FunctionsFound int
	FunctionsHit   int
	BranchesFound  int
	BranchesHit    int
	LinesFound     int
	LinesHit       int
}

func (cs *CoverageSummary) PrintTable(writer io.Writer) {
	formatCell := func(hit, found int) string {
		percent := 100.0
		if found != 0 {
			percent = (float64(hit) * 100) / float64(found)
		}
		return fmt.Sprintf("%d / %d %8s", hit, found, fmt.Sprintf("(%.1f%%)", percent))
	}

	// create table data for pterm table
	tableData := pterm.TableData{{"File", "Functions Hit/Found", "Lines Hit/Found", "Branches Hit/Found"}}
	for _, file := range cs.Files {
		tableData = append(tableData, []string{
			fileutil.PrettifyPath(file.Filename),
			formatCell(file.Coverage.FunctionsHit, file.Coverage.FunctionsFound),
			formatCell(file.Coverage.LinesHit, file.Coverage.LinesFound),
			formatCell(file.Coverage.BranchesHit, file.Coverage.BranchesFound),
		},
		)
	}
	tableData = append(tableData, []string{"", "", "", ""})
	// repeat the header for the case that the original header scrolled
	// of the terminal window
	tableData = append(tableData, []string{"", "Functions Hit/Found", "Lines Hit/Found", "Branches Hit/Found"})
	tableData = append(tableData, []string{
		"Total",
		fmt.Sprintf("%d / %d", cs.Total.FunctionsHit, cs.Total.FunctionsFound),
		fmt.Sprintf("%d / %d", cs.Total.LinesHit, cs.Total.LinesFound),
		fmt.Sprintf("%d / %d", cs.Total.BranchesHit, cs.Total.BranchesFound),
	},
	)
	table := pterm.DefaultTable.WithWriter(writer).WithHasHeader().WithData(tableData).WithRightAlignment()

	log.Print("\n")
	log.Successf("Coverage Report:\n")
	if err := table.Render(); err != nil {
		log.Error(err, "Unable to print coverage table")
	}
	log.Print("\n")
}
