package rplanlib

import (
	"fmt"
	"math"
	"os"
	//"regexp"
	"strings"
	"testing"
	"time"

	"github.com/willauld/lpsimplex"
)

//
// Testing for results_output.go
//

var sipSingle = map[string]string{
	"setName":                    "activeParams",
	"filingStatus":               "single",
	"key1":                       "retiree1",
	"key2":                       "",
	"eT_Age1":                    "54",
	"eT_Age2":                    "",
	"eT_RetireAge1":              "65",
	"eT_RetireAge2":              "",
	"eT_PlanThroughAge1":         "75",
	"eT_PlanThroughAge2":         "",
	"eT_PIA1":                    "",
	"eT_PIA2":                    "",
	"eT_SS_Start1":               "",
	"eT_SS_Start2":               "",
	"eT_TDRA1":                   "200", // 200k
	"eT_TDRA2":                   "",
	"eT_TDRA_Rate1":              "",
	"eT_TDRA_Rate2":              "",
	"eT_TDRA_Contrib1":           "",
	"eT_TDRA_Contrib2":           "",
	"eT_TDRA_ContribStartAge1":   "",
	"eT_TDRA_ContribStartAge2":   "",
	"eT_TDRA_ContribEndAge1":     "",
	"eT_TDRA_ContribEndAge2":     "",
	"eT_Roth1":                   "",
	"eT_Roth2":                   "",
	"eT_Roth_Rate1":              "",
	"eT_Roth_Rate2":              "",
	"eT_Roth_Contrib1":           "",
	"eT_Roth_Contrib2":           "",
	"eT_Roth_ContribStartAge1":   "",
	"eT_Roth_ContribStartAge2":   "",
	"eT_Roth_ContribEndAge1":     "",
	"eT_Roth_ContribEndAge2":     "",
	"eT_Aftatax":                 "",
	"eT_Aftatax_Rate":            "",
	"eT_Aftatax_Contrib":         "",
	"eT_Aftatax_ContribStartAge": "",
	"eT_Aftatax_ContribEndAge":   "",

	"eT_iRatePercent":    "2.5",
	"eT_rRatePercent":    "6",
	"eT_maximize":        "Spending", // or "PlusEstate"
	"dollarsInThousands": "true",
}

var sipJoint = map[string]string{
	"setName":                    "activeParams",
	"filingStatus":               "joint",
	"key1":                       "retiree1",
	"key2":                       "retiree2",
	"eT_Age1":                    "54",
	"eT_Age2":                    "54",
	"eT_RetireAge1":              "65",
	"eT_RetireAge2":              "65",
	"eT_PlanThroughAge1":         "75",
	"eT_PlanThroughAge2":         "75",
	"eT_PIA1":                    "",
	"eT_PIA2":                    "",
	"eT_SS_Start1":               "",
	"eT_SS_Start2":               "",
	"eT_TDRA1":                   "200", // 200k
	"eT_TDRA2":                   "",
	"eT_TDRA_Rate1":              "",
	"eT_TDRA_Rate2":              "",
	"eT_TDRA_Contrib1":           "",
	"eT_TDRA_Contrib2":           "",
	"eT_TDRA_ContribStartAge1":   "",
	"eT_TDRA_ContribStartAge2":   "",
	"eT_TDRA_ContribEndAge1":     "",
	"eT_TDRA_ContribEndAge2":     "",
	"eT_Roth1":                   "",
	"eT_Roth2":                   "",
	"eT_Roth_Rate1":              "",
	"eT_Roth_Rate2":              "",
	"eT_Roth_Contrib1":           "",
	"eT_Roth_Contrib2":           "",
	"eT_Roth_ContribStartAge1":   "",
	"eT_Roth_ContribStartAge2":   "",
	"eT_Roth_ContribEndAge1":     "",
	"eT_Roth_ContribEndAge2":     "",
	"eT_Aftatax":                 "",
	"eT_Aftatax_Rate":            "",
	"eT_Aftatax_Contrib":         "",
	"eT_Aftatax_ContribStartAge": "",
	"eT_Aftatax_ContribEndAge":   "",

	"eT_iRatePercent":    "2.5",
	"eT_rRatePercent":    "6",
	"eT_maximize":        "Spending", // or "PlusEstate"
	"dollarsInThousands": "true",
}
var sipSingle3Acc = map[string]string{
	"setName":                    "activeParams",
	"filingStatus":               "single",
	"key1":                       "retiree1",
	"key2":                       "",
	"eT_Age1":                    "54",
	"eT_Age2":                    "",
	"eT_RetireAge1":              "65",
	"eT_RetireAge2":              "",
	"eT_PlanThroughAge1":         "75",
	"eT_PlanThroughAge2":         "",
	"eT_PIA1":                    "",
	"eT_PIA2":                    "",
	"eT_SS_Start1":               "",
	"eT_SS_Start2":               "",
	"eT_TDRA1":                   "200", // 200k
	"eT_TDRA2":                   "",
	"eT_TDRA_Rate1":              "",
	"eT_TDRA_Rate2":              "",
	"eT_TDRA_Contrib1":           "",
	"eT_TDRA_Contrib2":           "",
	"eT_TDRA_ContribStartAge1":   "",
	"eT_TDRA_ContribStartAge2":   "",
	"eT_TDRA_ContribEndAge1":     "",
	"eT_TDRA_ContribEndAge2":     "",
	"eT_Roth1":                   "10", //10K
	"eT_Roth2":                   "",
	"eT_Roth_Rate1":              "",
	"eT_Roth_Rate2":              "",
	"eT_Roth_Contrib1":           "",
	"eT_Roth_Contrib2":           "",
	"eT_Roth_ContribStartAge1":   "",
	"eT_Roth_ContribStartAge2":   "",
	"eT_Roth_ContribEndAge1":     "",
	"eT_Roth_ContribEndAge2":     "",
	"eT_Aftatax":                 "50", //50k
	"eT_Aftatax_Rate":            "",
	"eT_Aftatax_Contrib":         "",
	"eT_Aftatax_ContribStartAge": "",
	"eT_Aftatax_ContribEndAge":   "",

	"eT_iRatePercent":    "2.5",
	"eT_rRatePercent":    "6",
	"eT_maximize":        "Spending", // or "PlusEstate"
	"dollarsInThousands": "true",
}

//def precheck_consistancy():
func TestPreCheckConsistancy(t *testing.T) {
	fmt.Printf("TestPreCheckConsistancy() Not Yet Implemented\n")
}

func TestCheckConsistancy(t *testing.T) {
	fmt.Printf("TestCheckConsistancy() Not Yet Implemented\n")
}

//func (ms ModelSpecs) activitySummaryHeader(fieldwidth int)
func TestActivitySummaryHeader(t *testing.T) {
	tests := []struct {
		sip    map[string]string
		expect string
	}{
		{
			sip: sipJoint,
			expect: `retiree1/retiree2
    age     fIRA    tIRA  RMDref   fRoth   tRoth fAftaTx tAftaTx   o_inc      SS NetASal Expense TFedTax Spndble`,
		},
		{
			sip: sipSingle,
			expect: `retiree1
 age     fIRA    tIRA  RMDref   fRoth   tRoth fAftaTx tAftaTx   o_inc      SS NetASal Expense TFedTax Spndble`,
		},
	}
	for i, elem := range tests {
		ip, err := NewInputParams(elem.sip, nil)
		if err != nil {
			t.Errorf("TestActivitySummaryHeader case %d: %s", i, err)
			continue
		}
		//fmt.Printf("InputParams: %#v\n", ip)
		ti := NewTaxInfo(ip.FilingStatus, 2017)
		taxbins := len(*ti.Taxtable)
		cgbins := len(*ti.Capgainstable)
		vindx, err := NewVectorVarIndex(ip.Numyr, taxbins,
			cgbins, ip.Accmap, os.Stdout)
		if err != nil {
			t.Errorf("TestActivitySummaryHeader case %d: %s", i, err)
			continue
		}
		csvfile := (*os.File)(nil)
		logfile := os.Stdout
		RoundToOneK := false
		ms, err := NewModelSpecs(vindx, ti, *ip,
			RoundToOneK, false, false,
			os.Stderr, logfile, csvfile, logfile, nil)
		if err != nil {
			t.Errorf("TestActivitySummaryHeader case %d: %s", i, err)
			continue
		}

		mychan := make(chan string)
		DoNothing := false //true
		oldout, w, err := ms.RedirectModelSpecsTable(mychan, DoNothing)
		if err != nil {
			t.Errorf("RedirectModelSpecsTable: %s\n", err)
			return // should this be continue?
		}

		fieldwidth := 7
		ms.activitySummaryHeader(fieldwidth)

		str := ms.RestoreModelSpecsTable(mychan, oldout, w, DoNothing)
		strn := strings.TrimSpace(str)
		//strn := stripWhitespace(str)
		//warningmes := stripWhitespace(elem.warningmes)
		if elem.expect != strn {
			t.Errorf("TestActivitySummaryHeader case %d:  expected output:\n\t '%s'\n\tbut found:\n\t'%s'\n", i, elem.expect, strn)
		}
	}
}

//func (ms ModelSpecs) printActivitySummary(xp *[]float64)
func TestActivitySummary(t *testing.T) {
	tests := []struct {
		expect string
		sip    map[string]string
		xp     *[]float64
	}{
		{ // Case 0
			expect: `Activity Summary:

 retiree1
 age     fIRA    tIRA  RMDref   fRoth   tRoth fAftaTx tAftaTx   o_inc      SS NetASal Expense TFedTax Spndble
  65:   40594       0       0       0       0       0       0       0       0       0       0    3431   37164 
  66:   41609       0       0       0       0       0       0       0       0       0       0    3516   38093 
  67:   42650       0       0       0       0       0       0       0       0       0       0    3604   39045 
  68:   43716       0       0       0       0       0       0       0       0       0       0    3694   40021 
  69:   44809       0       0       0       0       0       0       0       0       0       0    3787   41022 
  70:   45929       0    9263       0       0       0       0       0       0       0       0    3881   42048 
  71:   47077       0    8315       0       0       0       0       0       0       0       0    3978   43099 
  72:   48254       0    7174       0       0       0       0       0       0       0       0    4078   44176 
  73:   49460       0    5811       0       0       0       0       0       0       0       0    4180   45281 
  74:   50697       0    4190       0       0       0       0       0       0       0       0    4284   46413 
  75:   51964       0    2269       0       0       0       0       0       0       0       0    4391   47573 
 retiree1
 age     fIRA    tIRA  RMDref   fRoth   tRoth fAftaTx tAftaTx   o_inc      SS NetASal Expense TFedTax Spndble`,
			sip: sipSingle,
			xp:  xpSingle,
		},
		{ // Case 1
			expect: `Activity Summary:

retiree1/retiree2
    age     fIRA    tIRA  RMDref   fRoth   tRoth fAftaTx tAftaTx   o_inc      SS NetASal Expense TFedTax Spndble
 65/ 65:   40594       0       0       0       0       0       0       0       0       0       0    1330   39264 
 66/ 66:   41609       0       0       0       0       0       0       0       0       0       0    1364   40246 
 67/ 67:   42650       0       0       0       0       0       0       0       0       0       0    1398   41252 
 68/ 68:   43716       0       0       0       0       0       0       0       0       0       0    1433   42283 
 69/ 69:   44809       0       0       0       0       0       0       0       0       0       0    1468   43340 
 70/ 70:   45929       0    9263       0       0       0       0       0       0       0       0    1505   44424 
 71/ 71:   47077       0    8315       0       0       0       0       0       0       0       0    1543   45534 
 72/ 72:   48254       0    7174       0       0       0       0       0       0       0       0    1581   46673 
 73/ 73:   49460       0    5811       0       0       0       0       0       0       0       0    1621   47840 
 74/ 74:   50697       0    4190       0       0       0       0       0       0       0       0    1661   49036 
 75/ 75:   51964       0    2269       0       0       0       0       0       0       0       0    1703   50261 
retiree1/retiree2
    age     fIRA    tIRA  RMDref   fRoth   tRoth fAftaTx tAftaTx   o_inc      SS NetASal Expense TFedTax Spndble`,
			sip: sipJoint,
			xp:  xpJoint,
		},
	}
	for i, elem := range tests {
		//fmt.Printf("================ CASE %d ==================\n", i)
		ip, err := NewInputParams(elem.sip, nil)
		if err != nil {
			t.Errorf("TestActivitySummary case %d: %s", i, err)
			continue
		}
		//fmt.Printf("InputParams: %#v\n", ip)
		ti := NewTaxInfo(ip.FilingStatus, 2017)
		taxbins := len(*ti.Taxtable)
		cgbins := len(*ti.Capgainstable)
		vindx, err := NewVectorVarIndex(ip.Numyr, taxbins,
			cgbins, ip.Accmap, os.Stdout)
		if err != nil {
			t.Errorf("TestActivitySummary case %d: %s", i, err)
			continue
		}
		csvfile := (*os.File)(nil)
		logfile := os.Stdout
		RoundToOneK := false
		ms, err := NewModelSpecs(vindx, ti, *ip,
			RoundToOneK, false, false,
			os.Stderr, logfile, csvfile, logfile, nil)
		if err != nil {
			t.Errorf("TestActivitySummaryHeader case %d: %s", i, err)
			continue
		}

		mychan := make(chan string)
		donothing := false
		oldout, w, err := ms.RedirectModelSpecsTable(mychan, donothing)
		if err != nil {
			t.Errorf("RedirectModelSpecsTable: %s\n", err)
			return // should this be continue?
		}

		//xp := &[]float64{0.0, 0.0}
		ms.PrintActivitySummary(elem.xp)

		str := ms.RestoreModelSpecsTable(mychan, oldout, w, donothing)
		strnn := strings.TrimSpace(str)
		expect := elem.expect
		//re := regexp.MustCompile("\r")
		//strnn := re.ReplaceAllString(strn, "")
		//expect := re.ReplaceAllString(elem.expect, "")
		//r:=NewReplacer(U+0010, '')
		//strnn := r.Replace(strn)
		//expect := r.Replace(elem.expect)
		//strn := stripWhitespace(str)
		//warningmes := stripWhitespace(elem.warningmes)
		if expect != strnn {
			str := fmt.Sprintf("TestActivitySummary Case %d:", i)
			showStrMismatch(str, expect, strnn)
			t.Errorf("TestActivitySummary case %d:  expected output:\n\t '%s'\n\tbut found:\n\t'%s'\n", i, elem.expect, strnn)
		}
	}
}

