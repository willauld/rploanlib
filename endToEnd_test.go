package rplanlib_test

import (
	"fmt"
	"os"
	"strings"
	//"regexp"
	//"strings"
	"path/filepath"
	"testing"
	"time"

	csvtag "github.com/artonge/go-csv-tag"
	"github.com/willauld/lpsimplex"
	"github.com/willauld/rplanlib"
)

type errorCode int

type acase struct {
	Testfile         string    `csv:"Testfile"`
	ErrorType        errorCode `csv:"Error Type"`
	SpendableAtLeast int       `csv:"Spendable At Least"`
}

type errorType struct {
	ErrorType errorCode `csv:"Error Code"`
	ErrorStr  string    `csv:"Error String"`
	CustomStr string    `csv:"Custom String"`
}

func createErrorTypeCSV() {
	errorTypeTable := []errorType{
		{
			ErrorType: 0,
			ErrorStr:  "no error",
		},
		{
			ErrorType: 1,
			ErrorStr:  "check spenable and no aftertax accounts",
		},
		{
			ErrorType: 2,
			ErrorStr:  "check spenable with aftertax accounts",
		},
		{
			ErrorType: 3,
			ErrorStr:  "configuration input error",
			CustomStr: "checkNames: name",
		},
	}
	err := csvtag.DumpToFile(errorTypeTable, "testdata/errortypes.csv")
	if err != nil {
		// cry
	}
}