//func (ms ModelSpecs) printIncomeHeader(headerkeylist []string, countlist []int, incomeCat []string, fieldwidth int)
func TestPrintIncomeHeader(t *testing.T) {
	tests := []struct {
		sip        map[string]string
		expect     string
		headerlist []string
		countlist  []int
		csvfile    *os.File
		tablefile  *os.File
	}{
		{ //case 0
			sip:       sipJoint,
			countlist: []int{0, 0, 0, 0},
			expect: `retiree1
    age`,
			csvfile:   (*os.File)(nil),
			tablefile: os.Stdout,
		},
		{ //case 1
			sip:       sipSingle,
			countlist: []int{0, 0, 0, 0},
			expect: `retir
 age`,
			csvfile:   (*os.File)(nil),
			tablefile: os.Stdout,
		},
		{ //case 2
			sip:        sipJoint,
			countlist:  []int{1, 0, 0, 0},
			headerlist: []string{"nokey"},
			expect: `retiree1 SSincome
    age`,
			csvfile:   (*os.File)(nil),
			tablefile: os.Stdout,
		},
		{ //case 3
			sip:       sipJoint,
			countlist: []int{3, 3, 3, 3},
			headerlist: []string{
				"SS1", "SS2", "SS3", "income1", "income2",
				"income3", "asset1", "asset2", "asset3",
				"expense1", "expense2", "expense3",
			},
			expect: `retiree1 SSincome:                  Income:                    AssetSale:                 Expense:                  
    age       SS1      SS2      SS3  income1  income2  income3   asset1   asset2   asset3 expense1 expense2 expense3`,
			csvfile:   (*os.File)(nil),
			tablefile: os.Stdout,
		},
	}
	for i, elem := range tests {
		//fmt.Printf("=============== Case %d =================\n", i)
		ip, err := NewInputParams(elem.sip, nil)
		if err != nil {
			t.Errorf("TestPrintIncomeHeader case %d: %s", i, err)
			continue
		}
		ms := ModelSpecs{
			Ip:        *ip,
			Logfile:   os.Stdout,
			Errfile:   os.Stderr,
			Ao:        NewAppOutput(elem.csvfile, elem.tablefile),
			AssetSale: make([][]float64, 0),
		}
		ms.AssetSale = append(ms.AssetSale, make([]float64, ip.Numyr))

		mychan := make(chan string)
		DoNothing := false //true
		oldout, w, err := ms.RedirectModelSpecsTable(mychan, DoNothing)
		if err != nil {
			t.Errorf("RedirectModelSpecsTable: %s\n", err)
			return // should this be continue?
		}

		incomeCat := []string{"SSincome:", "Income:", "AssetSale:", "Expense:"}
		fieldwidth := 8
		ms.printIncomeHeader(elem.headerlist, elem.countlist, incomeCat, fieldwidth)

		str := ms.RestoreModelSpecsTable(mychan, oldout, w, DoNothing)
		strn := strings.TrimSpace(str)
		if elem.expect != strn {
			str := fmt.Sprintf("TestPrintIncomeHeader Case %d:", i)
			showStrMismatch(str, elem.expect, strn)
			t.Errorf("TestPrintIncomeHeader case %d:  expected output:\n\t '%s'\n\tbut found:\n\t'%s'\n", i, elem.expect, strn)
		}
	}
}

//func (ms ModelSpecs) getSSIncomeAssetExpenseList() ([]string, []int, [][]float64)
func TestGetIncomeAssetExpenseList(t *testing.T) {
	tests := []struct {
		sip            map[string]string
		incomeStreams  int
		expenseStreams int
		SSStreams      int
		AssetStreams   int
	}{
		{ //case 0
			incomeStreams:  3,
			expenseStreams: 3,
			SSStreams:      3,
			AssetStreams:   3,
			sip:            sipJoint,
		},
		{ //case 1
			incomeStreams:  1,
			expenseStreams: 0,
			SSStreams:      2,
			AssetStreams:   4,
			sip:            sipJoint,
		},
	}
	for i, elem := range tests {
		//fmt.Printf("=============== Case %d =================\n", i)
		ip, err := NewInputParams(elem.sip, nil)
		if err != nil {
			t.Errorf("TestGetIncomeAssetExpenseList case %d: %s", i, err)
			continue
		}
		ms := ModelSpecs{
			//ip:      ip,
			//logfile: os.Stdout,
			//errfile: os.Stderr,
			//ao:      NewAppOutput(elem.csvfile, elem.tablefile),

			SS:     make([][]float64, 0),
			SStags: make([]string, 0),

			Income:     make([][]float64, 0),
			Incometags: make([]string, 0),

			AssetSale: make([][]float64, 0),
			Assettags: make([]string, 0),

			Expenses:    make([][]float64, 0),
			Expensetags: make([]string, 0),
		}
		for i := 0; i <= elem.SSStreams; i++ {
			ms.SS = append(ms.SS, make([]float64, ip.Numyr))
			str := fmt.Sprintf("SS%d", i)
			ms.SStags = append(ms.SStags, str)
		}
		for i := 0; i <= elem.incomeStreams; i++ {
			ms.Income = append(ms.Income, make([]float64, ip.Numyr))
			str := fmt.Sprintf("income%d", i)
			ms.Incometags = append(ms.Incometags, str)
		}
		for i := 0; i <= elem.AssetStreams; i++ {
			ms.AssetSale = append(ms.AssetSale, make([]float64, ip.Numyr))
			str := fmt.Sprintf("asset%d", i)
			ms.Assettags = append(ms.Assettags, str)
		}
		for i := 0; i <= elem.expenseStreams; i++ {
			ms.Expenses = append(ms.Expenses, make([]float64, ip.Numyr))
			str := fmt.Sprintf("expense%d", i)
			ms.Expensetags = append(ms.Expensetags, str)
		}
		ms.AssetSale[1][7] = 50000

		headerlist, countlist, matrix := ms.getSSIncomeAssetExpenseList()
		//fmt.Printf("headerlist: %#v\n", headerlist)
		//fmt.Printf("countlist: %#v\n", countlist)
		//fmt.Printf("matrix: %#v\n", matrix)
		htot := elem.SSStreams + elem.incomeStreams + elem.AssetStreams + elem.expenseStreams
		if htot != len(headerlist) {
			t.Errorf("TestGetIncomeAssetExpenseList case %d: expected %d headers but found %d\n", i, htot, len(headerlist))
		}
		if htot != len(matrix) {
			t.Errorf("TestGetIncomeAssetExpenseList case %d: expected %d vectors but found %d\n", i, htot, len(matrix))
		}
		if elem.SSStreams != countlist[0] {
			t.Errorf("TestGetIncomeAssetExpenseList case %d:  expected %d SS streams but found %d streams\n", i, elem.SSStreams, countlist[0])
		}
		if elem.incomeStreams != countlist[1] {
			t.Errorf("TestGetIncomeAssetExpenseList case %d:  expected %d income streams but found %d streams\n", i, elem.SSStreams, countlist[0])
		}
		if elem.AssetStreams != countlist[2] {
			t.Errorf("TestGetIncomeAssetExpenseList case %d:  expected %d asset streams but found %d streams\n", i, elem.SSStreams, countlist[0])
		}
		if elem.expenseStreams != countlist[3] {
			t.Errorf("TestGetIncomeAssetExpenseList case %d:  expected %d expense streams but found %d streams\n", i, elem.SSStreams, countlist[0])
		}
	}
}

//func (ms ModelSpecs) printIncomeExpenseDetails()
func TestPrintIncomeExpenseDetails(t *testing.T) {
	tests := []struct {
		sip            map[string]string
		incomeStreams  int
		expenseStreams int
		SSStreams      int
		AssetStreams   int
		expect         string
		onek           float64
	}{
		{ //case 0
			incomeStreams:  3,
			expenseStreams: 3,
			SSStreams:      3,
			AssetStreams:   3,
			sip:            sipJoint,
			expect: `Income and Expense Summary:

retiree1 SSincome:                  Income:                    AssetSale:                 Expense:                  
    age       SS1      SS2      SS3  income1  income2  income3   asset1   asset2   asset3 expense1 expense2 expense3
 65/ 65:        1        2        3        1        2        3        1        2        3        1        2        3
 66/ 66:     1000     1000     1000     1000     1000     1000     1000     1000     1000     1000     1000     1000
 67/ 67:     2000     2000     2000     2000     2000     2000     2000     2000     2000     2000     2000     2000
 68/ 68:     3000     3000     3000     3000     3000     3000     3000     3000     3000     3000     3000     3000
 69/ 69:     4000     4000     4000     4000     4000     4000     4000     4000     4000     4000     4000     4000
 70/ 70:     5000     5000     5000     5000     5000     5000     5000     5000     5000     5000     5000     5000
 71/ 71:     6000     6000     6000     6000     6000     6000     6000     6000     6000     6000     6000     6000
 72/ 72:     7000     7000     7000     7000     7000     7000    50000     7000     7000     7000     7000     7000
 73/ 73:     8000     8000     8000     8000     8000     8000     8000     8000     8000     8000     8000     8000
 74/ 74:     9000     9000     9000     9000     9000     9000     9000     9000     9000     9000     9000     9000
 75/ 75:    10000    10000    10000    10000    10000    10000    10000    10000    10000    10000    10000    10000
retiree1 SSincome:                  Income:                    AssetSale:                 Expense:                  
    age       SS1      SS2      SS3  income1  income2  income3   asset1   asset2   asset3 expense1 expense2 expense3`,
			onek: 1,
		},
		{ //case 1
			incomeStreams:  1,
			expenseStreams: 0,
			SSStreams:      2,
			AssetStreams:   4,
			sip:            sipSingle,
			expect: `Income and Expense Summary:

retir SSincome:         Income:  AssetSale:                         
 age       SS1      SS2  income1   asset1   asset2   asset3   asset4
  65:        0        0        0        0        0        0        0
  66:        1        1        1        1        1        1        1
  67:        2        2        2        2        2        2        2
  68:        3        3        3        3        3        3        3
  69:        4        4        4        4        4        4        4
  70:        5        5        5        5        5        5        5
  71:        6        6        6        6        6        6        6
  72:        7        7        7       50        7        7        7
  73:        8        8        8        8        8        8        8
  74:        9        9        9        9        9        9        9
  75:       10       10       10       10       10       10       10
retir SSincome:         Income:  AssetSale:                         
 age       SS1      SS2  income1   asset1   asset2   asset3   asset4`,
			onek: 1000,
		},
	}
	for i, elem := range tests {
		//fmt.Printf("=============== Case %d =================\n", i)
		ip, err := NewInputParams(elem.sip, nil)
		if err != nil {
			t.Errorf("TestTestPrintIncomeExpenseDetails case %d: %s", i, err)
			continue
		}
		csvfile := (*os.File)(nil)
		tablefile := os.Stdout
		ms := ModelSpecs{
			Ip:      *ip,
			Logfile: os.Stdout,
			Errfile: os.Stderr,
			Ao:      NewAppOutput(csvfile, tablefile),

			SS:     make([][]float64, 0),
			SStags: make([]string, 0),

			Income:     make([][]float64, 0),
			Incometags: make([]string, 0),

			AssetSale: make([][]float64, 0),
			Assettags: make([]string, 0),

			Expenses:    make([][]float64, 0),
			Expensetags: make([]string, 0),

			OneK: elem.onek,
		}
		for i := 0; i <= elem.SSStreams; i++ {
			v := make([]float64, ip.Numyr)
			for j := 1; j < ip.Numyr; j++ {
				v[j] = float64(j * 1000)
			}
			v[0] = float64(i)
			ms.SS = append(ms.SS, v)
			str := fmt.Sprintf("SS%d", i)
			ms.SStags = append(ms.SStags, str)
		}
		for i := 0; i <= elem.incomeStreams; i++ {
			v := make([]float64, ip.Numyr)
			for j := 1; j < ip.Numyr; j++ {
				v[j] = float64(j * 1000)
			}
			v[0] = float64(i)
			ms.Income = append(ms.Income, v)
			str := fmt.Sprintf("income%d", i)
			ms.Incometags = append(ms.Incometags, str)
		}
		for i := 0; i <= elem.AssetStreams; i++ {
			v := make([]float64, ip.Numyr)
			for j := 1; j < ip.Numyr; j++ {
				v[j] = float64(j * 1000)
			}
			v[0] = float64(i)
			ms.AssetSale = append(ms.AssetSale, v)
			str := fmt.Sprintf("asset%d", i)
			ms.Assettags = append(ms.Assettags, str)
		}
		for i := 0; i <= elem.expenseStreams; i++ {
			v := make([]float64, ip.Numyr)
			for j := 1; j < ip.Numyr; j++ {
				v[j] = float64(j * 1000)
			}
			v[0] = float64(i)
			ms.Expenses = append(ms.Expenses, v)
			str := fmt.Sprintf("expense%d", i)
			ms.Expensetags = append(ms.Expensetags, str)
		}
		ms.AssetSale[1][7] = 50000.0

		//headerlist, countlist, matrix := ms.getSSIncomeAssetExpenseList()
		//fmt.Printf("headerlist: %#v\n", headerlist)
		//fmt.Printf("countlist: %#v\n", countlist)
		//fmt.Printf("matrix: %#v\n", matrix)

		mychan := make(chan string)
		DoNothing := false //true
		oldout, w, err := ms.RedirectModelSpecsTable(mychan, DoNothing)
		if err != nil {
			t.Errorf("RedirectModelSpecsTable: %s\n", err)
			return // should this be continue?
		}

		ms.PrintIncomeExpenseDetails()

		str := ms.RestoreModelSpecsTable(mychan, oldout, w, DoNothing)
		strn := strings.TrimSpace(str)
		if elem.expect != strn {
			str := fmt.Sprintf("TestTestPrintIncomeExpenseDetails Case %d:", i)
			showStrMismatch(str, elem.expect, strn)
			t.Errorf("TestTestPrintIncomeExpenseDetails case %d:  expected output:\n\t '%s'\n\tbut found:\n\t'%s'\n", i, elem.expect, strn)
		}
	}
}

//func printAccHeader()
func TestPrintAccHeader(t *testing.T) {
	tests := []struct {
		sip    map[string]string
		expect string
		onek   float64
	}{
		{ //case 0
			sip: sipJoint,
			expect: `retiree1/retiree2
    age      IRA    fIRA    tIRA  RMDref`,
			onek: 1,
		},
		{ //case 1
			sip: sipSingle,
			expect: `retiree1
 age      IRA    fIRA    tIRA  RMDref`,
			onek: 1,
		},
		{ //case 2
			sip: sipSingle3Acc,
			expect: `retiree1
 age      IRA    fIRA    tIRA  RMDref    Roth   fRoth   tRoth  AftaTx fAftaTx tAftaTx`,
			onek: 1,
		},
	}
	for i, elem := range tests {
		//fmt.Printf("=============== Case %d =================\n", i)
		ip, err := NewInputParams(elem.sip, nil)
		if err != nil {
			t.Errorf("TestTestPrintIncomeExpenseDetails case %d:  %s", i, err)
			continue
		}
		csvfile := (*os.File)(nil)
		tablefile := os.Stdout
		ms := ModelSpecs{
			Ip:      *ip,
			Logfile: os.Stdout,
			Errfile: os.Stderr,
			Ao:      NewAppOutput(csvfile, tablefile),
			OneK:    elem.onek,
		}

		mychan := make(chan string)
		DoNothing := false //true
		oldout, w, err := ms.RedirectModelSpecsTable(mychan, DoNothing)
		if err != nil {
			t.Errorf("RedirectModelSpecsTable: %s\n", err)
			return // should this be continue?
		}

		ms.printAccHeader()

		str := ms.RestoreModelSpecsTable(mychan, oldout, w, DoNothing)
		strn := strings.TrimSpace(str)
		if elem.expect != strn {
			str := fmt.Sprintf("TestPrintIncomeExpenseDetails Case %d:", i)
			showStrMismatch(str, elem.expect, strn)
			t.Errorf("TestPrintIncomeExpenseDetails case %d:  expected output:\n\t '%s'\n\tbut found:\n\t'%s'\n", i, elem.expect, strn)
		}
	}
}

//func (ms ModelSpecs) printAccountTrans(xp *[]float64)
func TestPrintAccountTrans(t *testing.T) {
	tests := []struct {
		sip    map[string]string
		sxp    *[]float64
		expect string
	}{
		{ //case 0
			sip: sipJoint,
			sxp: xpJoint,
			expect: `Account Transactions Summary:

retiree1/retiree2
    age      IRA    fIRA    tIRA  RMDref
 54/ 54:  200000       0       0       0
Plan Start: ---------
 65/ 65:  379660   40594       0       0
 66/ 66:  359409   41609       0       0
 67/ 67:  336868   42650       0       0
 68/ 68:  311871   43716       0       0
 69/ 69:  284245   44809       0       0
 70/ 70:  253803   45929       0    9263
 71/ 71:  220346   47077       0    8315
 72/ 72:  183665   48254       0    7174
 73/ 73:  143536   49460       0    5811
 74/ 74:   99720   50697       0    4190
 75/ 75:   51964   51964       0    2269
Plan End: -----------
 76/ 76:       0       0       0       0
retiree1/retiree2
    age      IRA    fIRA    tIRA  RMDref`,
		},
		{ //case 1
			sip: sipSingle,
			sxp: xpSingle,
			expect: `Account Transactions Summary:

retiree1
 age      IRA    fIRA    tIRA  RMDref
  54:  200000       0       0       0
Plan Start: ---------
  65:  379660   40594       0       0
  66:  359409   41609       0       0
  67:  336868   42650       0       0
  68:  311871   43716       0       0
  69:  284245   44809       0       0
  70:  253803   45929       0    9263
  71:  220346   47077       0    8315
  72:  183665   48254       0    7174
  73:  143536   49460       0    5811
  74:   99720   50697       0    4190
  75:   51964   51964       0    2269
Plan End: -----------
  76:       0       0       0       0
retiree1
 age      IRA    fIRA    tIRA  RMDref`,
		},
		{ //case 2
			sip: sipSingle3Acc,
			sxp: xpSingle3Acc,
			expect: `Account Transactions Summary:

retiree1
 age      IRA    fIRA    tIRA  RMDref    Roth   fRoth   tRoth  AftaTx fAftaTx tAftaTx
  54:  200000       0       0       0   10000       0       0   50000       0       0
Plan Start: ---------
  65:  379660   54922       0       0   18983       0       0   94915       0       0
  66:  344222   56295       0       0   20122       0       0  100610       0       0
  67:  305203   57702       0       0   21329       0       0  106646       0       0
  68:  262350   59145       0       0   22609       0       0  113045       0       0
  69:  215398   60623       0       0   23966       0       0  119828       0       0
  70:  164061   31531       0    5988   25404    1791       0  127018   24225       0
  71:  140481   30014       0    5301   25029       0       0  108960   28627       0
  72:  117095   30764       0    4574   26531       0       0   85153   29343       0
  73:   91511   31533       0    3705   28123       0       0   59159   30076       0
  74:   63576   32322       0    2671   29810       0       0   30828   30828       0
  75:   33130   33130       0    1447   31599   31599       0       0       0       0
Plan End: -----------
  76:       0       0       0       0       0       0       0       0       0       0
retiree1
 age      IRA    fIRA    tIRA  RMDref    Roth   fRoth   tRoth  AftaTx fAftaTx tAftaTx`,
		},
	}
	for i, elem := range tests {
		fmt.Printf("=============== Case %d =================\n", i)
		ip, err := NewInputParams(elem.sip, nil)
		if err != nil {
			fmt.Printf("TestPrintAccountTrans: %s\n", err)
		}
		ti := NewTaxInfo(ip.FilingStatus, 2017)
		taxbins := len(*ti.Taxtable)
		cgbins := len(*ti.Capgainstable)
		vindx, err := NewVectorVarIndex(ip.Numyr, taxbins,
			cgbins, ip.Accmap, os.Stdout)
		if err != nil {
			t.Errorf("TestPrintAccountTrans case %d: %s", i, err)
			continue
		}
		logfile := os.Stdout
		csvfile := (*os.File)(nil)
		RoundToOneK := false
		ms, err := NewModelSpecs(vindx, ti, *ip,
			RoundToOneK, false, false,
			os.Stderr, logfile, csvfile, logfile, nil)
		if err != nil {
			t.Errorf("TestActivityAccountTrans case %d: %s", i, err)
			continue
		}

		mychan := make(chan string)
		DoNothing := false //true
		oldout, w, err := ms.RedirectModelSpecsTable(mychan, DoNothing)
		if err != nil {
			t.Errorf("RedirectModelSpecsTable: %s\n", err)
			return // should this be continue?
		}

		ms.PrintAccountTrans(elem.sxp)

		str := ms.RestoreModelSpecsTable(mychan, oldout, w, DoNothing)
		strn := strings.TrimSpace(str)
		if elem.expect != strn {
			str := fmt.Sprintf("TestPrintAccountTrans Case %d:", i)
			showStrMismatch(str, elem.expect, strn)
			t.Errorf("TestPrintAccountTrans case %d:  expected output:\n\t '%s'\n\tbut found:\n\t'%s'\n", i, elem.expect, strn)
		}
	}
}

//func (ms ModelSpecs) printheaderTax()
func TestPrintHeaderTax(t *testing.T) {
	tests := []struct {
		sip    map[string]string
		expect string
	}{
		{ // Case 0
			sip: sipSingle,
			expect: `retiree1
 age     fIRA    tIRA  TxbleO TxbleSS  deduct   T_inc  earlyP  fedtax  mTaxB% fAftaTx TxblASl  cgTax%   cgTax TFedTax spndble`,
		},
		{ // Case 1
			sip: sipJoint,
			expect: `retiree1/retiree2
    age     fIRA    tIRA  TxbleO TxbleSS  deduct   T_inc  earlyP  fedtax  mTaxB% fAftaTx TxblASl  cgTax%   cgTax TFedTax spndble`,
		},
	}
	for i, elem := range tests {
		//fmt.Printf("=============== Case %d =================\n", i)
		ip, err := NewInputParams(elem.sip, nil)
		if err != nil {
			fmt.Printf("TestPrintHeaderTax: %s\n", err)
		}
		csvfile := (*os.File)(nil)
		tablefile := os.Stdout
		ms := ModelSpecs{
			Ip:      *ip,
			Logfile: os.Stdout,
			Errfile: os.Stderr,
			Ao:      NewAppOutput(csvfile, tablefile),
		}

		mychan := make(chan string)
		DoNothing := false //true
		oldout, w, err := ms.RedirectModelSpecsTable(mychan, DoNothing)
		if err != nil {
			t.Errorf("RedirectModelSpecsTable: %s\n", err)
			return // should this be continue?
		}

		ms.printHeaderTax()

		str := ms.RestoreModelSpecsTable(mychan, oldout, w, DoNothing)
		strn := strings.TrimSpace(str)
		if elem.expect != strn {
			str := fmt.Sprintf("TestPrintHeaderTax Case %d:", i)
			showStrMismatch(str, elem.expect, strn)
			t.Errorf("TestPrintHeaderTax case %d:  expected output:\n\t '%s'\n\tbut found:\n\t'%s'\n", i, elem.expect, strn)
		}
	}
}

//def print_tax(res):
func TestPrintTax(t *testing.T) {
	tests := []struct {
		sip    map[string]string
		sxp    *[]float64
		expect string
	}{
		{ // Case 0
			sip: sipSingle,
			sxp: xpSingle,
			expect: `Tax Summary:

retiree1
 age     fIRA    tIRA  TxbleO TxbleSS  deduct   T_inc  earlyP  fedtax  mTaxB% fAftaTx TxblASl  cgTax%   cgTax TFedTax spndble
  65:   40594       0       0       0   13646   26949      0     3431      15       0       0     100       0    3431   37164
  66:   41609       0       0       0   13987   27622      0     3516      15       0       0     100       0    3516   38093
  67:   42650       0       0       0   14337   28313      0     3604      15       0       0     100       0    3604   39045
  68:   43716       0       0       0   14695   29021      0     3694      15       0       0     100       0    3694   40021
  69:   44809       0       0       0   15062   29746      0     3787      15       0       0     100       0    3787   41022
  70:   45929       0       0       0   15439   30490      0     3881      15       0       0     100       0    3881   42048
  71:   47077       0       0       0   15825   31252      0     3978      15       0       0     100       0    3978   43099
  72:   48254       0       0       0   16220   32034      0     4078      15       0       0     100       0    4078   44176
  73:   49460       0       0       0   16626   32834      0     4180      15       0       0     100       0    4180   45281
  74:   50697       0       0       0   17042   33655      0     4284      15       0       0     100       0    4284   46413
  75:   51964       0       0       0   17468   34497      0     4391      15       0       0     100       0    4391   47573
retiree1
 age     fIRA    tIRA  TxbleO TxbleSS  deduct   T_inc  earlyP  fedtax  mTaxB% fAftaTx TxblASl  cgTax%   cgTax TFedTax spndble`,
		},
		{ // Case 1
			sip: sipJoint,
			sxp: xpJoint,
			expect: `Tax Summary:

retiree1/retiree2
    age     fIRA    tIRA  TxbleO TxbleSS  deduct   T_inc  earlyP  fedtax  mTaxB% fAftaTx TxblASl  cgTax%   cgTax TFedTax spndble
 65/ 65:   40594       0       0       0   27291   13303      0     1330      10       0       0     100       0    1330   39264
 66/ 66:   41609       0       0       0   27974   13636      0     1364      10       0       0     100       0    1364   40246
 67/ 67:   42650       0       0       0   28673   13977      0     1398      10       0       0     100       0    1398   41252
 68/ 68:   43716       0       0       0   29390   14326      0     1433      10       0       0     100       0    1433   42283
 69/ 69:   44809       0       0       0   30125   14684      0     1468      10       0       0     100       0    1468   43340
 70/ 70:   45929       0       0       0   30878   15051      0     1505      10       0       0     100       0    1505   44424
 71/ 71:   47077       0       0       0   31650   15427      0     1543      10       0       0     100       0    1543   45534
 72/ 72:   48254       0       0       0   32441   15813      0     1581      10       0       0     100       0    1581   46673
 73/ 73:   49460       0       0       0   33252   16208      0     1621      10       0       0     100       0    1621   47840
 74/ 74:   50697       0       0       0   34083   16614      0     1661      10       0       0     100       0    1661   49036
 75/ 75:   51964       0       0       0   34935   17029      0     1703      10       0       0     100       0    1703   50261
retiree1/retiree2
    age     fIRA    tIRA  TxbleO TxbleSS  deduct   T_inc  earlyP  fedtax  mTaxB% fAftaTx TxblASl  cgTax%   cgTax TFedTax spndble`,
		},
		{ // Case 2
			sip: sipSingle3Acc,
			sxp: xpSingle3Acc,
			expect: `Tax Summary:

retiree1
 age     fIRA    tIRA  TxbleO TxbleSS  deduct   T_inc  earlyP  fedtax  mTaxB% fAftaTx TxblASl  cgTax%   cgTax TFedTax spndble
  65:   54922       0       0       0   13646   41276      0     5580      15       0       0     100       0    5580   49342
  66:   56295       0       0       0   13987   42308      0     5719      15       0       0     100       0    5719   50576
  67:   57702       0       0       0   14337   43366      0     5862      15       0       0     100       0    5862   51840
  68:   59145       0       0       0   14695   44450      0     6009      15       0       0     100       0    6009   53136
  69:   60623       0       0       0   15062   45561      0     6159      15       0       0     100       0    6159   54465
  70:   31531       0       0       0   15439   16093      0     1722      15   24225       0     100       0    1722   55826
  71:   30014       0       0       0   15825   14189      0     1419      10   28627       0     100       0    1419   57222
  72:   30764       0       0       0   16220   14544      0     1454      10   29343       0     100       0    1454   58652
  73:   31533       0       0       0   16626   14907      0     1491      10   30076       0     100       0    1491   60119
  74:   32322       0       0       0   17042   15280      0     1528      10   30828       0     100       0    1528   61622
  75:   33130       0       0       0   17468   15662      0     1566      10       0       0     100       0    1566   63162
retiree1
 age     fIRA    tIRA  TxbleO TxbleSS  deduct   T_inc  earlyP  fedtax  mTaxB% fAftaTx TxblASl  cgTax%   cgTax TFedTax spndble`,
		},
	}
	for i, elem := range tests {
		//fmt.Printf("=============== Case %d =================\n", i)
		ip, err := NewInputParams(elem.sip, nil)
		if err != nil {
			fmt.Printf("TestPrintTax: %s\n", err)
		}
		ti := NewTaxInfo(ip.FilingStatus, 2017)
		taxbins := len(*ti.Taxtable)
		cgbins := len(*ti.Capgainstable)
		vindx, err := NewVectorVarIndex(ip.Numyr, taxbins,
			cgbins, ip.Accmap, os.Stdout)
		if err != nil {
			t.Errorf("TestPrintTax case %d: %s", i, err)
			continue
		}
		logfile := os.Stdout
		csvfile := (*os.File)(nil)
		RoundToOneK := false
		ms, err := NewModelSpecs(vindx, ti, *ip,
			RoundToOneK, false, false,
			os.Stderr, logfile, csvfile, logfile, nil)
		if err != nil {
			t.Errorf("TestPrintTax case %d: %s", i, err)
			continue
		}

		mychan := make(chan string)
		DoNothing := false //true
		oldout, w, err := ms.RedirectModelSpecsTable(mychan, DoNothing)
		if err != nil {
			t.Errorf("RedirectModelSpecsTable: %s\n", err)
			return // should this be continue?
		}

		ms.PrintTax(elem.sxp)

		str := ms.RestoreModelSpecsTable(mychan, oldout, w, DoNothing)
		strn := strings.TrimSpace(str)
		if elem.expect != strn {
			str := fmt.Sprintf("TestPrintTax Case %d:", i)
			showStrMismatch(str, elem.expect, strn)
			t.Errorf("TestPrintTax case %d:  expected output:\n\t '%s'\n\tbut found:\n\t'%s'\n", i, elem.expect, strn)
		}
	}
}

//func (ms ModelSpecs) printHeaderTaxBrackets()
func TestPrintHeaderTaxBrackets(t *testing.T) {
	tests := []struct {
		sip    map[string]string
		expect string
	}{
		{ // Case 0
			sip: sipSingle,
			expect: `Marginal Rate(%):     10     15     25     28     33     35     40
retiree1
 age     fIRA    tIRA  TxbleO TxbleSS  deduct   T_inc  fedtax brckt0 brckt1 brckt2 brckt3 brckt4 brckt5 brckt6 brkTot`,
		},
		{ // Case 1
			sip: sipJoint,
			expect: `Marginal Rate(%):     10     15     25     28     33     35     40
retiree1/retiree2
    age     fIRA    tIRA  TxbleO TxbleSS  deduct   T_inc  fedtax brckt0 brckt1 brckt2 brckt3 brckt4 brckt5 brckt6 brkTot`,
		},
	}
	for i, elem := range tests {
		//fmt.Printf("=============== Case %d =================\n", i)
		ip, err := NewInputParams(elem.sip, nil)
		if err != nil {
			fmt.Printf("TestHeaderTaxBrackets: %s\n", err)
		}
		ti := NewTaxInfo(ip.FilingStatus, 2017)
		csvfile := (*os.File)(nil)
		tablefile := os.Stdout
		ms := ModelSpecs{
			Ip:      *ip,
			Ti:      ti,
			Logfile: os.Stdout,
			Errfile: os.Stderr,
			Ao:      NewAppOutput(csvfile, tablefile),
		}

		mychan := make(chan string)
		DoNothing := false //true
		oldout, w, err := ms.RedirectModelSpecsTable(mychan, DoNothing)
		if err != nil {
			t.Errorf("RedirectModelSpecsTable: %s\n", err)
			return // should this be continue?
		}

		ms.printHeaderTaxBrackets()

		str := ms.RestoreModelSpecsTable(mychan, oldout, w, DoNothing)
		strn := strings.TrimSpace(str)
		if elem.expect != strn {
			str := fmt.Sprintf("TestPrintHeaderTaxBrackets Case %d:", i)
			showStrMismatch(str, elem.expect, strn)
			t.Errorf("TestPrintHeaderTaxBrackests case %d:  expected output:\n\t '%s'\n\tbut found:\n\t'%s'\n", i, elem.expect, strn)
		}
	}
}