func TestE2E(t *testing.T) {
	var updateExpectFile bool

	/*
		if !(testing.Short() && testing.Verbose()) { //Skip unless set "-v -short"
			t.Skip("TestResultsOutput() (full runs): skipping unless set '-v -short'")
		}
	*/
	/*
		//
		// Bring this back in to make sure all configuration files are
		// accounted for. Need to code this up
		//
		strmapfiles, err := filepath.Glob("./testdata/strmap/*.strmap")
		if err != nil {
			t.Errorf("TestE2E Error: %s", err)
		}
		tomlfiles, err := filepath.Glob("./testdata/toml/*.toml")
		if err != nil {
			t.Errorf("TestE2E Error: %s", err)
		}
		paramfiles := append(tomlfiles, strmapfiles...)
	*/

	cases := []acase{}
	err := csvtag.Load(csvtag.Config{ // Load your csv with configuration
		Path: "testdata/expect.csv", // Path of the csv file
		Dest: &cases,                // A pointer to the create slice
	})
	/*
				for _, ifile := range paramfiles {
					c := acase{
						testfile:         ifile,
						errorType:        0,
						spendableAtLeast: 0,
					}
					cases = append(cases, c)
				}
				err = csvtag.DumpToFile(cases, "testdata/expect.csv")
				return
		//return
	*/
	//createErrorTypeCSV() // Uncomment to update errortypes.csv
	errorTypeTable := []errorType{}
	err = csvtag.Load(csvtag.Config{ // Load your csv with configuration
		Path: "testdata/errortypes.csv", // Path of the csv file
		Dest: &errorTypeTable,           // A pointer to the create slice
	})
	//updateExpectFile = true

	for i, curCase := range cases {
		//if i != 23 {
		//	continue
		//}
		ifilecore := strings.TrimSuffix(filepath.Base(curCase.Testfile), filepath.Ext(curCase.Testfile))
		ifileext := filepath.Ext(curCase.Testfile)

		fmt.Printf("======== CASE %d - %s ========\n", i, curCase.Testfile)

		var ipsmp *map[string]string
		msgList := rplanlib.NewWarnErrorList()
		// curCase.testfile can be .toml or .strmap, Toml file is assumed
		if filepath.Ext(curCase.Testfile) == ".strmap" {
			ipsmp, err = rplanlib.GetInputStrStrMapFromFile(curCase.Testfile)
		} else {
			ipsmp, err = rplanlib.GetInputStringsMapFromToml(curCase.Testfile)
		}
		if err != nil {
			if curCase.ErrorType == 3 && strings.Contains(err.Error(), errorTypeTable[3].CustomStr) {
				// expected error
				continue
			}
			t.Errorf("TestE2E case %d: configuration file error (%s): %s", i, curCase.Testfile, err)
			rplanlib.PrintAndClearMsg(os.Stdout, msgList)
			continue
		}
		ip, err := rplanlib.NewInputParams(*ipsmp, msgList)
		if err != nil {
			t.Errorf("TestE2E case %d: %s", i, err)
			rplanlib.PrintAndClearMsg(os.Stdout, msgList)
			continue
		}
		//fmt.Printf("InputParams: %#v\n", ip)
		taxYear := 2018
		ti := rplanlib.NewTaxInfo(ip.FilingStatus, taxYear)
		taxbins := len(*ti.Taxtable)
		cgbins := len(*ti.Capgainstable)
		vindx, err := rplanlib.NewVectorVarIndex(ip.Numyr, taxbins,
			cgbins, ip.Accmap, os.Stdout)
		if err != nil {
			t.Errorf("TestE2E case %d: %s", i, err)
			continue
		}
		logname := "./testdata/" + ifileext[1:] + "_test_output/" + ifilecore + ".log"
		logfile, err := os.Create(logname)
		if err != nil {
			t.Errorf("TestE2E case %d: %s", i, err)
			continue
		}
		//csvfile := (*os.File)(nil)
		csvname := "./testdata/" + ifileext[1:] + "_test_output/" + ifilecore + ".csv"
		csvfile, err := os.Create(csvname)
		if err != nil {
			t.Errorf("TestE2E case %d: %s", i, err)
			continue
		}
		RoundToOneK := false
		allowDeposits := false
		developerInfo := true
		fourPercentRule := false
		ms, err := rplanlib.NewModelSpecs(vindx, ti, *ip,
			allowDeposits, RoundToOneK, developerInfo, fourPercentRule,
			os.Stderr, logfile, csvfile, logfile, msgList)
		if err != nil {
			t.Errorf("TestE2E case %d: %s", i, err)
			rplanlib.PrintAndClearMsg(logfile, msgList)
			continue
		}
		//fmt.Printf("ModelSpecs: %#v\n", ms)

		c, a, b, notes := ms.BuildModel()
		//c, a, b, _ := ms.BuildModel()

		var aprime *[][]float64
		var bprime *[]float64
		var Optelapsed time.Duration
		DoModelOptimizationTest := false //true //false
		if DoModelOptimizationTest {
			Optstart := time.Now()
			var oinfo *[]rplanlib.OptInfo
			aprime, bprime, oinfo = ms.OptimizeLPModel(&a, &b)
			//aprime, bprime, _ = ms.OptimizeLPModel(&a, &b)
			Optelapsed = time.Since(Optstart)

			ms.PrintModelMatrix(c, a, b, notes, nil, false, oinfo) // TODO FIXME need to make this print somewhere else for examining the optimized model
		}

		tol := 1.0e-7

		bland := false
		maxiter := 4000

		callback := lpsimplex.Callbackfunc(nil)
		//callback := lpsimplex.LPSimplexVerboseCallback
		//callback := lpsimplex.LPSimplexTerseCallback
		disp := false //true
		start := time.Now()
		res := lpsimplex.LPSimplex(c, a, b, nil, nil, nil, callback, disp, maxiter, tol, bland)
		elapsed := time.Since(start)

		var Oelapsed time.Duration
		var resPrime lpsimplex.OptResult
		if DoModelOptimizationTest {
			Ostart := time.Now()
			resPrime = lpsimplex.LPSimplex(c, *aprime, *bprime, nil, nil, nil, callback, disp, maxiter, tol, bland)
			Oelapsed = time.Since(Ostart)
		}
		/*
			err = BinDumpModel(c, a, b, res.X, "./RPlanModelgo.datX")
			if err != nil {
				t.Errorf("TestE2E case %d: %s", i, err)
				continue
			}
			BinCheckModelFiles("./RPlanModelgo.datX", "./RPlanModelpython.datX", &vindx)
		*/

		DisplayOutputAndTiming := false //true
		if DisplayOutputAndTiming || DoModelOptimizationTest {
			//fmt.Printf("Res: %#v\n", res)
			str := fmt.Sprintf("Message: %v\n", res.Message)
			fmt.Printf(str)
			if DoModelOptimizationTest {
				str = fmt.Sprintf("Message ResPrime: %v\n", resPrime.Message)
				fmt.Printf(str)
				if res.Fun == resPrime.Fun {
					fmt.Printf("Object functions match: %f\n", res.Fun)
				} else {
					fmt.Printf("Object functions DO NOT match, Standard: %f, Optimized: %f\n", res.Fun, resPrime.Fun)
				}
			}
			str = fmt.Sprintf("Time: LPSimplex() took %s\n", elapsed)
			fmt.Printf(str)
			if DoModelOptimizationTest {
				str = fmt.Sprintf("Time: Opt took %s, LPSimplex() took %s\n", Optelapsed, Oelapsed)
				fmt.Printf(str)
			}
			fmt.Printf("Called LPSimplex() for m:%d x n:%d model\n", len(a), len(a[0]))
		}
		if res.Success {
			//OK := ms.ConsistencyCheck(os.Stdout, &res.X)
			OK := ms.ConsistencyCheckBrackets(logfile, &res.X)
			if !OK {
				t.Errorf("TestE2E case %d: Check Brackets found issues with %s", i, curCase.Testfile)
			}
			OK = ms.ConsistencyCheckSpendable(logfile, &res.X)
			if !(OK || curCase.ErrorType != 0) {
				if !OK {
					if ms.Ip.Accmap[rplanlib.Aftertax] > 0 {
						foundError := 2
						if curCase.ErrorType != 2 {
							t.Errorf("TestE2E case %d: %s for file %s", i, errorTypeTable[foundError].ErrorStr, curCase.Testfile)
						}
					} else {
						foundError := 1
						if curCase.ErrorType != 1 {
							// actual error does not match expected error
							t.Errorf("TestE2E case %d: %s for file %s", i, errorTypeTable[foundError].ErrorStr, curCase.Testfile)
						}
					}
				} else {
					t.Errorf("TestE2E case %d: did not generate expected error for file %s", curCase.Testfile)
				}
			}
			s := res.X[ms.Vindx.S(0)]
			if updateExpectFile {
				newVal := int(s - 5.0)
				cases[i].SpendableAtLeast = newVal
				curCase.SpendableAtLeast = newVal
			}
			if s < float64(curCase.SpendableAtLeast) {
				t.Errorf("TestE2E case %d: first year spendable (%6.0f) is less than expected (%d) for file %s", i, s, curCase.SpendableAtLeast, curCase.Testfile)
			}

			ms.PrintActivitySummary(&res.X)
			ms.PrintIncomeExpenseDetails()
			ms.PrintAccountTrans(&res.X)
			ms.PrintTax(&res.X)
			ms.PrintTaxBrackets(&res.X)
			ms.PrintShadowTaxBrackets(&res.X)
			ms.PrintCapGainsBrackets(&res.X)
			ms.PrintAssetSummary()
			ms.PrintBaseConfig(&res.X)

			ms.PrintAccountWithdrawals(&res.X) // TESTING TESTING TESTING FIXME TODO
		}
		//createDefX(&res.X)
	}
	if updateExpectFile {
		err = csvtag.DumpToFile(cases, "testdata/expect.csv")
		if err != nil {
			t.Errorf("***** Update of expect.csv failed *******")
			// TODO FIXME
		}
	}
}