//def print_tax_brackets(res):
func TestPrintTaxBrackets(t *testing.T) {
	tests := []struct {
		sip    map[string]string
		sxp    *[]float64
		expect string
	}{
		{ // Case 0
			sip: sipSingle,
			sxp: xpSingle,
			expect: `Overall Tax Bracket Summary:
                                            Marginal Rate(%):     10     15     25     28     33     35     40
retiree1
 age     fIRA    tIRA  TxbleO TxbleSS  deduct   T_inc  fedtax brckt0 brckt1 brckt2 brckt3 brckt4 brckt5 brckt6 brkTot
  65:   40594       0       0       0   13646   26949    3431  12235  14714      0      0      0      0      0  26949
  66:   41609       0       0       0   13987   27622    3516  12541  15081      0      0      0      0      0  27622
  67:   42650       0       0       0   14337   28313    3604  12855  15458      0      0      0      0      0  28313
  68:   43716       0       0       0   14695   29021    3694  13176  15845      0      0      0      0      0  29021
  69:   44809       0       0       0   15062   29746    3787  13505  16241      0      0      0      0      0  29746
  70:   45929       0       0       0   15439   30490    3881  13843  16647      0      0      0      0      0  30490
  71:   47077       0       0       0   15825   31252    3978  14189  17063      0      0      0      0      0  31252
  72:   48254       0       0       0   16220   32034    4078  14544  17490      0      0      0      0      0  32034
  73:   49460       0       0       0   16626   32834    4180  14907  17927      0      0      0      0      0  32834
  74:   50697       0       0       0   17042   33655    4284  15280  18375      0      0      0      0      0  33655
  75:   51964       0       0       0   17468   34497    4391  15662  18835      0      0      0      0      0  34497
                                            Marginal Rate(%):     10     15     25     28     33     35     40
retiree1
 age     fIRA    tIRA  TxbleO TxbleSS  deduct   T_inc  fedtax brckt0 brckt1 brckt2 brckt3 brckt4 brckt5 brckt6 brkTot`,
		},
		{ // Case 1
			sip: sipJoint,
			sxp: xpJoint,
			expect: `Overall Tax Bracket Summary:
                                               Marginal Rate(%):     10     15     25     28     33     35     40
retiree1/retiree2
    age     fIRA    tIRA  TxbleO TxbleSS  deduct   T_inc  fedtax brckt0 brckt1 brckt2 brckt3 brckt4 brckt5 brckt6 brkTot
 65/ 65:   40594       0       0       0   27291   13303    1330  13303      0      0      0      0      0      0  13303
 66/ 66:   41609       0       0       0   27974   13636    1364  13636      0      0      0      0      0      0  13636
 67/ 67:   42650       0       0       0   28673   13977    1398  13977      0      0      0      0      0      0  13977
 68/ 68:   43716       0       0       0   29390   14326    1433  14326      0      0      0      0      0      0  14326
 69/ 69:   44809       0       0       0   30125   14684    1468  14684      0      0      0      0      0      0  14684
 70/ 70:   45929       0       0       0   30878   15051    1505  15051      0      0      0      0      0      0  15051
 71/ 71:   47077       0       0       0   31650   15427    1543  15427      0      0      0      0      0      0  15427
 72/ 72:   48254       0       0       0   32441   15813    1581  15813      0      0      0      0      0      0  15813
 73/ 73:   49460       0       0       0   33252   16208    1621  16208      0      0      0      0      0      0  16208
 74/ 74:   50697       0       0       0   34083   16614    1661  16614      0      0      0      0      0      0  16614
 75/ 75:   51964       0       0       0   34935   17029    1703  17029      0      0      0      0      0      0  17029
                                               Marginal Rate(%):     10     15     25     28     33     35     40
retiree1/retiree2
    age     fIRA    tIRA  TxbleO TxbleSS  deduct   T_inc  fedtax brckt0 brckt1 brckt2 brckt3 brckt4 brckt5 brckt6 brkTot`,
		},
		{ // Case 2
			sip: sipSingle3Acc,
			sxp: xpSingle3Acc,
			expect: `Overall Tax Bracket Summary:
                                            Marginal Rate(%):     10     15     25     28     33     35     40
retiree1
 age     fIRA    tIRA  TxbleO TxbleSS  deduct   T_inc  fedtax brckt0 brckt1 brckt2 brckt3 brckt4 brckt5 brckt6 brkTot
  65:   54922       0       0       0   13646   41276    5580  12235  29041      0      0      0      0      0  41276
  66:   56295       0       0       0   13987   42308    5719  12541  29767      0      0      0      0      0  42308
  67:   57702       0       0       0   14337   43366    5862  12855  30511      0      0      0      0      0  43366
  68:   59145       0       0       0   14695   44450    6009  13176  31274      0      0      0      0      0  44450
  69:   60623       0       0       0   15062   45561    6159  13505  32056      0      0      0      0      0  45561
  70:   31531       0       0       0   15439   16093    1722  13843   2250      0      0      0      0      0  16093
  71:   30014       0       0       0   15825   14189    1419  14189      0      0      0      0      0      0  14189
  72:   30764       0       0       0   16220   14544    1454  14544      0      0      0      0      0      0  14544
  73:   31533       0       0       0   16626   14907    1491  14907      0      0      0      0      0      0  14907
  74:   32322       0       0       0   17042   15280    1528  15280      0      0      0      0      0      0  15280
  75:   33130       0       0       0   17468   15662    1566  15662      0      0      0      0      0      0  15662
                                            Marginal Rate(%):     10     15     25     28     33     35     40
retiree1
 age     fIRA    tIRA  TxbleO TxbleSS  deduct   T_inc  fedtax brckt0 brckt1 brckt2 brckt3 brckt4 brckt5 brckt6 brkTot`,
		},
	}
	for i, elem := range tests {
		//fmt.Printf("=============== Case %d =================\n", i)
		ip, err := NewInputParams(elem.sip, nil)
		if err != nil {
			fmt.Printf("TestPrintTaxBrackets: %s\n", err)
		}
		ti := NewTaxInfo(ip.FilingStatus, 2017)
		taxbins := len(*ti.Taxtable)
		cgbins := len(*ti.Capgainstable)
		vindx, err := NewVectorVarIndex(ip.Numyr, taxbins,
			cgbins, ip.Accmap, os.Stdout)
		if err != nil {
			t.Errorf("TestPrintTaxBrackets case %d: %s", i, err)
			continue
		}
		logfile := os.Stdout
		csvfile := (*os.File)(nil)
		RoundToOneK := false
		ms, err := NewModelSpecs(vindx, ti, *ip,
			RoundToOneK, false, false,
			os.Stderr, logfile, csvfile, logfile, nil)
		if err != nil {
			t.Errorf("TestPrintTaxBrackets case %d: %s", i, err)
			continue
		}

		mychan := make(chan string)
		DoNothing := false //true
		oldout, w, err := ms.RedirectModelSpecsTable(mychan, DoNothing)
		if err != nil {
			t.Errorf("RedirectModelSpecsTable: %s\n", err)
			return // should this be continue?
		}

		ms.PrintTaxBrackets(elem.sxp)

		str := ms.RestoreModelSpecsTable(mychan, oldout, w, DoNothing)
		strn := strings.TrimSpace(str)
		if elem.expect != strn {
			str := fmt.Sprintf("TestPrintTaxBrackets Case %d:", i)
			showStrMismatch(str, elem.expect, strn)
			t.Errorf("TestPrintTaxBrackets case %d:  expected output:\n\t '%s'\n\tbut found:\n\t'%s'\n", i, elem.expect, strn)
		}
	}
}

func TestPrintHeaderCapGainsBrackets(t *testing.T) {
	tests := []struct {
		sip    map[string]string
		expect string
	}{
		{ // Case 0
			sip: sipSingle,
			expect: `Marginal Rate(%):      0     15     20
retiree1
 age  fAftaTx TblASle  cgTax% cgTaxbl   T_inc   cgTax brckt0 brckt1 brckt2 brkTot`,
		},
		{ // Case 1
			sip: sipJoint,
			expect: `Marginal Rate(%):      0     15     20
retiree1/retiree2
    age  fAftaTx TblASle  cgTax% cgTaxbl   T_inc   cgTax brckt0 brckt1 brckt2 brkTot`,
		},
	}
	for i, elem := range tests {
		//fmt.Printf("=============== Case %d =================\n", i)
		ip, err := NewInputParams(elem.sip, nil)
		if err != nil {
			fmt.Printf("TestPrintHeaderCapGainsBrackets: %s\n", err)
		}
		ti := NewTaxInfo(ip.FilingStatus, 2017)
		csvfile := (*os.File)(nil)
		tablefile := os.Stdout
		ms := ModelSpecs{
			Ip:      *ip,
			Ti:      ti,
			Logfile: os.Stdout,
			Errfile: os.Stderr,
			Ao:      NewAppOutput(csvfile, tablefile),
		}

		mychan := make(chan string)
		DoNothing := false //true
		oldout, w, err := ms.RedirectModelSpecsTable(mychan, DoNothing)
		if err != nil {
			t.Errorf("RedirectModelSpecsTable: %s\n", err)
			return // should this be continue?
		}

		ms.printHeaderCapgainsBrackets()

		str := ms.RestoreModelSpecsTable(mychan, oldout, w, DoNothing)
		strn := strings.TrimSpace(str)
		if elem.expect != strn {
			str := fmt.Sprintf("TestPrintHeaderCapgainsBrackets Case %d:", i)
			showStrMismatch(str, elem.expect, strn)
			t.Errorf("TestPrintHeaderCapgainsBrackests case %d:  expected output:\n\t '%s'\n\tbut found:\n\t'%s'\n", i, elem.expect, strn)
		}
	}
}

//def print_cap_gains_brackets(res):
func TestPrintCapGainsBrackets(t *testing.T) {
	tests := []struct {
		sip    map[string]string
		sxp    *[]float64
		expect string
	}{
		{ // Case 0
			sip: sipSingle,
			sxp: xpSingle,
			expect: `Overall Capital Gains Bracket Summary:
                                    Marginal Rate(%):      0     15     20
retiree1
 age  fAftaTx TblASle  cgTax% cgTaxbl   T_inc   cgTax brckt0 brckt1 brckt2 brkTot
  65:       0       0     100       0   26949       0      0      0      0      0
  66:       0       0     100       0   27622       0      0      0      0      0
  67:       0       0     100       0   28313       0      0      0      0      0
  68:       0       0     100       0   29021       0      0      0      0      0
  69:       0       0     100       0   29746       0      0      0      0      0
  70:       0       0     100       0   30490       0      0      0      0      0
  71:       0       0     100       0   31252       0      0      0      0      0
  72:       0       0     100       0   32034       0      0      0      0      0
  73:       0       0     100       0   32834       0      0      0      0      0
  74:       0       0     100       0   33655       0      0      0      0      0
  75:       0       0     100       0   34497       0      0      0      0      0
                                    Marginal Rate(%):      0     15     20
retiree1
 age  fAftaTx TblASle  cgTax% cgTaxbl   T_inc   cgTax brckt0 brckt1 brckt2 brkTot`,
		},
		{ // Case 1
			sip: sipJoint,
			sxp: xpJoint,
			expect: `Overall Capital Gains Bracket Summary:
                                       Marginal Rate(%):      0     15     20
retiree1/retiree2
    age  fAftaTx TblASle  cgTax% cgTaxbl   T_inc   cgTax brckt0 brckt1 brckt2 brkTot
 65/ 65:       0       0     100       0   13303       0      0      0      0      0
 66/ 66:       0       0     100       0   13636       0      0      0      0      0
 67/ 67:       0       0     100       0   13977       0      0      0      0      0
 68/ 68:       0       0     100       0   14326       0      0      0      0      0
 69/ 69:       0       0     100       0   14684       0      0      0      0      0
 70/ 70:       0       0     100       0   15051       0      0      0      0      0
 71/ 71:       0       0     100       0   15427       0      0      0      0      0
 72/ 72:       0       0     100       0   15813       0      0      0      0      0
 73/ 73:       0       0     100       0   16208       0      0      0      0      0
 74/ 74:       0       0     100       0   16614       0      0      0      0      0
 75/ 75:       0       0     100       0   17029       0      0      0      0      0
                                       Marginal Rate(%):      0     15     20
retiree1/retiree2
    age  fAftaTx TblASle  cgTax% cgTaxbl   T_inc   cgTax brckt0 brckt1 brckt2 brkTot`,
		},
		{ // Case 2
			sip: sipSingle3Acc,
			sxp: xpSingle3Acc,
			expect: `Overall Capital Gains Bracket Summary:
                                    Marginal Rate(%):      0     15     20
retiree1
 age  fAftaTx TblASle  cgTax% cgTaxbl   T_inc   cgTax brckt0 brckt1 brckt2 brkTot
  65:       0       0     100       0   41276       0      0      0      0      0
  66:       0       0     100       0   42308       0      0      0      0      0
  67:       0       0     100       0   43366       0      0      0      0      0
  68:       0       0     100       0   44450       0      0      0      0      0
  69:       0       0     100       0   45561       0      0      0      0      0
  70:   24225       0     100   24225   16093       0  24225      0      0  24225
  71:   28627       0     100   28627   14189       0  28627      0      0  28627
  72:   29343       0     100   29343   14544       0  29343      0      0  29343
  73:   30076       0     100   30076   14907       0  30076      0      0  30076
  74:   30828       0     100   30828   15280       0  30828      0      0  30828
  75:       0       0     100       0   15662       0      0      0      0      0
                                    Marginal Rate(%):      0     15     20
retiree1
 age  fAftaTx TblASle  cgTax% cgTaxbl   T_inc   cgTax brckt0 brckt1 brckt2 brkTot`,
		},
	}
	for i, elem := range tests {
		//fmt.Printf("=============== Case %d =================\n", i)
		ip, err := NewInputParams(elem.sip, nil)
		if err != nil {
			fmt.Printf("TestPrintCapGainsBrackets: case %d: %s\n", i, err)
		}
		ti := NewTaxInfo(ip.FilingStatus, 2017)
		taxbins := len(*ti.Taxtable)
		cgbins := len(*ti.Capgainstable)
		vindx, err := NewVectorVarIndex(ip.Numyr, taxbins,
			cgbins, ip.Accmap, os.Stdout)
		if err != nil {
			t.Errorf("TestPrintCapGainsBrackets: case %d: %s\n", i, err)
			continue
		}
		logfile := os.Stdout
		csvfile := (*os.File)(nil)
		RoundToOneK := false
		ms, err := NewModelSpecs(vindx, ti, *ip,
			RoundToOneK, false, false,
			os.Stderr, logfile, csvfile, logfile, nil)
		if err != nil {
			t.Errorf("TestPrintCapGainsBrackets case %d: %s", i, err)
			continue
		}

		mychan := make(chan string)
		DoNothing := false //true
		oldout, w, err := ms.RedirectModelSpecsTable(mychan, DoNothing)
		if err != nil {
			t.Errorf("RedirectModelSpecsTable: %s\n", err)
			return // should this be continue?
		}

		ms.PrintCapGainsBrackets(elem.sxp)

		str := ms.RestoreModelSpecsTable(mychan, oldout, w, DoNothing)
		strn := strings.TrimSpace(str)
		if elem.expect != strn {
			str := fmt.Sprintf("TestPrintCapgainsBrackets Case %d:", i)
			showStrMismatch(str, elem.expect, strn)
			t.Errorf("TestPrintCapGainsBrackets case %d:  expected output:\n\t '%s'\n\tbut found:\n\t'%s'\n", i, elem.expect, strn)
		}
	}
}

func showStrMismatch(title, s1, s2 string) { // TODO move to Utility functions
	for i := 0; i < intMin(len(s1), len(s2)); i++ {
		if s1[i] != s2[i] {
			fmt.Printf("%s:\n", title)
			fmt.Printf("	Strings don't match, miscompare between '[]':\n")
			fmt.Printf("	Char#: %d, CharVals1: %c, CharInts1: %d, CharVals2: %c, CharInts2: %d\n", i, s1[i], s1[i], s2[i], s2[i])
			fmt.Printf("expect:\n'%s'['%s']\n", s1[:i], s1[i:])
			fmt.Printf(" found:\n'%s'['%s']\n", s2[:i], s2[i:])
			break
		}
	}
}

//TODO complete TestOrdinaryTaxable after more print function are working
//func (ms ModelSpecs) ordinaryTaxable(year int, xp *[]float64) float64
func TestOrdinaryTaxable(t *testing.T) {
	tests := []struct {
		sip    map[string]string
		sxp    *[]float64
		year   int
		expect int
	}{
		{
			sip:    sipSingle,
			sxp:    xpSingle,
			year:   7,
			expect: 32033,
		},
	}
	for i, elem := range tests {
		//fmt.Printf("======== CASE %d ========\n", i)
		ip, err := NewInputParams(elem.sip, nil)
		if err != nil {
			t.Errorf("TestOrdinaryTaxable: Case %d: %s\n", i, err)
			continue
		}
		//fmt.Printf("InputParams: %#v\n", ip)
		ti := NewTaxInfo(ip.FilingStatus, 2017)
		taxbins := len(*ti.Taxtable)
		cgbins := len(*ti.Capgainstable)
		vindx, err := NewVectorVarIndex(ip.Numyr, taxbins,
			cgbins, ip.Accmap, os.Stdout)
		if err != nil {
			t.Errorf("TestOrdinaryTaxable case %d: %s", i, err)
			continue
		}
		logfile := os.Stdout
		csvfile := (*os.File)(nil)
		RoundToOneK := false
		ms, err := NewModelSpecs(vindx, ti, *ip,
			RoundToOneK, false, false, os.Stderr, logfile, csvfile, logfile, nil)
		if err != nil {
			t.Errorf("TestOrdinaryTaxable case %d: %s", i, err)
			continue
		}
		ot := ms.ordinaryTaxable(elem.year, elem.sxp)
		if int(ot) != elem.expect {
			t.Errorf("TestOrdinaryTaxable case %d: expected %d, found %d\n", i, elem.expect, int(ot))
		}
	}
}

//func (ms ModelSpecs) IncomeSummary(year int, xp *[]float64) (T, spendable, tax, rate, ncgtax, earlytax float64, rothearly bool)
func TestIncomeSummary(t *testing.T) {
	fmt.Printf("TestIncomeSummary() Not Yet Implemented\n")
}

//func (ms ModelSpecs) getResultTotals(xp *[]float64) (twithd, tcombined, tT, ttax, tcgtax, tearlytax, tspendable, tbeginbal, tendbal float64)

func TestGetResultTotals(t *testing.T) {
	fmt.Printf("TestGetResultTotals() Not Yet Implemented\n")
}

//func (ms ModelSpecs) printBaseConfig(xp *[]float64)  // input is res.x
func TestPrintBaseConfig(t *testing.T) {
	tests := []struct {
		sip    map[string]string
		sxp    *[]float64
		expect string
	}{
		{ // case 0
			sip: sipSingle,
			sxp: xpSingle,
			expect: `======
Optimized for Spending with Single status
	starting at age 65 with an estate of $379_660 liquid and $0 illiquid

No desired minimum or maximum amount specified

After tax yearly income: $37_164 adjusting for inflation
	and final estate at age 76 with $0 liquid and $0 illiquid

total withdrawals: $506_759
total ordinary taxable income $336_414
total ordinary tax on all taxable income: $42_825 (12.7%) of taxable income
total income (withdrawals + other) $506_759
total cap gains tax: $0
total all tax on all income: $42_825 (8.5%)
Total spendable (after tax money): $463_934`,
		},
		{ // case 1
			sip: sipJoint,
			sxp: xpJoint,
			expect: `======
Optimized for Spending with Joint status
	starting at age 65 with an estate of $379_660 liquid and $0 illiquid

No desired minimum or maximum amount specified

After tax yearly income: $39_264 adjusting for inflation
	and final estate at age 76 with $0 liquid and $0 illiquid

total withdrawals: $506_759
total ordinary taxable income $166_068
total ordinary tax on all taxable income: $16_607 (10.0%) of taxable income
total income (withdrawals + other) $506_759
total cap gains tax: $0
total all tax on all income: $16_607 (3.3%)
Total spendable (after tax money): $490_153`,
		},
	}
	for i, elem := range tests {
		//fmt.Printf("======== CASE %d ========\n", i)
		ip, err := NewInputParams(elem.sip, nil)
		if err != nil {
			t.Errorf("TestPrintBaseConfig case %d: %s", i, err)
			continue
		}
		//fmt.Printf("InputParams: %#v\n", ip)
		ti := NewTaxInfo(ip.FilingStatus, 2017)
		taxbins := len(*ti.Taxtable)
		cgbins := len(*ti.Capgainstable)
		vindx, err := NewVectorVarIndex(ip.Numyr, taxbins,
			cgbins, ip.Accmap, os.Stdout)
		if err != nil {
			t.Errorf("TestPrintBaseConfig case %d: %s", i, err)
			continue
		}
		logfile := os.Stdout
		csvfile := (*os.File)(nil)
		RoundToOneK := false
		ms, err := NewModelSpecs(vindx, ti, *ip,
			RoundToOneK, false, false,
			os.Stderr, logfile, csvfile, logfile, nil)
		if err != nil {
			t.Errorf("TestPrintBaseConfig case %d: %s", i, err)
			continue
		}

		mychan := make(chan string)
		DoNothing := false //true
		oldout, w, err := ms.RedirectModelSpecsTable(mychan, DoNothing)
		if err != nil {
			t.Errorf("RedirectModelSpecsTable: %s\n", err)
			return // should this be continue?
		}

		ms.PrintBaseConfig(elem.sxp)

		str := ms.RestoreModelSpecsTable(mychan, oldout, w, DoNothing)
		strn := strings.TrimSpace(str)
		//strn := stripWhitespace(str)
		//ot := ms.ordinaryTaxable(elem.year, elem.sxp)
		if strn != elem.expect {
			str := fmt.Sprintf("TestBaseConfig Case %d:", i)
			showStrMismatch(str, elem.expect, strn)
			t.Errorf("TestPrintBaseConfig case %d: expected\n'%s'\nfound '%s'\n", i, elem.expect, strn)
		}
	}
}

//def verifyInputs( c , A , b ):
func TestVerifyInputs(t *testing.T) {
	fmt.Printf("TestVerifyInputs() Not Yet Implemented\n")
}
func TestAssetByTagAndField(t *testing.T) {
	tests := []struct {
		ip         map[string]string
		verbose    bool
		iRate      float64
		CheckAsset struct {
			Value            float64
			BrokeragePercent float64
		}
	}{
		{ // case 0 // driver generated from AWill.toml
			ip: map[string]string{
				"key1":                         "will",
				"key2":                         "yuli",
				"eT_Age1":                      "56",
				"eT_Age2":                      "54",
				"eT_RetireAge1":                "57",
				"eT_RetireAge2":                "62",
				"eT_PlanThroughAge1":           "100",
				"eT_PlanThroughAge2":           "100",
				"eT_PIA1":                      "31000",
				"eT_PIA2":                      "-1",
				"eT_SS_Start1":                 "70",
				"eT_SS_Start2":                 "68",
				"eT_TDRA1":                     "1400000",
				"eT_TDRA2":                     "18000",
				"eT_TDRA_Contrib1":             "0",
				"eT_Roth1":                     "0",
				"eT_Roth2":                     "0",
				"eT_Aftatax":                   "700000",
				"eT_Aftatax_Basis":             "400000",
				"eT_iRatePercent":              "2.5",
				"eT_rRatePercent":              "6",
				"dollarsInThousands":           "false",
				"eT_Income1":                   "rental_Fessenden",
				"eT_IncomeAmount1":             "36000",
				"eT_IncomeStartAge1":           "57",
				"eT_IncomeEndAge1":             "75",
				"eT_IncomeInflate1":            "true",
				"eT_IncomeTax1":                "true",
				"eT_Expense1":                  "mortgage",
				"eT_ExpenseAmount1":            "37131",
				"eT_ExpenseStartAge1":          "56",
				"eT_ExpenseEndAge1":            "61",
				"eT_Asset1":                    "rental_VanHoutin",
				"eT_AssetValue1":               "700000",
				"eT_AssetAgeToSell1":           "80",
				"eT_AssetCostAndImprovements1": "425000",
				"eT_AssetOwedAtAgeToSell1":     "0",
				"eT_AssetPrimaryResidence1":    "false",
				"eT_AssetRRatePercent1":        "4",
				"eT_Income2":                   "rental_VanHoutin",
				"eT_IncomeAmount2":             "24000",
				"eT_IncomeStartAge2":           "67",
				"eT_IncomeEndAge2":             "80",
				"eT_IncomeInflate2":            "true",
				"eT_IncomeTax2":                "true",
				"eT_Expense2":                  "college",
				"eT_ExpenseAmount2":            "30000",
				"eT_ExpenseStartAge2":          "56",
				"eT_ExpenseEndAge2":            "59",
				"eT_ExpenseInflate2":           "false",
				"eT_Asset2":                    "home",
				"eT_AssetValue2":               "550000",
				"eT_AssetAgeToSell2":           "0",
				"eT_AssetCostAndImprovements2": "300000",
				"eT_AssetOwedAtAgeToSell2":     "0",
				"eT_AssetPrimaryResidence2":    "true",
				"eT_AssetRRatePercent2":        "4",
				"eT_Asset3":                    "rental_Fessenden",
				"eT_AssetValue3":               "900000",
				"eT_AssetAgeToSell3":           "75",
				"eT_AssetCostAndImprovements3": "450000",
				"eT_AssetOwedAtAgeToSell3":     "0",
				"eT_AssetPrimaryResidence3":    "false",
				"eT_AssetRRatePercent3":        "4",
				//Added
				"filingStatus": "joint",
				"eT_maximize":  "Spending", // or "PlusEstate"
			},
			verbose: true,
			iRate:   1.025,
		},
		{ // Case 1  // case to match AWill.toml Hand coded
			ip: map[string]string{
				"setName":                    "AWill.toml",
				"filingStatus":               "joint",
				"key1":                       "will",
				"key2":                       "yuli",
				"eT_Age1":                    "56",
				"eT_Age2":                    "54",
				"eT_RetireAge1":              "57",
				"eT_RetireAge2":              "62",
				"eT_PlanThroughAge1":         "100",
				"eT_PlanThroughAge2":         "100",
				"eT_PIA1":                    "31000", //31K
				"eT_PIA2":                    "-1",
				"eT_SS_Start1":               "70",
				"eT_SS_Start2":               "68",
				"eT_TDRA1":                   "1400000", //1.4M
				"eT_TDRA2":                   "18000",   //18K
				"eT_TDRA_Rate1":              "",
				"eT_TDRA_Rate2":              "",
				"eT_TDRA_Contrib1":           "",
				"eT_TDRA_Contrib2":           "", // contribute 5k per year
				"eT_TDRA_ContribStartAge1":   "",
				"eT_TDRA_ContribStartAge2":   "",
				"eT_TDRA_ContribEndAge1":     "",
				"eT_TDRA_ContribEndAge2":     "",
				"eT_Roth1":                   "0",
				"eT_Roth2":                   "0",
				"eT_Roth_Rate1":              "",
				"eT_Roth_Rate2":              "",
				"eT_Roth_Contrib1":           "",
				"eT_Roth_Contrib2":           "",
				"eT_Roth_ContribStartAge1":   "",
				"eT_Roth_ContribStartAge2":   "",
				"eT_Roth_ContribEndAge1":     "",
				"eT_Roth_ContribEndAge2":     "",
				"eT_Aftatax":                 "700000", //700k
				"eT_Aftatax_Basis":           "400000", //400k
				"eT_Aftatax_Rate":            "",
				"eT_Aftatax_Contrib":         "",
				"eT_Aftatax_ContribStartAge": "",
				"eT_Aftatax_ContribEndAge":   "",

				"eT_iRatePercent":    "2.5",
				"eT_rRatePercent":    "6",
				"eT_maximize":        "Spending", // or "PlusEstate"
				"dollarsInThousands": "false",

				"eT_Income1":         "rental_fessinden",
				"eT_IncomeAmount1":   "36000",
				"eT_IncomeStartAge1": "57",
				"eT_IncomeEndAge1":   "75",
				"eT_IncomeInflate1":  "true",
				"eT_IncomeTax1":      "true",

				//prototype entries below
				"eT_Income2":         "rental_Van_Houten",
				"eT_IncomeAmount2":   "24000",
				"eT_IncomeStartAge2": "67",
				"eT_IncomeEndAge2":   "80",
				"eT_IncomeInflate2":  "true",
				"eT_IncomeTax2":      "true",

				//prototype entries below
				"eT_Expense1":         "morgage",
				"eT_ExpenseAmount1":   "37131",
				"eT_ExpenseStartAge1": "56",
				"eT_ExpenseEndAge1":   "61",
				"eT_ExpenseInflate1":  "",
				"eT_ExpenseTax1":      "", //ignored, or should be

				//prototype entries below
				"eT_Expense2":         "college",
				"eT_ExpenseAmount2":   "30000",
				"eT_ExpenseStartAge2": "56",
				"eT_ExpenseEndAge2":   "59",
				"eT_ExpenseInflate2":  "false",
				"eT_ExpenseTax2":      "", //ignored, or should be

				//prototype entries below
				"eT_Asset1":                    "rental_fessenden",
				"eT_AssetValue1":               "900000",
				"eT_AssetAgeToSell1":           "75",
				"eT_AssetCostAndImprovements1": "450000",
				"eT_AssetOwedAtAgeToSell1":     "0",
				"eT_AssetPrimaryResidence1":    "false",
				"eT_AssetRRatePercent1":        "4",
				"eT_AssetBrokeragePercent1":    "",

				"eT_Asset2":                    "rental_van_houten",
				"eT_AssetValue2":               "700000",
				"eT_AssetAgeToSell2":           "80",
				"eT_AssetCostAndImprovements2": "425000",
				"eT_AssetOwedAtAgeToSell2":     "0",
				"eT_AssetPrimaryResidence2":    "false",
				"eT_AssetRRatePercent2":        "4", // python defaults to global rate
				"eT_AssetBrokeragePercent2":    "",

				"eT_Asset3":                    "home",
				"eT_AssetValue3":               "550000",
				"eT_AssetAgeToSell3":           "0",
				"eT_AssetCostAndImprovements3": "300000",
				"eT_AssetOwedAtAgeToSell3":     "0",
				"eT_AssetPrimaryResidence3":    "true",
				"eT_AssetRRatePercent3":        "4", // python defaults to global rate
				"eT_AssetBrokeragePercent3":    "",
			},
			verbose: true,
			iRate:   1.025,
		},
		{ // Case 2 // case to match mobile.toml
			ip: map[string]string{
				"setName":                    "activeParams",
				"filingStatus":               "joint",
				"key1":                       "retiree1",
				"key2":                       "retiree2",
				"eT_Age1":                    "54",
				"eT_Age2":                    "54",
				"eT_RetireAge1":              "65",
				"eT_RetireAge2":              "65",
				"eT_PlanThroughAge1":         "75",
				"eT_PlanThroughAge2":         "75",
				"eT_PIA1":                    "20", //20K
				"eT_PIA2":                    "-1",
				"eT_SS_Start1":               "70",
				"eT_SS_Start2":               "70",
				"eT_TDRA1":                   "200", // 200k
				"eT_TDRA2":                   "",
				"eT_TDRA_Rate1":              "",
				"eT_TDRA_Rate2":              "",
				"eT_TDRA_Contrib1":           "",
				"eT_TDRA_Contrib2":           "5", // contribute 5k per year
				"eT_TDRA_ContribStartAge1":   "",
				"eT_TDRA_ContribStartAge2":   "63",
				"eT_TDRA_ContribEndAge1":     "",
				"eT_TDRA_ContribEndAge2":     "64",
				"eT_Roth1":                   "5",
				"eT_Roth2":                   "20", //20k
				"eT_Roth_Rate1":              "",
				"eT_Roth_Rate2":              "",
				"eT_Roth_Contrib1":           "",
				"eT_Roth_Contrib2":           "",
				"eT_Roth_ContribStartAge1":   "",
				"eT_Roth_ContribStartAge2":   "",
				"eT_Roth_ContribEndAge1":     "",
				"eT_Roth_ContribEndAge2":     "",
				"eT_Aftatax":                 "60", //60k
				"eT_Aftatax_Rate":            "",
				"eT_Aftatax_Contrib":         "10", //10K
				"eT_Aftatax_ContribStartAge": "63",
				"eT_Aftatax_ContribEndAge":   "67",

				"eT_iRatePercent": "2.5",
				"eT_rRatePercent": "6",
				"eT_maximize":     "Spending", // or "PlusEstate"

				//prototype entries below
				"eT_Income1":         "rental1",
				"eT_IncomeAmount1":   "1",
				"eT_IncomeStartAge1": "63",
				"eT_IncomeEndAge1":   "67",
				"eT_IncomeInflate1":  "true",
				"eT_IncomeTax1":      "true",

				//prototype entries below
				"eT_Income2":         "rental2",
				"eT_IncomeAmount2":   "2",
				"eT_IncomeStartAge2": "62",
				"eT_IncomeEndAge2":   "70",
				"eT_IncomeInflate2":  "false",
				"eT_IncomeTax2":      "true",

				//prototype entries below
				"eT_Expense1":         "exp1",
				"eT_ExpenseAmount1":   "1",
				"eT_ExpenseStartAge1": "63",
				"eT_ExpenseEndAge1":   "67",
				"eT_ExpenseInflate1":  "true",
				"eT_ExpenseTax1":      "true", //ignored, or should be

				//prototype entries below
				"eT_Expense2":         "exp2",
				"eT_ExpenseAmount2":   "2",
				"eT_ExpenseStartAge2": "62",
				"eT_ExpenseEndAge2":   "70",
				"eT_ExpenseInflate2":  "false",
				"eT_ExpenseTax2":      "true", //ignored, or should be

				//prototype entries below
				"eT_Asset1":                    "ass1",
				"eT_AssetValue1":               "100",
				"eT_AssetAgeToSell1":           "73",
				"eT_AssetCostAndImprovements1": "20",
				"eT_AssetOwedAtAgeToSell1":     "10",
				"eT_AssetPrimaryResidence1":    "True",
				"eT_AssetRRatePercent1":        "4",
				"eT_AssetBrokeragePercent1":    "4",

				//prototype entries below
				"eT_Asset2":                    "ass2",
				"eT_AssetValue2":               "100",
				"eT_AssetAgeToSell2":           "73",
				"eT_AssetCostAndImprovements2": "20",
				"eT_AssetOwedAtAgeToSell2":     "10",
				"eT_AssetPrimaryResidence2":    "false",
				"eT_AssetRRatePercent2":        "6", // python defaults to global rate
				"eT_AssetBrokeragePercent2":    "",
			},
			verbose: true,
			iRate:   1.025,
		},
		{ // Case 3 // case to match mobile.toml
			ip: map[string]string{
				"setName":                    "activeParams",
				"filingStatus":               "single",
				"key1":                       "retiree1",
				"key2":                       "",
				"eT_Age1":                    "54",
				"eT_Age2":                    "",
				"eT_RetireAge1":              "65",
				"eT_RetireAge2":              "",
				"eT_PlanThroughAge1":         "75",
				"eT_PlanThroughAge2":         "",
				"eT_PIA1":                    "",
				"eT_PIA2":                    "",
				"eT_SS_Start1":               "",
				"eT_SS_Start2":               "",
				"eT_TDRA1":                   "200", // 200k
				"eT_TDRA2":                   "",
				"eT_TDRA_Rate1":              "",
				"eT_TDRA_Rate2":              "",
				"eT_TDRA_Contrib1":           "",
				"eT_TDRA_Contrib2":           "",
				"eT_TDRA_ContribStartAge1":   "",
				"eT_TDRA_ContribStartAge2":   "",
				"eT_TDRA_ContribEndAge1":     "",
				"eT_TDRA_ContribEndAge2":     "",
				"eT_Roth1":                   "",
				"eT_Roth2":                   "",
				"eT_Roth_Rate1":              "",
				"eT_Roth_Rate2":              "",
				"eT_Roth_Contrib1":           "",
				"eT_Roth_Contrib2":           "",
				"eT_Roth_ContribStartAge1":   "",
				"eT_Roth_ContribStartAge2":   "",
				"eT_Roth_ContribEndAge1":     "",
				"eT_Roth_ContribEndAge2":     "",
				"eT_Aftatax":                 "",
				"eT_Aftatax_Rate":            "",
				"eT_Aftatax_Contrib":         "",
				"eT_Aftatax_ContribStartAge": "",
				"eT_Aftatax_ContribEndAge":   "",

				"eT_iRatePercent": "2.5",
				"eT_rRatePercent": "6",
				"eT_maximize":     "Spending", // or "PlusEstate"
			},
			verbose: true,
			iRate:   1.025,
		},
	}
	if testing.Short() { //Skip if set "-short"
		t.Skip("TestAssetByTagAndField() (full runs): skipping when set '-short'")
	}
	for i, elem := range tests {
		if i != 1 {
			continue
		}
		fmt.Printf("======== CASE %d ========\n", i)
		ip, err := NewInputParams(elem.ip, nil)
		if err != nil {
			t.Errorf("TestResultsOutput case %d: %s", i, err)
			continue
		}
		//fmt.Printf("InputParams: %#v\n", ip)
		ti := NewTaxInfo(ip.FilingStatus, 2017)
		taxbins := len(*ti.Taxtable)
		cgbins := len(*ti.Capgainstable)
		vindx, err := NewVectorVarIndex(ip.Numyr, taxbins,
			cgbins, ip.Accmap, os.Stdout)
		if err != nil {
			t.Errorf("TestResultsOutput case %d: %s", i, err)
			continue
		}
		logfile, err := os.Create("ModelMatixPP.log")
		if err != nil {
			t.Errorf("TestResultsOutput case %d: %s", i, err)
			continue
		}
		//csvfile := (*os.File)(nil)
		csvfile, err := os.Create("ModelOutput.csv")
		if err != nil {
			t.Errorf("TestResultsOutput case %d: %s", i, err)
			continue
		}
		RoundToOneK := false
		ms, err := NewModelSpecs(vindx, ti, *ip,
			RoundToOneK, false, false,
			os.Stderr, logfile, csvfile, logfile, nil)
		if err != nil {
			t.Errorf("TestResultsOutput case %d: %s", i, err)
			continue
		}
		//fmt.Printf("ModelSpecs: %#v\n", ms)

		c, a, b, _ /*notes*/ := ms.BuildModel()

		tol := 1.0e-7
		bland := false
		maxiter := 4000
		callback := lpsimplex.Callbackfunc(nil)
		//callback := lpsimplex.LPSimplexVerboseCallback
		//callback := lpsimplex.LPSimplexTerseCallback
		disp := true // false //true
		start := time.Now()
		res := lpsimplex.LPSimplex(c, a, b, nil, nil, nil, callback, disp, maxiter, tol, bland)
		elapsed := time.Since(start)

		/*
			err = BinDumpModel(c, a, b, res.X, "./RPlanModelgo.datX")
			if err != nil {
				t.Errorf("TestResultsOutput case %d: %s", i, err)
				continue
			}
			BinCheckModelFiles("./RPlanModelgo.datX", "./RPlanModelpython.datX", &vindx)
		*/

		//fmt.Printf("Res: %#v\n", res)
		str := fmt.Sprintf("Message: %v\n", res.Message)
		fmt.Printf(str)
		str = fmt.Sprintf("Time: LPSimplex() took %s\n", elapsed)
		fmt.Printf(str)
		fmt.Printf("Called LPSimplex() for m:%d x n:%d model\n", len(a), len(a[0]))

		if res.Success {

			for i, asset := range ms.Ip.Assets {
				ipTag := asset.Tag
				fmt.Printf(" ====== Asset %d %s ======== \n", i, ipTag)
				rvalue := ms.AssetByTagAndField(ipTag, "Value")
				if asset.Value != int(rvalue) {
					t.Errorf("TestAssetByTagAndField case %d: Asset %s Value Expect: %d, found: %d\n", i, ipTag, asset.Value, int(rvalue))
				}
				rvalue = ms.AssetByTagAndField(ipTag, "CostAndImprovements")
				if asset.CostAndImprovements != int(rvalue) {
					t.Errorf("TestAssetByTagAndField case %d: Asset %s CostAndImprovements Expect: %d, found: %d\n", i, ipTag, asset.CostAndImprovements, int(rvalue))
				}
				rvalue = ms.AssetByTagAndField(ipTag, "AgeToSell")
				if asset.AgeToSell != int(rvalue) {
					t.Errorf("TestAssetByTagAndField case %d: Asset %s AgeToSell Expect: %d, found: %d\n", i, ipTag, asset.AgeToSell, int(rvalue))
				}
				rvalue = ms.AssetByTagAndField(ipTag, "OwedAtAgeToSell")
				if asset.OwedAtAgeToSell != int(rvalue) {
					t.Errorf("TestAssetByTagAndField case %d: Asset %s OwedAtAgeToSell Expect: %d, found: %d\n", i, ipTag, asset.OwedAtAgeToSell, int(rvalue))
				}
				rvalue = ms.AssetByTagAndField(ipTag, "PrimaryResidence")
				if asset.PrimaryResidence != (rvalue == 1.0) {
					t.Errorf("TestAssetByTagAndField case %d: Asset %s PrimaryResidence Expect: %v, found: %v\n", i, ipTag, asset.PrimaryResidence, rvalue == 1.0)
				}
				rvalue = ms.AssetByTagAndField(ipTag, "AssetRRate")
				if asset.AssetRRate != rvalue {
					t.Errorf("TestAssetByTagAndField case %d: Asset %s AssetRRate Expect: %6.2f, found: %6.2f\n", i, ipTag, asset.AssetRRate, rvalue)
				}
				rvalue = ms.AssetByTagAndField(ipTag, "BrokeragePercent")
				if asset.BrokeragePercent != rvalue {
					t.Errorf("TestAssetByTagAndField case %d: Asset %s BrokeragePercent Expect: %6.2f, found: %6.2f\n", i, ipTag, asset.BrokeragePercent, rvalue)
				}
				//if ms.Ip.PlanStart <= asset.AgeToSell && ms.Ip.PlanEnd >= asset.AgeToSell {
				value, brate, assetRR, basis, owed, prime, age := ms.AssetByTag(ipTag)
				ageToSell := int(age)
				year := ageToSell - ms.Ip.StartPlan
				if year < 0 {
					year = 0
				}
				price := value * math.Pow(assetRR, float64(ageToSell-ms.Ip.Age1))
				bfee := price * brate
				net := price*(1-brate) - owed
				if net < 0.0 || ageToSell < ms.Ip.StartPlan || ageToSell > ms.Ip.EndPlan {
					net = 0.0
				}
				exclude := 0.0
				taxable := price*(1-brate) - basis
				if prime == 1.0 {
					exclude = ms.Ti.Primeresidence * math.Pow(ms.Ip.IRate, float64(ageToSell-ms.Ip.Age1))
					if exclude > taxable {
						exclude = taxable
					}
					taxable -= exclude
				}
				if taxable < 0.0 {
					taxable = 0.0
				}
				rvalue = ms.AssetByTagAndField(ipTag, "SellNet")
				if net != rvalue {
					t.Errorf("TestAssetByTagAndField case %d: Asset %s SellNet Expect: %6.2f, found: %6.2f\n", i, ipTag, net, rvalue)
				}
				rvalue = ms.AssetByTagAndField(ipTag, "SellPrice")
				if price != rvalue {
					t.Errorf("TestAssetByTagAndField case %d: Asset %s SellPrice Expect: %6.2f, found: %6.2f\n", i, ipTag, price, rvalue)
				}
				rvalue = ms.AssetByTagAndField(ipTag, "BroderFee")
				if bfee != rvalue {
					t.Errorf("TestAssetByTagAndField case %d: Asset %s BroderFee Expect: %6.2f, found: %6.2f\n", i, ipTag, bfee, rvalue)
				}
				rvalue = ms.AssetByTagAndField(ipTag, "Taxable")
				if taxable != rvalue {
					t.Errorf("TestAssetByTagAndField case %d: Asset %s Taxable Expect: %6.2f, found: %6.2f\n", i, ipTag, taxable, rvalue)
				}
				rvalue = ms.AssetByTagAndField(ipTag, "MaxTaxableExclution")
				if exclude != rvalue {
					t.Errorf("TestAssetByTagAndField case %d: Asset %s MaxTaxableExclution Expect: %6.2f, found: %6.2f\n", i, ipTag, exclude, rvalue)
				}
				//}

				//fmt.Printf("Asset %s Value: %6.2f\n", ms.Assettags[1], ms.AssetByTagAndField(ms.Assettags[1], "Value"))
			}

			NoOutput := false // true // will want to turn this off later
			if !NoOutput {
				ms.ConsistencyCheckBrackets(&res.X)
				ms.ConsistencyCheckSpendable(&res.X)

				ms.PrintActivitySummary(&res.X)
				ms.PrintIncomeExpenseDetails()
				ms.PrintAccountTrans(&res.X)
				ms.PrintTax(&res.X)
				ms.PrintTaxBrackets(&res.X)
				ms.PrintCapGainsBrackets(&res.X)
				/*
					//ms.print_model_results(res.x)
						        if args.verboseincome:
						            print_income_expense_details()
						        if args.verboseaccounttrans:
						            print_account_trans(res)
						        if args.verbosetax:
						            print_tax(res)
						        if args.verbosetaxbrackets:
						            print_tax_brackets(res)
									print_cap_gains_brackets(res)
				*/
				ms.PrintBaseConfig(&res.X)
			}
		}
		//createDefX(&res.X)
	}
}

func TestResultsOutput(t *testing.T) {
	tests := []struct {
		ip      map[string]string
		verbose bool
		iRate   float64
	}{
		{ // case 0 // driver generated from AWill.toml
			ip: map[string]string{
				"key1":                         "will",
				"key2":                         "yuli",
				"eT_Age1":                      "56",
				"eT_Age2":                      "54",
				"eT_RetireAge1":                "57",
				"eT_RetireAge2":                "62",
				"eT_PlanThroughAge1":           "100",
				"eT_PlanThroughAge2":           "100",
				"eT_PIA1":                      "31000",
				"eT_PIA2":                      "-1",
				"eT_SS_Start1":                 "70",
				"eT_SS_Start2":                 "68",
				"eT_TDRA1":                     "1400000",
				"eT_TDRA2":                     "18000",
				"eT_TDRA_Contrib1":             "0",
				"eT_Roth1":                     "0",
				"eT_Roth2":                     "0",
				"eT_Aftatax":                   "700000",
				"eT_Aftatax_Basis":             "400000",
				"eT_iRatePercent":              "2.5",
				"eT_rRatePercent":              "6",
				"dollarsInThousands":           "false",
				"eT_Income1":                   "rental_Fessenden",
				"eT_IncomeAmount1":             "36000",
				"eT_IncomeStartAge1":           "57",
				"eT_IncomeEndAge1":             "75",
				"eT_IncomeInflate1":            "true",
				"eT_IncomeTax1":                "true",
				"eT_Expense1":                  "mortgage",
				"eT_ExpenseAmount1":            "37131",
				"eT_ExpenseStartAge1":          "56",
				"eT_ExpenseEndAge1":            "61",
				"eT_Asset1":                    "rental_VanHoutin",
				"eT_AssetValue1":               "700000",
				"eT_AssetAgeToSell1":           "80",
				"eT_AssetCostAndImprovements1": "425000",
				"eT_AssetOwedAtAgeToSell1":     "0",
				"eT_AssetPrimaryResidence1":    "false",
				"eT_AssetRRatePercent1":        "4",
				"eT_Income2":                   "rental_VanHoutin",
				"eT_IncomeAmount2":             "24000",
				"eT_IncomeStartAge2":           "67",
				"eT_IncomeEndAge2":             "80",
				"eT_IncomeInflate2":            "true",
				"eT_IncomeTax2":                "true",
				"eT_Expense2":                  "college",
				"eT_ExpenseAmount2":            "30000",
				"eT_ExpenseStartAge2":          "56",
				"eT_ExpenseEndAge2":            "59",
				"eT_ExpenseInflate2":           "false",
				"eT_Asset2":                    "home",
				"eT_AssetValue2":               "550000",
				"eT_AssetAgeToSell2":           "0",
				"eT_AssetCostAndImprovements2": "300000",
				"eT_AssetOwedAtAgeToSell2":     "0",
				"eT_AssetPrimaryResidence2":    "true",
				"eT_AssetRRatePercent2":        "4",
				"eT_Asset3":                    "rental_Fessenden",
				"eT_AssetValue3":               "900000",
				"eT_AssetAgeToSell3":           "75",
				"eT_AssetCostAndImprovements3": "450000",
				"eT_AssetOwedAtAgeToSell3":     "0",
				"eT_AssetPrimaryResidence3":    "false",
				"eT_AssetRRatePercent3":        "4",
				//Added
				"filingStatus": "joint",
				"eT_maximize":  "Spending", // or "PlusEstate"
			},
			verbose: true,
			iRate:   1.025,
		},
		{ // Case 1  // case to match AWill.toml Hand coded
			ip: map[string]string{
				"setName":                    "AWill.toml",
				"filingStatus":               "joint",
				"key1":                       "will",
				"key2":                       "yuli",
				"eT_Age1":                    "56",
				"eT_Age2":                    "54",
				"eT_RetireAge1":              "57",
				"eT_RetireAge2":              "62",
				"eT_PlanThroughAge1":         "100",
				"eT_PlanThroughAge2":         "100",
				"eT_PIA1":                    "31000", //31K
				"eT_PIA2":                    "-1",
				"eT_SS_Start1":               "70",
				"eT_SS_Start2":               "68",
				"eT_TDRA1":                   "1400000", //1.4M
				"eT_TDRA2":                   "18000",   //18K
				"eT_TDRA_Rate1":              "",
				"eT_TDRA_Rate2":              "",
				"eT_TDRA_Contrib1":           "",
				"eT_TDRA_Contrib2":           "", // contribute 5k per year
				"eT_TDRA_ContribStartAge1":   "",
				"eT_TDRA_ContribStartAge2":   "",
				"eT_TDRA_ContribEndAge1":     "",
				"eT_TDRA_ContribEndAge2":     "",
				"eT_Roth1":                   "0",
				"eT_Roth2":                   "0",
				"eT_Roth_Rate1":              "",
				"eT_Roth_Rate2":              "",
				"eT_Roth_Contrib1":           "",
				"eT_Roth_Contrib2":           "",
				"eT_Roth_ContribStartAge1":   "",
				"eT_Roth_ContribStartAge2":   "",
				"eT_Roth_ContribEndAge1":     "",
				"eT_Roth_ContribEndAge2":     "",
				"eT_Aftatax":                 "700000", //700k
				"eT_Aftatax_Basis":           "400000", //400k
				"eT_Aftatax_Rate":            "",
				"eT_Aftatax_Contrib":         "",
				"eT_Aftatax_ContribStartAge": "",
				"eT_Aftatax_ContribEndAge":   "",

				"eT_iRatePercent":    "2.5",
				"eT_rRatePercent":    "6",
				"eT_maximize":        "Spending", // or "PlusEstate"
				"dollarsInThousands": "false",

				"eT_Income1":         "rental_fessinden",
				"eT_IncomeAmount1":   "36000",
				"eT_IncomeStartAge1": "57",
				"eT_IncomeEndAge1":   "75",
				"eT_IncomeInflate1":  "true",
				"eT_IncomeTax1":      "true",

				//prototype entries below
				"eT_Income2":         "rental_Van_Houten",
				"eT_IncomeAmount2":   "24000",
				"eT_IncomeStartAge2": "67",
				"eT_IncomeEndAge2":   "80",
				"eT_IncomeInflate2":  "true",
				"eT_IncomeTax2":      "true",

				//prototype entries below
				"eT_Expense1":         "morgage",
				"eT_ExpenseAmount1":   "37131",
				"eT_ExpenseStartAge1": "56",
				"eT_ExpenseEndAge1":   "61",
				"eT_ExpenseInflate1":  "",
				"eT_ExpenseTax1":      "", //ignored, or should be

				//prototype entries below
				"eT_Expense2":         "college",
				"eT_ExpenseAmount2":   "30000",
				"eT_ExpenseStartAge2": "56",
				"eT_ExpenseEndAge2":   "59",
				"eT_ExpenseInflate2":  "false",
				"eT_ExpenseTax2":      "", //ignored, or should be

				//prototype entries below
				"eT_Asset1":                    "rental_fessenden",
				"eT_AssetValue1":               "900000",
				"eT_AssetAgeToSell1":           "75",
				"eT_AssetCostAndImprovements1": "450000",
				"eT_AssetOwedAtAgeToSell1":     "0",
				"eT_AssetPrimaryResidence1":    "false",
				"eT_AssetRRatePercent1":        "4",
				"eT_AssetBrokeragePercent1":    "",

				"eT_Asset2":                    "rental_van_houten",
				"eT_AssetValue2":               "700000",
				"eT_AssetAgeToSell2":           "80",
				"eT_AssetCostAndImprovements2": "425000",
				"eT_AssetOwedAtAgeToSell2":     "0",
				"eT_AssetPrimaryResidence2":    "false",
				"eT_AssetRRatePercent2":        "4", // python defaults to global rate
				"eT_AssetBrokeragePercent2":    "",

				"eT_Asset3":                    "home",
				"eT_AssetValue3":               "550000",
				"eT_AssetAgeToSell3":           "0",
				"eT_AssetCostAndImprovements3": "300000",
				"eT_AssetOwedAtAgeToSell3":     "0",
				"eT_AssetPrimaryResidence3":    "true",
				"eT_AssetRRatePercent3":        "4", // python defaults to global rate
				"eT_AssetBrokeragePercent3":    "",
			},
			verbose: true,
			iRate:   1.025,
		},
		{ // Case 2 // case to match mobile.toml
			ip: map[string]string{
				"setName":                    "activeParams",
				"filingStatus":               "joint",
				"key1":                       "retiree1",
				"key2":                       "retiree2",
				"eT_Age1":                    "54",
				"eT_Age2":                    "54",
				"eT_RetireAge1":              "65",
				"eT_RetireAge2":              "65",
				"eT_PlanThroughAge1":         "75",
				"eT_PlanThroughAge2":         "75",
				"eT_PIA1":                    "20", //20K
				"eT_PIA2":                    "-1",
				"eT_SS_Start1":               "70",
				"eT_SS_Start2":               "70",
				"eT_TDRA1":                   "200", // 200k
				"eT_TDRA2":                   "",
				"eT_TDRA_Rate1":              "",
				"eT_TDRA_Rate2":              "",
				"eT_TDRA_Contrib1":           "",
				"eT_TDRA_Contrib2":           "5", // contribute 5k per year
				"eT_TDRA_ContribStartAge1":   "",
				"eT_TDRA_ContribStartAge2":   "63",
				"eT_TDRA_ContribEndAge1":     "",
				"eT_TDRA_ContribEndAge2":     "64",
				"eT_Roth1":                   "5",
				"eT_Roth2":                   "20", //20k
				"eT_Roth_Rate1":              "",
				"eT_Roth_Rate2":              "",
				"eT_Roth_Contrib1":           "",
				"eT_Roth_Contrib2":           "",
				"eT_Roth_ContribStartAge1":   "",
				"eT_Roth_ContribStartAge2":   "",
				"eT_Roth_ContribEndAge1":     "",
				"eT_Roth_ContribEndAge2":     "",
				"eT_Aftatax":                 "60", //60k
				"eT_Aftatax_Rate":            "",
				"eT_Aftatax_Contrib":         "10", //10K
				"eT_Aftatax_ContribStartAge": "63",
				"eT_Aftatax_ContribEndAge":   "67",

				"eT_iRatePercent": "2.5",
				"eT_rRatePercent": "6",
				"eT_maximize":     "Spending", // or "PlusEstate"

				//prototype entries below
				"eT_Income1":         "rental1",
				"eT_IncomeAmount1":   "1",
				"eT_IncomeStartAge1": "63",
				"eT_IncomeEndAge1":   "67",
				"eT_IncomeInflate1":  "true",
				"eT_IncomeTax1":      "true",

				//prototype entries below
				"eT_Income2":         "rental2",
				"eT_IncomeAmount2":   "2",
				"eT_IncomeStartAge2": "62",
				"eT_IncomeEndAge2":   "70",
				"eT_IncomeInflate2":  "false",
				"eT_IncomeTax2":      "true",

				//prototype entries below
				"eT_Expense1":         "exp1",
				"eT_ExpenseAmount1":   "1",
				"eT_ExpenseStartAge1": "63",
				"eT_ExpenseEndAge1":   "67",
				"eT_ExpenseInflate1":  "true",
				"eT_ExpenseTax1":      "true", //ignored, or should be

				//prototype entries below
				"eT_Expense2":         "exp2",
				"eT_ExpenseAmount2":   "2",
				"eT_ExpenseStartAge2": "62",
				"eT_ExpenseEndAge2":   "70",
				"eT_ExpenseInflate2":  "false",
				"eT_ExpenseTax2":      "true", //ignored, or should be

				//prototype entries below
				"eT_Asset1":                    "ass1",
				"eT_AssetValue1":               "100",
				"eT_AssetAgeToSell1":           "73",
				"eT_AssetCostAndImprovements1": "20",
				"eT_AssetOwedAtAgeToSell1":     "10",
				"eT_AssetPrimaryResidence1":    "True",
				"eT_AssetRRatePercent1":        "4",
				"eT_AssetBrokeragePercent1":    "4",

				//prototype entries below
				"eT_Asset2":                    "ass2",
				"eT_AssetValue2":               "100",
				"eT_AssetAgeToSell2":           "73",
				"eT_AssetCostAndImprovements2": "20",
				"eT_AssetOwedAtAgeToSell2":     "10",
				"eT_AssetPrimaryResidence2":    "false",
				"eT_AssetRRatePercent2":        "6", // python defaults to global rate
				"eT_AssetBrokeragePercent2":    "",
			},
			verbose: true,
			iRate:   1.025,
		},
		{ // Case 3 // case to match mobile.toml
			ip: map[string]string{
				"setName":                    "activeParams",
				"filingStatus":               "single",
				"key1":                       "retiree1",
				"key2":                       "",
				"eT_Age1":                    "54",
				"eT_Age2":                    "",
				"eT_RetireAge1":              "65",
				"eT_RetireAge2":              "",
				"eT_PlanThroughAge1":         "75",
				"eT_PlanThroughAge2":         "",
				"eT_PIA1":                    "",
				"eT_PIA2":                    "",
				"eT_SS_Start1":               "",
				"eT_SS_Start2":               "",
				"eT_TDRA1":                   "200", // 200k
				"eT_TDRA2":                   "",
				"eT_TDRA_Rate1":              "",
				"eT_TDRA_Rate2":              "",
				"eT_TDRA_Contrib1":           "",
				"eT_TDRA_Contrib2":           "",
				"eT_TDRA_ContribStartAge1":   "",
				"eT_TDRA_ContribStartAge2":   "",
				"eT_TDRA_ContribEndAge1":     "",
				"eT_TDRA_ContribEndAge2":     "",
				"eT_Roth1":                   "",
				"eT_Roth2":                   "",
				"eT_Roth_Rate1":              "",
				"eT_Roth_Rate2":              "",
				"eT_Roth_Contrib1":           "",
				"eT_Roth_Contrib2":           "",
				"eT_Roth_ContribStartAge1":   "",
				"eT_Roth_ContribStartAge2":   "",
				"eT_Roth_ContribEndAge1":     "",
				"eT_Roth_ContribEndAge2":     "",
				"eT_Aftatax":                 "",
				"eT_Aftatax_Rate":            "",
				"eT_Aftatax_Contrib":         "",
				"eT_Aftatax_ContribStartAge": "",
				"eT_Aftatax_ContribEndAge":   "",

				"eT_iRatePercent": "2.5",
				"eT_rRatePercent": "6",
				"eT_maximize":     "Spending", // or "PlusEstate"
			},
			verbose: true,
			iRate:   1.025,
		},
	}
	if !(testing.Short() && testing.Verbose()) { //Skip unless set "-v -short"
		/*
			ip0, err := NewInputParams(tests[0].ip, nil)
			if err != nil {
				t.Errorf("TestResultsOutput PRE case %d: %s", 0, err)
			}
			ip1, err := NewInputParams(tests[1].ip, nil)
			if err != nil {
				t.Errorf("TestResultsOutput PRE case %d: %s", 1, err)
			}
			m0 := structs.Map(ip0)
			m1 := structs.Map(ip1)
			for k, v := range m0 {
				c0 := fmt.Sprintf("%#v", v)
				c1 := fmt.Sprintf("%#v", m1[k])
				if c0 != c1 {
					fmt.Printf("\nNo match at: '%s', found\n\tm0: '%#v'\n\t m1: '%#v'\n\n", k, v, m1[k])
				}
			}*/
		t.Skip("TestResultsOutput() (full runs): skipping unless set '-v -short'")
	}
	for i, elem := range tests {
		if i != 1 {
			continue
		}
		fmt.Printf("======== CASE %d ========\n", i)
		ip, err := NewInputParams(elem.ip, nil)
		if err != nil {
			t.Errorf("TestResultsOutput case %d: %s", i, err)
			continue
		}
		//fmt.Printf("InputParams: %#v\n", ip)
		ti := NewTaxInfo(ip.FilingStatus, 2017)
		taxbins := len(*ti.Taxtable)
		cgbins := len(*ti.Capgainstable)
		vindx, err := NewVectorVarIndex(ip.Numyr, taxbins,
			cgbins, ip.Accmap, os.Stdout)
		if err != nil {
			t.Errorf("TestResultsOutput case %d: %s", i, err)
			continue
		}
		logfile, err := os.Create("ModelMatixPP.log")
		if err != nil {
			t.Errorf("TestResultsOutput case %d: %s", i, err)
			continue
		}
		//csvfile := (*os.File)(nil)
		csvfile, err := os.Create("ModelOutput.csv")
		if err != nil {
			t.Errorf("TestResultsOutput case %d: %s", i, err)
			continue
		}
		RoundToOneK := false
		ms, err := NewModelSpecs(vindx, ti, *ip,
			RoundToOneK, false, false,
			os.Stderr, logfile, csvfile, logfile, nil)
		if err != nil {
			t.Errorf("TestResultsOutput case %d: %s", i, err)
			continue
		}
		//fmt.Printf("ModelSpecs: %#v\n", ms)

		c, a, b, notes := ms.BuildModel()

		Optstart := time.Now()
		aprime, bprime, oinfo := ms.OptimizeLPModel(&a, &b)
		Optelapsed := time.Since(Optstart)

		ms.PrintModelMatrix(c, a, b, notes, nil, false, oinfo)

		tol := 1.0e-7

		bland := false
		maxiter := 4000

		callback := lpsimplex.Callbackfunc(nil)
		//callback := lpsimplex.LPSimplexVerboseCallback
		//callback := lpsimplex.LPSimplexTerseCallback
		disp := true // false //true
		start := time.Now()
		res := lpsimplex.LPSimplex(c, a, b, nil, nil, nil, callback, disp, maxiter, tol, bland)
		elapsed := time.Since(start)

		Ostart := time.Now()
		res = lpsimplex.LPSimplex(c, *aprime, *bprime, nil, nil, nil, callback, disp, maxiter, tol, bland)
		Oelapsed := time.Since(Ostart)
		/*
			err = BinDumpModel(c, a, b, res.X, "./RPlanModelgo.datX")
			if err != nil {
				t.Errorf("TestResultsOutput case %d: %s", i, err)
				continue
			}
			BinCheckModelFiles("./RPlanModelgo.datX", "./RPlanModelpython.datX", &vindx)
		*/

		//fmt.Printf("Res: %#v\n", res)
		str := fmt.Sprintf("Message: %v\n", res.Message)
		fmt.Printf(str)
		str = fmt.Sprintf("Time: LPSimplex() took %s\n", elapsed)
		fmt.Printf(str)
		str = fmt.Sprintf("Time: Opt took %s, LPSimplex() took %s\n", Optelapsed, Oelapsed)
		fmt.Printf(str)
		fmt.Printf("Called LPSimplex() for m:%d x n:%d model\n", len(a), len(a[0]))

		if res.Success {
			ms.ConsistencyCheckBrackets(&res.X)
			ms.ConsistencyCheckSpendable(&res.X)

			ms.PrintActivitySummary(&res.X)
			ms.PrintIncomeExpenseDetails()
			ms.PrintAccountTrans(&res.X)
			ms.PrintTax(&res.X)
			ms.PrintTaxBrackets(&res.X)
			ms.PrintCapGainsBrackets(&res.X)
			/*
				//ms.print_model_results(res.x)
					        if args.verboseincome:
					            print_income_expense_details()
					        if args.verboseaccounttrans:
					            print_account_trans(res)
					        if args.verbosetax:
					            print_tax(res)
					        if args.verbosetaxbrackets:
					            print_tax_brackets(res)
								print_cap_gains_brackets(res)
			*/
			ms.PrintBaseConfig(&res.X)
		}
		//createDefX(&res.X)
	}
}

func TestGenStockOutput(t *testing.T) {
	if !(testing.Short() && testing.Verbose()) { //Skip unless set "-v -short"
		t.Skip("GenStockOutput(): skipping unless set '-v -short'")
	} else {
		fmt.Printf("\nTo regenerate Stock output, use:: go test -run TestGenStockOutput -v -Short\n\n")
	}
	GenStockOutput(t)
}

func GenStockOutput(t *testing.T) {
	tests := []struct {
		ip      map[string]string
		name    string
		verbose bool
		iRate   float64
	}{
		{ // case 0 // driver generated from AWill.toml
			ip:      sipSingle,
			name:    "Single",
			verbose: true,
			iRate:   1.025,
		},
		{ // Case 1  // case to match AWill.toml Hand coded
			ip:      sipJoint,
			name:    "Joint",
			verbose: true,
			iRate:   1.025,
		},
		{ // Case 2 // case to match mobile.toml
			ip:      sipSingle3Acc,
			name:    "Single3Acc",
			verbose: true,
			iRate:   1.025,
		},
	}
	Stockfile, err := os.Create("./temp.go")
	if err != nil {
		t.Errorf("GenStockOutput():  %s", err)
		return
	}

	for i, elem := range tests {
		//if i != 0 {
		//	continue
		//}
		fmt.Printf("======== CASE %d ========\n", i)
		ip, err := NewInputParams(elem.ip, nil)
		if err != nil {
			t.Errorf("TestResultsOutput case %d: %s", i, err)
			continue
		}
		fmt.Printf("InputParams: %#v\n", ip)
		ti := NewTaxInfo(ip.FilingStatus, 2017)
		taxbins := len(*ti.Taxtable)
		cgbins := len(*ti.Capgainstable)
		vindx, err := NewVectorVarIndex(ip.Numyr, taxbins,
			cgbins, ip.Accmap, os.Stdout)
		if err != nil {
			t.Errorf("GenStockOutput case %d: %s", i, err)
			continue
		}
		logfile, err := os.Create("ModelMatixPP.log")
		if err != nil {
			t.Errorf("TestResultsOutput case %d: %s", i, err)
			continue
		}
		//csvfile := (*os.File)(nil)
		csvfile, err := os.Create("ModelOutput.csv")
		if err != nil {
			t.Errorf("TestResultsOutput case %d: %s", i, err)
			continue
		}
		RoundToOneK := false
		ms, err := NewModelSpecs(vindx, ti, *ip,
			RoundToOneK, false, false,
			os.Stderr, logfile, csvfile, logfile, nil)
		if err != nil {
			t.Errorf("TestResultsOutput case %d: %s", i, err)
			continue
		}
		//fmt.Printf("ModelSpecs: %#v\n", ms)

		c, a, b, notes := ms.BuildModel()
		ms.PrintModelMatrix(c, a, b, notes, nil, false, nil)

		tol := 1.0e-7

		bland := false
		maxiter := 4000

		callback := lpsimplex.Callbackfunc(nil)
		//callback := lpsimplex.LPSimplexVerboseCallback
		//callback := lpsimplex.LPSimplexTerseCallback
		disp := true // false //true
		start := time.Now()
		res := lpsimplex.LPSimplex(c, a, b, nil, nil, nil, callback, disp, maxiter, tol, bland)
		elapsed := time.Since(start)

		/*
			err = BinDumpModel(c, a, b, res.X, "./RPlanModelgo.datX")
			if err != nil {
				t.Errorf("TestResultsOutput case %d: %s", i, err)
				continue
			}
			BinCheckModelFiles("./RPlanModelgo.datX", "./RPlanModelpython.datX", &vindx)
		*/

		//fmt.Printf("Res: %#v\n", res)
		str := fmt.Sprintf("Message: %v\n", res.Message)
		fmt.Printf(str)
		str = fmt.Sprintf("Time: LPSimplex() took %s\n", elapsed)
		fmt.Printf(str)
		fmt.Printf("Called LPSimplex() for m:%d x n:%d model\n", len(a), len(a[0]))

		if res.Success {
			ms.PrintActivitySummary(&res.X)
			ms.PrintIncomeExpenseDetails()
			ms.PrintAccountTrans(&res.X)
			ms.PrintTax(&res.X)
			ms.PrintTaxBrackets(&res.X)
			ms.PrintCapGainsBrackets(&res.X)
			/*
				//ms.print_model_results(res.x)
					        if args.verboseincome:
					            print_income_expense_details()
					        if args.verboseaccounttrans:
					            print_account_trans(res)
					        if args.verbosetax:
					            print_tax(res)
					        if args.verbosetaxbrackets:
					            print_tax_brackets(res)
								print_cap_gains_brackets(res)
			*/
			ms.PrintBaseConfig(&res.X)
		}
		createDefX(&res.X, elem.name, Stockfile)
	}
	Stockfile.Close()
}

func createDefX(xp *[]float64, name string, f *os.File) {
	fmt.Fprintf(f, "var xp%s = &[]float64{ // sip%s InputParameters\n", name, name)
	count := 0
	for _, v := range *xp {
		count++
		fmt.Fprintf(f, "%v, ", v)
		if count > 4 {
			count = 0
			fmt.Fprintf(f, "\n")
		}

	}
	fmt.Fprintf(f, "}\n\n")
}

var xpSingle = &[]float64{ // sipSingle InputParameters
	12235.208083996806, 14713.53302947854, 0, 0, 0,
	0, 0, 12541.088286096725, 15081.371355215551, 0,
	0, 0, 0, 0, 12854.615493249143,
	15458.405639095905, 0, 0, 0, 0,
	0, 13175.98088058037, 15844.865780073236, 0, 0,
	0, 0, 0, 13505.380402594881, 16240.987424575038,
	0, 0, 0, 0, 0,
	13843.014912659752, 42493.97339140862, 0, 0, 0,
	0, 0, 14189.090285476243, 43556.32272619383, 0,
	0, 0, 0, 0, 14543.817542613151,
	9445.388367110701, 0, 0, 0, 0,
	0, 14907.412981178479, 0, 0, 0,
	0, 0, 0, 15280.09830570794, 0,
	0, 0, 0, 0, 0,
	15662.100763350638, 0, 0, 0, 0,
	0, 0, 26948.741113475324, 0, 0,
	27622.459641312245, 0, 0, 28313.02113234505, 0,
	0, 29020.84666065361, 0, 0, 29746.367827169914,
	0, 0, 56336.98830406837, 0, 0,
	57745.413011670076, 0, 0, 23989.205909723853, 0,
	0, 14907.412981178479, 0, 0, 15280.09830570794,
	0, 0, 15662.100763350638, 0, 0,
	0, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
	0, 0, 0, -1.769520233403414e-12, 0,
	0, 6837.7219214835195, 0, 0, 15237.95958436583,
	0, 0, 15618.908573974939, 0, 0,
	16009.381288324323, 0, 0, 40594.44235460848, 0,
	41609.30341347374, 0, 42649.53599881058, 0, 43715.77439878077,
	0, 44808.66875875027, 0, 71775.84675893822, 0,
	73570.24292791168, 0, 40209.656573871485, 6837.72192148352, 31533.37491192981,
	15237.95958436583, 32321.709284728055, 15618.908573974939, 33129.75201684625, 16009.381288324323,
	379659.71166708524, 0, 359409.18547122536, 0, 336867.87498121697,
	3.605778432844177e-11, 311871.4393213507, 4.307583599993721e-11, 284245.0048179243, 1.1543803924007598e-10,
	253802.51622272423, 5.2464349970412534e-11, 192948.2696316135, 23288.112114378517, 126540.70830592391,
	48555.71375847916, 91510.9148359757, 44221.07134721531, 63576.19231948863, 30722.098468620516,
	33129.75201684624, 16009.381288324323, 0, 0, 37163.89159178703,
	38092.98888158171, 39045.31360362125, 40021.44644371176, 41021.98260480454, 42047.53216992467,
	43098.7204741728, 44176.18848602709, 45280.59319817775, 46412.6080281322, 47572.92322883553,
	4.0073457699372324e-14, 0, 0, 0, 0,
	21969.917089036302, 22519.165016262174, 0, 0, 0,
	0}

var xpJoint = &[]float64{ // sipJoint InputParameters
	13303.039872342253, 0, 0, 0, 0,
	0, 0, 13635.61586915054, 0, 0,
	0, 0, 0, 0, 13976.506265879452,
	0, 0, 0, 0, 0,
	0, 14325.918922526429, 0, 0, 0,
	0, 0, 0, 14684.06689558955, 0,
	0, 0, 0, 0, 0,
	27686.029825319503, 0, 0, 0, 0,
	0, 0, 28378.180570952485, 0, 0,
	0, 0, 0, 0, 29087.6350852263,
	0, 0, 0, 0, 0,
	0, 3366.766501659296, 0, 0, 0,
	0, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
	0, 0, 13303.039872342202, 0, 0,
	13635.615869150564, 0, 0, 13976.506265879449, 0,
	0, 14325.918922526429, 0, 0, 14684.06689558955,
	0, 0, 27686.029825319503, 0, 0,
	28378.180570952485, 0, 0, 29087.6350852263, 0,
	0, 3366.766501659295, 0, 0, 0,
	0, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
	0, 0, 0, 0, 1.844100547076112e-13,
	0, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
	0, 1.0990154339522673e-12, 0, 0, 11557.526242043086,
	0, 0, 14952.306495874855, 0, 0,
	15326.114158271721, 0, 0, 40594.44235460859, 0,
	41609.30341347355, 0, 42649.53599881052, 0, 43715.77439878077,
	0, 44808.66875875025, 0, 58563.74673505922, 0,
	60027.84040343569, 0, 61528.536413521586, 0, 36618.69036316195,
	11557.526242043086, 34083.22195804023, 14952.306495874853, 34935.302506991226, 15326.114158271721,
	379659.7116670852, 0, 359409.1854712254, 6.298448531004706e-11, 336867.8749812167,
	0, 311871.43932135083, 4.8934395880033896e-11, 284245.00481792405, 3.426806610607001e-11,
	253802.51622272446, 5.209003853238387e-11, 206953.095656925, 12053.657639502653, 155740.77056869882,
	25131.876178362843, 99864.9682044879, 39303.66280656699, 67041.05451180553, 29410.90475839537,
	34935.302506991204, 15326.114158271721, 0, 0, 39264.138367374246,
	40245.74182655859, 41251.88537222254, 42283.18250652813, 43340.262069191274, 44423.76862092109,
	45534.36283644412, 46672.721907355255, 47839.5399550391, 49035.52845391506, 50261.41666526296,
	2.4061588015489116e-11, 0, 0, 0, 0,
	11371.375131606188, 11655.659509896337, 11947.05099764371, 0, 0,
	0}

var xpSingle3Acc = &[]float64{ // sipSingle3Acc InputParameters
	12235.208083996806, 28929.09187432343, 0, 0, 0,
	0, 0, 12541.088286096725, 29652.31917118132, 0,
	0, 0, 0, 0, 12854.615493249143,
	30393.627150460838, 0, 0, 0, 0,
	0, 13175.98088058037, 31153.46782922237, 0, 0,
	0, 0, 0, 13505.380402594881, 31932.304524952877,
	0, 0, 0, 0, 0,
	13843.014912659752, 26757.36822889616, 0, 0, 0,
	0, 0, 14189.09028547624, 0, 0,
	0, 0, 0, 0, 14543.81754261315,
	0, 0, 0, 0, 0,
	0, 14907.412981178479, 0, 0, 0,
	0, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
	0, 0, 41164.29995832022, 0, 0,
	42193.407457278045, 0, 0, 43248.24264370998, 0,
	0, 44329.44870980275, 0, 0, 45437.684927547765,
	0, 0, 40600.38314155591, 0, 0,
	14189.09028547624, 0, 0, 14543.81754261315, 0,
	0, 14907.412981178473, 0, 0, 0,
	0, 0, 0, 0, 0,
	1.8399220314502775e-11, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
	5077.257322803382, 0, 0, 28516.545825299338, 0,
	0, 29229.459470931804, 0, 0, 29960.195957705182,
	0, 0, 44461.28933178491, 0, 0,
	11577.18555962318, 0, 0, 54810.0011994534, 0,
	0, 56180.25122943954, 0, 2.6472742897426077e-12, 57584.75751017552,
	0, 0, 59024.37644792991, 0, 0,
	60499.985859128115, 0, 0, 56039.24159642577, 0,
	5077.257322803382, 30013.920201717847, 0, 28516.545825299338, 30764.268206760793,
	0, 29229.459470931804, 31533.374911929808, 0, 29960.195957705182,
	17041.61097902012, 0, 44461.28933178491, 17467.651253495613, 33995.63600545628,
	11577.18555962318, 379659.71166708524, 18982.98558335426, 94914.92791677131, 344340.69309568976,
	20121.964718355506, 100609.82359177763, 305450.0683782256, 21329.28260145685, 106646.41300728427,
	262737.2295201328, 22609.03955754426, 113045.19778772134, 215935.62425653526, 23965.58193099694,
	119827.90965498464, 164761.77670125137, 25403.51684685674, 127017.58423428368, 115245.88721111536,
	26927.727857668167, 129256.74652616912, 90345.8850299614, 28543.391529128272, 106784.61274292195,
	63156.51383259268, 30255.99502087598, 82208.46246830963, 33520.5272559028, 32071.354722128573,
	55383.16250124074, 17467.651253495616, 33995.63600545628, 11577.18555962318, 0,
	0, 0, 49247.116609905024, 50478.29452515266, 51740.2518882815,
	53033.75818548851, 54359.602140125695, 55718.59219362882, 57111.556998469554, 58539.34592343128,
	60002.82957151705, 61502.900310805024, 63040.47281857509, 0, 0,
	0, 0, 0, 0, 0,
	0, 0, 0, 0}
