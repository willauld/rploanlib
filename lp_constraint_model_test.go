package rplanlib

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"unicode"
)

//
// Testing for lp_constraint_model.go
//

func TestIntMax(t *testing.T) {
	tests := []struct {
		a   int
		b   int
		max int
	}{
		{ // case 0
			a:   5,
			b:   6,
			max: 6,
		},
		{ // case 1
			a:   7,
			b:   6,
			max: 7,
		},
		{ // case 2
			a:   6,
			b:   6,
			max: 6,
		},
		{ // case 3
			a:   -10,
			b:   -6,
			max: -6,
		},
	}
	for i, elem := range tests {
		rmax := intMax(elem.a, elem.b)
		if rmax != elem.max {
			t.Errorf("intMax case %d: Failed - Expected %v but found %v\n", i, elem.max, rmax)
		}
	}
}

func TestIntMin(t *testing.T) {
	tests := []struct {
		a   int
		b   int
		min int
	}{
		{ // case 0
			a:   5,
			b:   6,
			min: 5,
		},
		{ // case 1
			a:   7,
			b:   6,
			min: 6,
		},
		{ // case 2
			a:   6,
			b:   6,
			min: 6,
		},
		{ // case 3
			a:   -10,
			b:   -6,
			min: -10,
		},
	}
	for i, elem := range tests {
		rmin := intMin(elem.a, elem.b)
		if rmin != elem.min {
			t.Errorf("intMin case %d: Failed - Expected %v but found %v\n", i, elem.min, rmin)
		}
	}
}

func TestCheckStrconvError(t *testing.T) {
	tests := []struct {
		//err    error
		errstr string
	}{
		{ // case 0
			errstr: "case 0",
		},
		{ // case 1
			errstr: "case 1",
		},
		{ // case 2
			errstr: "",
		},
	}
	for i, elem := range tests {
		var err error
		err = nil
		if elem.errstr != "" {
			err = fmt.Errorf(elem.errstr)
		}
		func() {
			defer func() {
				r := recover()
				if r == nil && elem.errstr != "" {
					t.Errorf("checkStrconvError case %d.a should have panicked", i)
				} else if elem.errstr == "" && r != nil {
					t.Errorf("checkStrconvError case %d.b should have panicked", i)
				} else if r != nil {
					errstr := fmt.Sprintf("%s", r)
					if errstr != elem.errstr {
						t.Errorf("checkStrconvError case %d panicked with err '%s' but should have err '%s'", i, errstr, elem.errstr)
					}
				}
			}()
			// This function should cause a panic
			checkStrconvError(err)
		}()
	}
}

func TestMergeVectors(t *testing.T) {
	tests := []struct {
		a      []float64
		b      []float64
		errstr string
	}{
		{ // Case 0
			a:      []float64{5, 2, -2, 388886, 0},
			b:      []float64{20, 30, 40, 50, 60},
			errstr: "",
		},
		{ // Case 1
			a:      []float64{5, 2, -2, 388886, 0, 20},
			b:      []float64{20, 30, 40, 50, 60},
			errstr: "mergeVectors: Can not merge, lengths do not match, 6 vs 5",
		},
	}
	for i, elem := range tests {
		newv, err := mergeVectors(elem.a, elem.b)
		if err != nil {
			if len(elem.a) == len(elem.b) {
				t.Errorf("mergeVectors case %d failed but should not have!", i)
			}
			s := fmt.Sprintf("%v", err)
			if s != elem.errstr {
				t.Errorf("mergeVectors case %d failed with incorrect err string\n\tExpected: '%s' but found: '%s'", i, elem.errstr, s)
			}
			continue
		}
		for i := 0; i < len(newv); i++ {
			if newv[i] != elem.a[i]+elem.b[i] {
				t.Errorf("mergeVectors case %d merged values do no sum", i)
			}
		}
	}
}

func TestBuildVector(t *testing.T) {
	tests := []struct {
		yearly      int
		startAge    int
		endAge      int
		vecStartAge int
		vecEndAge   int
		rate        float64
		baseAge     int
		errstr      string
	}{
		{ // case 0 // over beginning of vec
			yearly:      1,
			startAge:    45,
			endAge:      66,
			vecStartAge: 62,
			vecEndAge:   100,
			rate:        1.025,
			baseAge:     40,
			errstr:      "",
		},
		{ // case 1 // over ending of vec
			yearly:      1,
			startAge:    70,
			endAge:      102,
			vecStartAge: 62,
			vecEndAge:   100,
			rate:        1.025,
			baseAge:     40,
			errstr:      "",
		},
		{ // case 2 // in the middle of vec
			yearly:      1,
			startAge:    66,
			endAge:      68,
			vecStartAge: 62,
			vecEndAge:   100,
			rate:        1.025,
			baseAge:     40,
			errstr:      "",
		},
		{ // case 3 // all above vec
			yearly:      1,
			startAge:    145,
			endAge:      166,
			vecStartAge: 62,
			vecEndAge:   100,
			rate:        1.025,
			baseAge:     40,
			errstr:      "",
		},
		{ // case 4 // all below vec
			yearly:      1,
			startAge:    45,
			endAge:      60,
			vecStartAge: 62,
			vecEndAge:   100,
			rate:        1.025,
			baseAge:     40,
			errstr:      "",
		},
		{ // case 5 // all match vec
			yearly:      1,
			startAge:    62,
			endAge:      100,
			vecStartAge: 62,
			vecEndAge:   100,
			rate:        1.025,
			baseAge:     40,
			errstr:      "",
		},
		{ // case 6 // vec start > vec end
			yearly:      1,
			startAge:    62,
			endAge:      100,
			vecStartAge: 100,
			vecEndAge:   62,
			rate:        1.025,
			baseAge:     40,
			errstr:      "vec start age (100) is greater than vec end age (62)",
		},
		{ // case 7 // start age > end age
			yearly:      1,
			startAge:    100,
			endAge:      62,
			vecStartAge: 62,
			vecEndAge:   100,
			rate:        1.025,
			baseAge:     40,
			errstr:      "start age (100) is greater than end age (62)",
		},
	}
	for i, elem := range tests {
		newv, err := buildVector(elem.yearly, elem.startAge, elem.endAge, elem.vecStartAge, elem.vecEndAge, elem.rate, elem.baseAge)
		if err != nil {
			es := fmt.Sprintf("%s", err)
			if elem.errstr != es {
				t.Errorf("buildVector case %d: expected errstr '%s', found '%s'", i, elem.errstr, es)
			}
			// tbd TODO fix this
			//fmt.Printf("&&&&&&&&&& buildVector() returned and err for case %d: %s\n", i, err)
			continue
		}
		fnz := -1
		if elem.startAge < elem.vecEndAge && elem.startAge >= elem.vecStartAge {
			fnz = elem.startAge - elem.vecStartAge
		} else if elem.startAge < elem.vecStartAge && elem.endAge > elem.vecStartAge /*elem.vecEndAge*/ {
			fnz = 0
		}
		lnz := len(newv) + 1
		if elem.endAge < elem.vecEndAge && elem.endAge >= elem.vecStartAge {
			lnz = elem.endAge - elem.vecStartAge
		} else if elem.endAge >= elem.vecEndAge && elem.startAge <= elem.vecEndAge {
			lnz = len(newv) - 1
		}
		//fmt.Printf("CASE %d: ===================================\n", i)
		//fmt.Printf("### endAge(%d) < vecEndAge(%d) && endAge(%d) >= vecStartAge(%d)\n", elem.endAge, elem.vecEndAge, elem.endAge, elem.vecStartAge)
		//fmt.Printf("*** endAge(%d) >= vecEndAge(%d) && startAge(%d) <= vecEndAge(%d)\n", elem.endAge, elem.vecEndAge, elem.startAge, elem.vecEndAge)
		firstNonZero := -1
		lastNonZero := len(newv) + 1
		for f := 0; f < len(newv); f++ {
			if newv[f] != 0 && firstNonZero < 0 {
				firstNonZero = f
			}
			if newv[f] != 0 && firstNonZero >= 0 {
				lastNonZero = f
			}
		}
		if fnz != firstNonZero {
			t.Errorf("buildVector case %d: firstNonZero is incorrect, expected %d, found %d", i, fnz, firstNonZero)
		}
		if lnz != lastNonZero {
			t.Errorf("buildVector case %d: lastNonZero is incorrect, expected %d, found %d", i, lnz, lastNonZero)
		}
		//fmt.Printf("Case %d: newv len:%d ============\n", i, len(newv))
		//fmt.Printf("firstNonZero: %d, lastNonZero: %d\n", firstNonZero, lastNonZero)
		//fmt.Printf("fnz: %d, lnz: %d\n", fnz, lnz)
		//fmt.Printf("Case %d: %v\n", i, newv)
	}
}

func TestNewModelSpecs(t *testing.T) {
	tests := []struct {
		years   int
		ip      map[string]string
		verbose bool
		iRate   float64
	}{
		{ // Case 0 // joint
			years: 10,
			ip: map[string]string{
				"setName":                    "activeParams",
				"filingStatus":               "joint",
				"key1":                       "retiree1",
				"key2":                       "retiree2",
				"eT_Age1":                    "65",
				"eT_Age2":                    "63",
				"eT_RetireAge1":              "66",
				"eT_RetireAge2":              "66",
				"eT_PlanThroughAge1":         "100",
				"eT_PlanThroughAge2":         "100",
				"eT_PIA1":                    "30", // 30k
				"eT_PIA2":                    "-1",
				"eT_SS_Start1":               "70",
				"eT_SS_Start2":               "66",
				"eT_TDRA1":                   "200", // 200k
				"eT_TDRA2":                   "100", // 100k
				"eT_TDRA_Rate1":              "",
				"eT_TDRA_Rate2":              "",
				"eT_TDRA_Contrib1":           "",
				"eT_TDRA_Contrib2":           "",
				"eT_TDRA_ContribStartAge1":   "",
				"eT_TDRA_ContribStartAge2":   "",
				"eT_TDRA_ContribEndAge1":     "",
				"eT_TDRA_ContribEndAge2":     "",
				"eT_Roth1":                   "",
				"eT_Roth2":                   "50", // 50k
				"eT_Roth_Rate1":              "",
				"eT_Roth_Rate2":              "",
				"eT_Roth_Contrib1":           "",
				"eT_Roth_Contrib2":           "",
				"eT_Roth_ContribStartAge1":   "",
				"eT_Roth_ContribStartAge2":   "",
				"eT_Roth_ContribEndAge1":     "",
				"eT_Roth_ContribEndAge2":     "",
				"eT_Aftatax":                 "50", // 50k
				"eT_Aftatax_Rate":            "7.25",
				"eT_Aftatax_Contrib":         "",
				"eT_Aftatax_ContribStartAge": "",
				"eT_Aftatax_ContribEndAge":   "",
			},
			verbose: false,
			iRate:   1.025,
		},
		{ // Case 1 // mseparate
			years: 10,
			ip: map[string]string{
				"setName":                    "activeParams",
				"filingStatus":               "mseparate",
				"key1":                       "retiree1",
				"key2":                       "retiree2",
				"eT_Age1":                    "",
				"eT_Age2":                    "",
				"eT_RetireAge1":              "",
				"eT_RetireAge2":              "",
				"eT_PlanThroughAge1":         "",
				"eT_PlanThroughAge2":         "",
				"eT_PIA1":                    "",
				"eT_PIA2":                    "",
				"eT_SS_Start1":               "",
				"eT_SS_Start2":               "",
				"eT_TDRA1":                   "",
				"eT_TDRA2":                   "100", //100k
				"eT_TDRA_Rate1":              "",
				"eT_TDRA_Rate2":              "",
				"eT_TDRA_Contrib1":           "",
				"eT_TDRA_Contrib2":           "",
				"eT_TDRA_ContribStartAge1":   "",
				"eT_TDRA_ContribStartAge2":   "",
				"eT_TDRA_ContribEndAge1":     "",
				"eT_TDRA_ContribEndAge2":     "",
				"eT_Roth1":                   "100", //100k
				"eT_Roth2":                   "",
				"eT_Roth_Rate1":              "",
				"eT_Roth_Rate2":              "",
				"eT_Roth_Contrib1":           "",
				"eT_Roth_Contrib2":           "",
				"eT_Roth_ContribStartAge1":   "",
				"eT_Roth_ContribStartAge2":   "",
				"eT_Roth_ContribEndAge1":     "",
				"eT_Roth_ContribEndAge2":     "",
				"eT_Aftatax":                 "30", //30k
				"eT_Aftatax_Rate":            "",
				"eT_Aftatax_Contrib":         "",
				"eT_Aftatax_ContribStartAge": "",
				"eT_Aftatax_ContribEndAge":   "",
			},
			verbose: false,
			iRate:   1.025,
		},
		{ // Case 2 // single
			years: 10,
			//ip:            map[string]string{"filingStatus": "single"},
			ip: map[string]string{
				"setName":                    "activeParams",
				"filingStatus":               "single",
				"key1":                       "retiree1",
				"key2":                       "retiree2",
				"eT_Age1":                    "",
				"eT_Age2":                    "",
				"eT_RetireAge1":              "",
				"eT_RetireAge2":              "",
				"eT_PlanThroughAge1":         "",
				"eT_PlanThroughAge2":         "",
				"eT_PIA1":                    "",
				"eT_PIA2":                    "",
				"eT_SS_Start1":               "",
				"eT_SS_Start2":               "",
				"eT_TDRA1":                   "40", // 40k
				"eT_TDRA2":                   "",
				"eT_TDRA_Rate1":              "",
				"eT_TDRA_Rate2":              "",
				"eT_TDRA_Contrib1":           "",
				"eT_TDRA_Contrib2":           "",
				"eT_TDRA_ContribStartAge1":   "",
				"eT_TDRA_ContribStartAge2":   "",
				"eT_TDRA_ContribEndAge1":     "",
				"eT_TDRA_ContribEndAge2":     "",
				"eT_Roth1":                   "40", // 40k
				"eT_Roth2":                   "",
				"eT_Roth_Rate1":              "",
				"eT_Roth_Rate2":              "",
				"eT_Roth_Contrib1":           "",
				"eT_Roth_Contrib2":           "",
				"eT_Roth_ContribStartAge1":   "",
				"eT_Roth_ContribStartAge2":   "",
				"eT_Roth_ContribEndAge1":     "",
				"eT_Roth_ContribEndAge2":     "",
				"eT_Aftatax":                 "20", // 20k
				"eT_Aftatax_Rate":            "",
				"eT_Aftatax_Contrib":         "",
				"eT_Aftatax_ContribStartAge": "",
				"eT_Aftatax_ContribEndAge":   "",
			},
			verbose: false,
			iRate:   1.025,
		},
	}
	if !(testing.Short() && testing.Verbose()) { //Skip unless set "-v -short"
		t.Skip("TestNewModelSpecs(): skipping unless set '-v -short'")
	}
	for i, elem := range tests {
		ip, err := NewInputParams(elem.ip, nil)
		if err != nil {
			fmt.Printf("TestNewModelSpecs: %s\n", err)
			continue
		}
		ti := NewTaxInfo(ip.FilingStatus, 2017)
		taxbins := len(*ti.Taxtable)
		cgbins := len(*ti.Capgainstable)
		vindx, err := NewVectorVarIndex(ip.Numyr, taxbins, cgbins, ip.Accmap, os.Stdout)
		if err != nil {
			t.Errorf("TestNewModelSpecs case %d: %s", i, err)
			continue
		}
		RoundToOneK := false
		ms, err := NewModelSpecs(vindx, ti, *ip,
			RoundToOneK, false, false,
			os.Stderr, os.Stdout, nil, nil, nil)
		if err != nil {
			t.Errorf("TestNewModelSpecs case %d: %s\n", i, err)
		}
		if ms.Ip.IRate != elem.iRate {
			t.Errorf("TestNewModelSpecs case %d: iRate expected %f, found %f", i, elem.iRate, ms.Ip.IRate)
		}
	}
}

func TestBuildModel(t *testing.T) {
	tests := []struct {
		ip      map[string]string
		verbose bool
		iRate   float64
	}{
		{ // Case 0 // joint
			ip: map[string]string{
				"setName":                    "activeParams",
				"filingStatus":               "joint",
				"key1":                       "retiree1",
				"key2":                       "retiree2",
				"eT_Age1":                    "65",
				"eT_Age2":                    "63",
				"eT_RetireAge1":              "66",
				"eT_RetireAge2":              "66",
				"eT_PlanThroughAge1":         "100",
				"eT_PlanThroughAge2":         "100",
				"eT_PIA1":                    "30", // 30k
				"eT_PIA2":                    "-1",
				"eT_SS_Start1":               "70",
				"eT_SS_Start2":               "66",
				"eT_TDRA1":                   "200", // 200k
				"eT_TDRA2":                   "100", // 100k
				"eT_TDRA_Rate1":              "",
				"eT_TDRA_Rate2":              "",
				"eT_TDRA_Contrib1":           "",
				"eT_TDRA_Contrib2":           "",
				"eT_TDRA_ContribStartAge1":   "",
				"eT_TDRA_ContribStartAge2":   "",
				"eT_TDRA_ContribEndAge1":     "",
				"eT_TDRA_ContribEndAge2":     "",
				"eT_Roth1":                   "",
				"eT_Roth2":                   "50", // 50k
				"eT_Roth_Rate1":              "",
				"eT_Roth_Rate2":              "",
				"eT_Roth_Contrib1":           "",
				"eT_Roth_Contrib2":           "",
				"eT_Roth_ContribStartAge1":   "",
				"eT_Roth_ContribStartAge2":   "",
				"eT_Roth_ContribEndAge1":     "",
				"eT_Roth_ContribEndAge2":     "",
				"eT_Aftatax":                 "50", // 50k
				"eT_Aftatax_Rate":            "7.25",
				"eT_Aftatax_Contrib":         "",
				"eT_Aftatax_ContribStartAge": "",
				"eT_Aftatax_ContribEndAge":   "",
				"dollarsInThousands":         "true",
			},
			verbose: true,
			iRate:   1.025,
		},
		{ // Case 1 // joint
			ip: map[string]string{
				"setName":                    "activeParams",
				"filingStatus":               "joint",
				"key1":                       "retiree1",
				"key2":                       "retiree2",
				"eT_Age1":                    "65",
				"eT_Age2":                    "63",
				"eT_RetireAge1":              "66",
				"eT_RetireAge2":              "66",
				"eT_PlanThroughAge1":         "100",
				"eT_PlanThroughAge2":         "100",
				"eT_PIA1":                    "30", // 30k
				"eT_PIA2":                    "-1",
				"eT_SS_Start1":               "70",
				"eT_SS_Start2":               "66",
				"eT_TDRA1":                   "200", // 200k
				"eT_TDRA2":                   "100", // 100k
				"eT_TDRA_Rate1":              "",
				"eT_TDRA_Rate2":              "",
				"eT_TDRA_Contrib1":           "",
				"eT_TDRA_Contrib2":           "",
				"eT_TDRA_ContribStartAge1":   "",
				"eT_TDRA_ContribStartAge2":   "",
				"eT_TDRA_ContribEndAge1":     "",
				"eT_TDRA_ContribEndAge2":     "",
				"eT_Roth1":                   "",
				"eT_Roth2":                   "50", // 50k
				"eT_Roth_Rate1":              "",
				"eT_Roth_Rate2":              "",
				"eT_Roth_Contrib1":           "",
				"eT_Roth_Contrib2":           "",
				"eT_Roth_ContribStartAge1":   "",
				"eT_Roth_ContribStartAge2":   "",
				"eT_Roth_ContribEndAge1":     "",
				"eT_Roth_ContribEndAge2":     "",
				"eT_Aftatax":                 "50", // 50k
				"eT_Aftatax_Rate":            "7.25",
				"eT_Aftatax_Contrib":         "",
				"eT_Aftatax_ContribStartAge": "",
				"eT_Aftatax_ContribEndAge":   "",
				"eT_maximize":                "PlusEstate",
				"dollarsInThousands":         "true",
			},
			verbose: true,
			iRate:   1.025,
		},
	}
	if !(testing.Short() && testing.Verbose()) { //Skip unless set "-v -short"
		t.Skip("TestBuildModel(): skipping unless set '-v -short'")
	}
	for i, elem := range tests {
		status, err := verifyFilingStatus(elem.ip["filingStatus"])
		if err != nil {
			fmt.Printf("TestNewModelSpecs: %s\n", err)
			continue
		}
		ti := NewTaxInfo(status, 2017)
		ip, err := NewInputParams(elem.ip, nil)
		if err != nil {
			fmt.Printf("TestNewModelSpecs: %s\n", err)
			continue
		}
		fmt.Printf("ip.Maximize: %#v\n", ip.Maximize)
		taxbins := len(*ti.Taxtable)
		cgbins := len(*ti.Capgainstable)
		vindx, err := NewVectorVarIndex(ip.Numyr, taxbins,
			cgbins, ip.Accmap, os.Stdout)
		if err != nil {
			t.Errorf("BuildModel case %d: %s", i, err)
			continue
		}
		logfile, err := os.Create("ModelMatixPP.log")
		RoundToOneK := false
		ms, err := NewModelSpecs(vindx, ti, *ip,
			RoundToOneK, false, false,
			os.Stderr, logfile, nil, nil, nil)
		if err != nil {
			t.Errorf("TestNewModelSpecs case %d: %s\n", i, err)
		}
		/**/
		c, A, b, notes := ms.BuildModel()
		ms.PrintModelMatrix(c, A, b, notes, nil, false, nil)
		/**/
		if ms.Ip.IRate != elem.iRate {
			t.Errorf("BuildModel case %d: iRate expected %f, found %f", i, elem.iRate, ms.Ip.IRate)
		}
	}
}

func TestAccountOwnerAge(t *testing.T) {
	tests := []struct {
		ms    ModelSpecs
		index int
		year  int
	}{
		{ // case 0
			ms: ModelSpecs{
				Retirees: []retiree{
					{ // retireeindx == 0
						age:        56,
						ageAtStart: 57,
						throughAge: 100,
						mykey:      "retiree1",
					},
					{ // retireeindx == 1
						age:        54,
						ageAtStart: 55,
						throughAge: 100,
						mykey:      "retiree2",
					},
				},
				Accounttable: []account{
					{
						Bal:           30,
						Basis:         0,
						Contributions: []float64{},
						RRate:         1.06,
						Acctype:       IRA,
						mykey:         "retiree2",
					},
				},
			},
			index: 1,
			year:  10,
		},
		{ // case 1
			ms: ModelSpecs{
				Retirees: []retiree{
					{ // retireeindx == 0
						age:        56,
						ageAtStart: 57,
						throughAge: 100,
						mykey:      "retiree1",
					},
					{ // retireeindx == 1
						age:        54,
						ageAtStart: 55,
						throughAge: 100,
						mykey:      "retiree2",
					},
				},
				Accounttable: []account{
					{
						Bal:           30,
						Basis:         0,
						Contributions: []float64{},
						RRate:         1.06,
						Acctype:       IRA,
						mykey:         "retiree2",
					},
				},
			},
			index: 1,
			year:  7,
		},
	}
	for i, elem := range tests {
		ownerAge := elem.ms.accountOwnerAge(elem.year, elem.ms.Accounttable[0])
		calcage := elem.ms.Retirees[elem.index].ageAtStart + elem.year
		if ownerAge != calcage {
			t.Errorf("AccountOwnerAge case %d: age does not match, expected %d, found %d", i, calcage, ownerAge)
		}
	}
}

func TestMatchRetiree(t *testing.T) {
	tests := []struct {
		ms  ModelSpecs
		key string
		age int
	}{
		{ // case 0
			ms: ModelSpecs{
				Retirees: []retiree{
					{ // retireeindx == 0
						age:        56,
						ageAtStart: 57,
						throughAge: 100,
						mykey:      "retiree1",
					},
					{ // retireeindx == 1
						age:        54,
						ageAtStart: 55,
						throughAge: 100,
						mykey:      "retiree2",
					},
				},
			},
			key: "retiree2",
			age: 54,
		},
		{ // case 1
			ms: ModelSpecs{
				Retirees: []retiree{
					{ // retireeindx == 0
						age:        56,
						ageAtStart: 57,
						throughAge: 100,
						mykey:      "retiree1",
					},
					{ // retireeindx == 1
						age:        54,
						ageAtStart: 55,
						throughAge: 100,
						mykey:      "retiree2",
					},
				},
			},
			key: "retiree1",
			age: 56,
		},
	}
	for i, elem := range tests {
		r := elem.ms.matchRetiree(elem.key, 0, false)
		if r.mykey != elem.key {
			t.Errorf("MatchRetiree case %d: key does not match, expected %s, found %s", i, elem.key, r.mykey)
		}
		if r.age != elem.age {
			t.Errorf("MatchRetiree case %d: age does not match, expected %d, found %d", i, elem.age, r.age)
		}
	}
}

func TestCgTaxableFraction(t *testing.T) { /* TODO:FIXME:IMPLEMENTME */
	tests := []struct {
		ms      ModelSpecs
		expectf float64
		year    int
	}{
		{ // case 0
			ms: ModelSpecs{
				Retirees: []retiree{
					{ // retireeindx == 0
						age:        56,
						ageAtStart: 57,
						throughAge: 100,
						mykey:      "retiree1",
					},
					{ // retireeindx == 1
						age:        54,
						ageAtStart: 55,
						throughAge: 100,
						mykey:      "retiree2",
					},
				},
				Accounttable: []account{
					{
						Bal:           30,
						Basis:         20,
						Contributions: []float64{},
						RRate:         1.06,
						Acctype:       IRA,
						mykey:         "retiree2",
					},
				},
				Ip: InputParams{
					Accmap: map[Acctype]int{
						IRA:      1,
						Roth:     0,
						Aftertax: 0,
					},
				},
			},
			expectf: 1, //no aftertax account
			year:    10,
		},
		{ // case 1
			ms: ModelSpecs{
				Retirees: []retiree{
					{ // retireeindx == 0
						age:        56,
						ageAtStart: 57,
						throughAge: 100,
						mykey:      "retiree1",
					},
					{ // retireeindx == 1
						age:        54,
						ageAtStart: 55,
						throughAge: 100,
						mykey:      "retiree2",
					},
				},
				Accounttable: []account{
					{
						Bal:           30,
						Basis:         10,
						Contributions: []float64{},
						RRate:         1.06,
						Acctype:       IRA,
						mykey:         "retiree2",
					},
					{
						Bal:           30,
						Basis:         10,
						Contributions: []float64{},
						RRate:         1.06,
						Acctype:       Aftertax,
						mykey:         "retiree2",
					},
				},
				Ip: InputParams{
					Accmap: map[Acctype]int{
						IRA:      1,
						Roth:     0,
						Aftertax: 1,
					},
				},
			},
			expectf: -1, //no aftertax account
			year:    7,
		},
	}
	for i, elem := range tests {
		f := elem.ms.cgTaxableFraction(elem.year)
		fprime := elem.expectf
		if elem.expectf < 0 {
			fprime = 1 - (elem.ms.Accounttable[0].Basis / (elem.ms.Accounttable[0].Bal * math.Pow(elem.ms.Accounttable[0].RRate, float64(elem.year+elem.ms.Ip.PrePlanYears))))
		}
		if f != fprime {
			t.Errorf("cgTaxableFraction case %d: expected %f, found %f", i, fprime, f)
		}
	}
}

func TestPrintModelMatrix(t *testing.T) {
	tests := []struct {
		ip       map[string]string
		testcase string
	}{
		{ // Case 0
			ip: map[string]string{
				"setName":                    "activeParams",
				"filingStatus":               "single",
				"key1":                       "retiree1",
				"key2":                       "retiree2",
				"eT_Age1":                    "60",
				"eT_Age2":                    "",
				"eT_RetireAge1":              "65",
				"eT_RetireAge2":              "",
				"eT_PlanThroughAge1":         "75",
				"eT_PlanThroughAge2":         "",
				"eT_PIA1":                    "20", // 20k
				"eT_PIA2":                    "",
				"eT_SS_Start1":               "70",
				"eT_SS_Start2":               "",
				"eT_TDRA1":                   "10", // 10k
				"eT_TDRA2":                   "",
				"eT_TDRA_Rate1":              "",
				"eT_TDRA_Rate2":              "",
				"eT_TDRA_Contrib1":           "",
				"eT_TDRA_Contrib2":           "",
				"eT_TDRA_ContribStartAge1":   "",
				"eT_TDRA_ContribStartAge2":   "",
				"eT_TDRA_ContribEndAge1":     "",
				"eT_TDRA_ContribEndAge2":     "",
				"eT_Roth1":                   "5", // 5k
				"eT_Roth2":                   "",
				"eT_Roth_Rate1":              "",
				"eT_Roth_Rate2":              "",
				"eT_Roth_Contrib1":           "",
				"eT_Roth_Contrib2":           "",
				"eT_Roth_ContribStartAge1":   "",
				"eT_Roth_ContribStartAge2":   "",
				"eT_Roth_ContribEndAge1":     "",
				"eT_Roth_ContribEndAge2":     "",
				"eT_Aftatax":                 "15", // 15k
				"eT_Aftatax_Rate":            "",
				"eT_Aftatax_Contrib":         "",
				"eT_Aftatax_ContribStartAge": "",
				"eT_Aftatax_ContribEndAge":   "",
				"dollarsInThousands":         "true",
			},
			testcase: "allzeros",
		},
		{ // Case 1
			ip: map[string]string{
				"setName":                    "activeParams",
				"filingStatus":               "single",
				"key1":                       "retiree1",
				"key2":                       "retiree2",
				"eT_Age1":                    "60",
				"eT_Age2":                    "",
				"eT_RetireAge1":              "65",
				"eT_RetireAge2":              "",
				"eT_PlanThroughAge1":         "75",
				"eT_PlanThroughAge2":         "",
				"eT_PIA1":                    "20", // 20k
				"eT_PIA2":                    "",
				"eT_SS_Start1":               "70",
				"eT_SS_Start2":               "",
				"eT_TDRA1":                   "10", // 10k
				"eT_TDRA2":                   "",
				"eT_TDRA_Rate1":              "",
				"eT_TDRA_Rate2":              "",
				"eT_TDRA_Contrib1":           "",
				"eT_TDRA_Contrib2":           "",
				"eT_TDRA_ContribStartAge1":   "",
				"eT_TDRA_ContribStartAge2":   "",
				"eT_TDRA_ContribEndAge1":     "",
				"eT_TDRA_ContribEndAge2":     "",
				"eT_Roth1":                   "5", // 5k
				"eT_Roth2":                   "",
				"eT_Roth_Rate1":              "",
				"eT_Roth_Rate2":              "",
				"eT_Roth_Contrib1":           "",
				"eT_Roth_Contrib2":           "",
				"eT_Roth_ContribStartAge1":   "",
				"eT_Roth_ContribStartAge2":   "",
				"eT_Roth_ContribEndAge1":     "",
				"eT_Roth_ContribEndAge2":     "",
				"eT_Aftatax":                 "15", // 15k
				"eT_Aftatax_Rate":            "",
				"eT_Aftatax_Contrib":         "",
				"eT_Aftatax_ContribStartAge": "",
				"eT_Aftatax_ContribEndAge":   "",
				"dollarsInThousands":         "true",
			},
			testcase: "counting",
		},
	}
	for i, elem := range tests {
		fmt.Printf("======== Case %d ==========\n", i)
		ip, err := NewInputParams(elem.ip, nil)
		if err != nil {
			fmt.Printf("TestNewModelSpecs: %s\n", err)
			continue
		}
		ti := NewTaxInfo(ip.FilingStatus, 2017)
		taxbins := len(*ti.Taxtable)
		cgbins := len(*ti.Capgainstable)
		vindx, err := NewVectorVarIndex(ip.Numyr, taxbins, cgbins, ip.Accmap, os.Stdout)
		if err != nil {
			t.Errorf("PrintConstraint case %d: %s", i, err)
			continue
		}
		numaccounts := 0
		for _, acc := range ip.Accmap {
			numaccounts += acc
		}
		ms := ModelSpecs{
			Ip:      *ip,
			Vindx:   vindx,
			Ti:      ti,
			Logfile: os.Stdout,
			Errfile: os.Stderr,
		}

		c := make([]float64, vindx.Vsize)
		row0 := make([]float64, vindx.Vsize)
		row1 := make([]float64, vindx.Vsize)
		A := make([][]float64, 2)
		b := make([]float64, 2)

		switch elem.testcase {
		case "allones":
			for indx := 0; indx < vindx.Vsize; indx++ {
				c[indx] = 1.0
			}
			for indx := 0; indx < vindx.Vsize; indx++ {
				row0[indx] = 1.0
			}
			for indx := 0; indx < vindx.Vsize; indx++ {
				row1[indx] = 1.0
			}
			for indx := 0; indx < 2; indx++ {
				b[indx] = 1.0
			}
		case "counting":
			for indx := 0; indx < vindx.Vsize; indx++ {
				c[indx] = float64(indx)
			}
			for indx := 0; indx < vindx.Vsize; indx++ {
				row0[indx] = float64(indx)
			}
			for indx := 0; indx < vindx.Vsize; indx++ {
				row1[indx] = float64(indx + vindx.Vsize)
			}
			for indx := 0; indx < 2; indx++ {
				b[indx] = float64(indx * vindx.Vsize)
			}
		case "allzeros":
			// nothing to change
		default:
			t.Errorf("TestPrintModelMatrix: Unexpected test case '%s'\n", elem.testcase)
			continue
		}
		//fmt.Printf("Vsize: %d\n", vindx.Vsize)

		mychan := make(chan string)
		oldout, w, err := ms.RedirectModelSpecsLog(mychan)
		if err != nil {
			t.Errorf("RedirectStdout: %s\n", err)
			return // should this be continue?
		}
		A[0] = row0
		A[1] = row1
		ms.PrintModelMatrix(c, A, b, nil, nil, false, nil) // TODO add cases with varying parameters 4 and 5

		str := ms.RestoreModelSpecsLog(mychan, oldout, w)
		if !checkPrintModelOutput(ms, A, b, c, str) {
			fmt.Printf("Missmatched string:\n%s\n", str)
			t.Errorf("PrintModelMatrix Case %d: did not match expected\n", i)
		}
	}
}

func checkPrintModelOutput(ms ModelSpecs, A [][]float64, b []float64, c []float64, str string) bool {
	//
	// This code is pretty auwful. I should be improved when I have a chance
	//
	// Maybe also add some checks to use ms.Vindx... or other data from ms.
	//
	pass := true
	type checkType int
	const (
		cobj        checkType = iota
		constraints checkType = iota
		rhs         checkType = iota
		unknown     checkType = iota
	)
	curCheck := unknown
	constraintNo := -1
	varNo := -1
	//c[194] = 4.0
	//A[0][20] = 5.0
	//b[0] = 20.0
	scanner := bufio.NewScanner(strings.NewReader(str))
	for scanner.Scan() {
		x := scanner.Text()
		if !strings.Contains(x, "B?") {
			for _, token := range strings.Split(strings.TrimSuffix(x, "\n"), ", ") {
				for _, elem := range strings.Split(strings.TrimSuffix(token, "\n"), ": ") {
					if curCheck == rhs {
						// elem should be a rhs value compare it with b[constraintNo]
						if constraintNo == -1 && varNo == -1 {
							// we are not in a full model printout just a row so let
							// constraintNo == 0
							constraintNo = 0
						}
						f, _ := strconv.ParseFloat(elem, 64)
						if math.Abs(f-b[constraintNo]) > 0.01 {
							// Print B miss compares ****************************
							fmt.Printf("PrintModel:: b[%d]: %f != print output vector value %f\n", constraintNo, b[constraintNo], f)
							pass = false
						}
						//fmt.Printf("**** rhs: %f \n", f)
						curCheck = unknown
						varNo = -1
						constraintNo = -1
					} else {
						if curCheck == unknown {
							matched, err := regexp.MatchString("c", elem)
							if matched || err != nil {
								if err != nil {
									fmt.Printf("MatchString()1 Error: %s\n", err)
								} else {
									curCheck = cobj
									varNo = 0
								}
							} else {
								matched, err := regexp.MatchString("[0123456789]+", elem)
								if matched || err != nil {
									if err != nil {
										fmt.Printf("MatchString()2 Error: %s\n", err)
									} else {
										curCheck = constraints
										varNo = 0
										ruleNo, _ := strconv.ParseInt(strings.TrimSpace(elem), 10, 32)
										constraintNo = int(ruleNo)
										//fmt.Printf("ConstraintNo: %d, elem: '%s'\n", constraintNo, elem)
									}
								} else {
									matched, err := regexp.MatchString("[.*]=", elem)
									if matched || err != nil {
										if err != nil {
											fmt.Printf("MatchString()4 Error: %s\n", err)
										} else {
											curCheck = constraints
											varNo = 0
											//ruleNo, _ := strconv.ParseInt(strings.TrimSpace(elem), 10, 32)
											constraintNo = 0
											//fmt.Printf("ConstraintNo: %d, elem: '%s'\n", constraintNo, elem)
										}
									}
								}
							}
						}
						if elem == "" {
							curCheck = unknown
							varNo = -1
							constraintNo = -1
						} else {
							matched, err := regexp.MatchString("<=", elem)
							if matched || err != nil {
								if err != nil {
									fmt.Printf("MatchString()3 Error: %s\n", err)
								} else {
									curCheck = rhs
								}
							}
						}
					}
					//fmt.Println(curCheck, constraintNo, varNo, elem)
					if curCheck == constraints || curCheck == cobj {
						pieces := strings.Split(elem, "=")
						//fmt.Println(pieces)
						if len(pieces) == 2 {
							f, _ := strconv.ParseFloat(strings.TrimSpace(pieces[1]), 64)
							// skip over values that are 0.0 and hence not printed
							val := 0.0
							for i := varNo; i < len(A[0]); i++ {
								if curCheck == constraints {
									val = A[constraintNo][varNo]
								} else {
									val = c[varNo]
								}
								if val != 0.0 {
									break
								}
								varNo++
							}
							if curCheck == constraints {
								//fmt.Printf("*** constraintNo: %d, varNo: %d\n", constraintNo, varNo)
								if math.Abs(A[constraintNo][varNo]-f) > 0.01 {
									// Print A miss compares ****************************
									fmt.Printf("PrintModel:: A[%d][%d]: %f != print output vector value %f\n", constraintNo, varNo, A[constraintNo][varNo], f)
									pass = false
								}
							} else {
								//fmt.Printf("************* c[%d]: %f, elemval: %f, diff of %f\n", varNo, c[varNo], f, c[varNo]-f)
								if math.Abs(c[varNo]-f) > 0.01 {
									// Print A miss compares ****************************
									fmt.Printf("PrintModel:: c[%d]: %f != print output vector value %f\n", varNo, c[varNo], f)
									pass = false
								}
							}
							varNo++
						}
					}
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("checkModelOutput: Error reading input string:", err)
		return false
	}
	return pass == true
}

func TestPrintConstraint(t *testing.T) {
	tests := []struct {
		ip        map[string]string
		b         float64
		expectstr string
		testcase  string
	}{
		{ // Case 0
			ip: map[string]string{
				"setName":                    "activeParams",
				"filingStatus":               "single",
				"key1":                       "retiree1",
				"key2":                       "retiree2",
				"eT_Age1":                    "60",
				"eT_Age2":                    "",
				"eT_RetireAge1":              "65",
				"eT_RetireAge2":              "",
				"eT_PlanThroughAge1":         "75",
				"eT_PlanThroughAge2":         "",
				"eT_PIA1":                    "20", // 20k
				"eT_PIA2":                    "",
				"eT_SS_Start1":               "70",
				"eT_SS_Start2":               "",
				"eT_TDRA1":                   "10", // 10k
				"eT_TDRA2":                   "",
				"eT_TDRA_Rate1":              "",
				"eT_TDRA_Rate2":              "",
				"eT_TDRA_Contrib1":           "",
				"eT_TDRA_Contrib2":           "",
				"eT_TDRA_ContribStartAge1":   "",
				"eT_TDRA_ContribStartAge2":   "",
				"eT_TDRA_ContribEndAge1":     "",
				"eT_TDRA_ContribEndAge2":     "",
				"eT_Roth1":                   "5", // 5k
				"eT_Roth2":                   "",
				"eT_Roth_Rate1":              "",
				"eT_Roth_Rate2":              "",
				"eT_Roth_Contrib1":           "",
				"eT_Roth_Contrib2":           "",
				"eT_Roth_ContribStartAge1":   "",
				"eT_Roth_ContribStartAge2":   "",
				"eT_Roth_ContribEndAge1":     "",
				"eT_Roth_ContribEndAge2":     "",
				"eT_Aftatax":                 "15", // 15k
				"eT_Aftatax_Rate":            "",
				"eT_Aftatax_Contrib":         "",
				"eT_Aftatax_ContribStartAge": "",
				"eT_Aftatax_ContribEndAge":   "",
				"dollarsInThousands":         "true",
			},
			b: 743.027,
			expectstr: `Row: [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
			b: 743.027
			<= b[]: 743.03
			`,
			testcase: "allzeros",
		},
		{ // Case 1
			ip: map[string]string{
				"setName":                    "activeParams",
				"filingStatus":               "single",
				"key1":                       "retiree1",
				"key2":                       "retiree2",
				"eT_Age1":                    "60",
				"eT_Age2":                    "",
				"eT_RetireAge1":              "65",
				"eT_RetireAge2":              "",
				"eT_PlanThroughAge1":         "75",
				"eT_PlanThroughAge2":         "",
				"eT_PIA1":                    "20", // 20k
				"eT_PIA2":                    "",
				"eT_SS_Start1":               "70",
				"eT_SS_Start2":               "",
				"eT_TDRA1":                   "10", // 10k
				"eT_TDRA2":                   "",
				"eT_TDRA_Rate1":              "",
				"eT_TDRA_Rate2":              "",
				"eT_TDRA_Contrib1":           "",
				"eT_TDRA_Contrib2":           "",
				"eT_TDRA_ContribStartAge1":   "",
				"eT_TDRA_ContribStartAge2":   "",
				"eT_TDRA_ContribEndAge1":     "",
				"eT_TDRA_ContribEndAge2":     "",
				"eT_Roth1":                   "5", // 5k
				"eT_Roth2":                   "",
				"eT_Roth_Rate1":              "",
				"eT_Roth_Rate2":              "",
				"eT_Roth_Contrib1":           "",
				"eT_Roth_Contrib2":           "",
				"eT_Roth_ContribStartAge1":   "",
				"eT_Roth_ContribStartAge2":   "",
				"eT_Roth_ContribEndAge1":     "",
				"eT_Roth_ContribEndAge2":     "",
				"eT_Aftatax":                 "15", // 15k
				"eT_Aftatax_Rate":            "",
				"eT_Aftatax_Contrib":         "",
				"eT_Aftatax_ContribStartAge": "",
				"eT_Aftatax_ContribEndAge":   "",
				"dollarsInThousands":         "true",
			},
			b: 743.027,
			expectstr: `Row: [0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31 32 33 34 35 36 37 38 39 40 41 42 43 44 45 46 47 48 49 50 51 52 53 54 55 56 57 58 59 60 61 62 63 64 65 66 67 68 69 70 71 72 73 74 75 76 77 78 79 80 81 82 83 84 85 86 87 88 89 90 91 92 93 94 95 96
			97 98 99 100 101 102 103 104 105 106 107 108 109 110 111 112 113 114 115 116 117 118 119 120 121 122 123 124 125 126 127 128 129 130 131 132 133 134 135 136 137 138 139 140 141 142 143 144 145 146 147 148 149 150 151 152 153 154 155 156 157 158 159 160 161 162 163 164 165 166 167 168 169 170 171 172 173 174 175 176 177 178 179 180 181 182 183 184 185 186 187 188 189 190 191 192 193 194 195 196 197 198 199 200 201 202 203 204 205 206 207 208 209 210 211 212 213 214 215 216 217 218 219 220 221 222 223 224 225 226 227 228 229 230 231 232 233 234 235 236 237 238 239 240 241 242 243 244 245 246 247 248 249 250 251 252 253 254 255]
							b: 743.027
							x[0,1]= 1.000, x[0,2]= 2.000, x[0,3]= 3.000, x[0,4]= 4.000, x[0,5]= 5.000, x[0,6]= 6.000, x[1,0]= 7.000, x[1,1]= 8.000, x[1,2]= 9.000, x[1,3]=10.000, x[1,4]=11.000, x[1,5]=12.000, x[1,6]=13.000, x[2,0]=14.000, x[2,1]=15.000, x[2,2]=16.000, x[2,3]=17.000, x[2,4]=18.000, x[2,5]=19.000, x[2,6]=20.000, x[3,0]=21.000, x[3,1]=22.000, x[3,2]=23.000, x[3,3]=24.000, x[3,4]=25.000, x[3,5]=26.000, x[3,6]=27.000, x[4,0]=28.000, x[4,1]=29.000, x[4,2]=30.000, x[4,3]=31.000, x[4,4]=32.000, x[4,5]=33.000, x[4,6]=34.000, x[5,0]=35.000, x[5,1]=36.000, x[5,2]=37.000, x[5,3]=38.000, x[5,4]=39.000, x[5,5]=40.000, x[5,6]=41.000, x[6,0]=42.000, x[6,1]=43.000, x[6,2]=44.000, x[6,3]=45.000, x[6,4]=46.000, x[6,5]=47.000, x[6,6]=48.000, x[7,0]=49.000, x[7,1]=50.000, x[7,2]=51.000, x[7,3]=52.000, x[7,4]=53.000, x[7,5]=54.000, x[7,6]=55.000, x[8,0]=56.000, x[8,1]=57.000, x[8,2]=58.000, x[8,3]=59.000, x[8,4]=60.000, x[8,5]=61.000, x[8,6]=62.000, x[9,0]=63.000, x[9,1]=64.000, x[9,2]=65.000, x[9,3]=66.000, x[9,4]=67.000, x[9,5]=68.000, x[9,6]=69.000, x[10,0]=70.000, x[10,1]=71.000, x[10,2]=72.000, x[10,3]=73.000, x[10,4]=74.000, x[10,5]=75.000, x[10,6]=76.000, sy[0,0]=77.000, sy[0,1]=78.000, sy[0,2]=79.000, sy[1,0]=80.000, sy[1,1]=81.000, sy[1,2]=82.000, sy[2,0]=83.000, sy[2,1]=84.000, sy[2,2]=85.000, sy[3,0]=86.000, sy[3,1]=87.000, sy[3,2]=88.000, sy[4,0]=89.000, sy[4,1]=90.000, sy[4,2]=91.000, sy[5,0]=92.000, sy[5,1]=93.000, sy[5,2]=94.000, sy[6,0]=95.000, sy[6,1]=96.000, sy[6,2]=97.000, sy[7,0]=98.000, sy[7,1]=99.000, sy[7,2]=100.000, sy[8,0]=101.000, sy[8,1]=102.000, sy[8,2]=103.000, sy[9,0]=104.000, sy[9,1]=105.000, sy[9,2]=106.000, sy[10,0]=107.000, sy[10,1]=108.000, sy[10,2]=109.000, y[0,0]=110.000, y[0,1]=111.000, y[0,2]=112.000, y[1,0]=113.000,
			y[1,1]=114.000, y[1,2]=115.000, y[2,0]=116.000, y[2,1]=117.000, y[2,2]=118.000, y[3,0]=119.000, y[3,1]=120.000, y[3,2]=121.000, y[4,0]=122.000, y[4,1]=123.000, y[4,2]=124.000, y[5,0]=125.000, y[5,1]=126.000,
			y[5,2]=127.000, y[6,0]=128.000, y[6,1]=129.000, y[6,2]=130.000, y[7,0]=131.000, y[7,1]=132.000, y[7,2]=133.000, y[8,0]=134.000, y[8,1]=135.000, y[8,2]=136.000, y[9,0]=137.000, y[9,1]=138.000, y[9,2]=139.000,
			y[10,0]=140.000, y[10,1]=141.000, y[10,2]=142.000, w[0,0]=143.000, w[0,1]=144.000, w[0,2]=145.000, w[1,0]=146.000, w[1,1]=147.000, w[1,2]=148.000, w[2,0]=149.000, w[2,1]=150.000, w[2,2]=151.000, w[3,0]=152.000, w[3,1]=153.000, w[3,2]=154.000, w[4,0]=155.000, w[4,1]=156.000, w[4,2]=157.000, w[5,0]=158.000, w[5,1]=159.000, w[5,2]=160.000, w[6,0]=161.000, w[6,1]=162.000, w[6,2]=163.000, w[7,0]=164.000, w[7,1]=165.000, w[7,2]=166.000, w[8,0]=167.000, w[8,1]=168.000, w[8,2]=169.000, w[9,0]=170.000, w[9,1]=171.000, w[9,2]=172.000, w[10,0]=173.000, w[10,1]=174.000, w[10,2]=175.000, b[0,0]=176.000, b[0,1]=177.000, b[0,2]=178.000, b[1,0]=179.000, b[1,1]=180.000, b[1,2]=181.000, b[2,0]=182.000, b[2,1]=183.000, b[2,2]=184.000, b[3,0]=185.000, b[3,1]=186.000, b[3,2]=187.000, b[4,0]=188.000, b[4,1]=189.000, b[4,2]=190.000, b[5,0]=191.000, b[5,1]=192.000, b[5,2]=193.000, b[6,0]=194.000, b[6,1]=195.000, b[6,2]=196.000, b[7,0]=197.000, b[7,1]=198.000, b[7,2]=199.000, b[8,0]=200.000, b[8,1]=201.000, b[8,2]=202.000, b[9,0]=203.000, b[9,1]=204.000, b[9,2]=205.000, b[10,0]=206.000, b[10,1]=207.000, b[10,2]=208.000, b[11,0]=209.000, b[11,1]=210.000, b[11,2]=211.000, s[0]=212.000, s[1]=213.000, s[2]=214.000, s[3]=215.000, s[4]=216.000, s[5]=217.000,
			s[6]=218.000, s[7]=219.000, s[8]=220.000, s[9]=221.000, s[10]=222.000, D[0,0]=223.000, D[0,1]=224.000, D[0,2]=225.000, D[1,0]=226.000, D[1,1]=227.000, D[1,2]=228.000, D[2,0]=229.000, D[2,1]=230.000, D[2,2]=231.000, D[3,0]=232.000, D[3,1]=233.000, D[3,2]=234.000, D[4,0]=235.000, D[4,1]=236.000, D[4,2]=237.000, D[5,0]=238.000, D[5,1]=239.000, D[5,2]=240.000, D[6,0]=241.000, D[6,1]=242.000, D[6,2]=243.000, D[7,0]=244.000, D[7,1]=245.000, D[7,2]=246.000, D[8,0]=247.000, D[8,1]=248.000, D[8,2]=249.000, D[9,0]=250.000, D[9,1]=251.000, D[9,2]=252.000, D[10,0]=253.000, D[10,1]=254.000, D[10,2]=255.000, <= b[]: 743.03
							`,
			testcase: "counting",
		},
	}
	for i, elem := range tests {
		ip, err := NewInputParams(elem.ip, nil)
		if err != nil {
			fmt.Printf("TestNewModelSpecs: %s\n", err)
			continue
		}
		ti := NewTaxInfo(ip.FilingStatus, 2017)
		taxbins := len(*ti.Taxtable)
		cgbins := len(*ti.Capgainstable)
		vindx, err := NewVectorVarIndex(ip.Numyr, taxbins, cgbins, ip.Accmap, os.Stdout)
		if err != nil {
			t.Errorf("PrintConstraint case %d: %s", i, err)
			continue
		}
		numaccounts := 0
		for _, acc := range ip.Accmap {
			numaccounts += acc
		}
		ms := ModelSpecs{
			Ip:      *ip,
			Vindx:   vindx,
			Ti:      ti,
			Logfile: os.Stdout,
			Errfile: os.Stderr,
		}

		row := make([]float64, vindx.Vsize)
		switch elem.testcase {
		case "allones":
			for indx := 0; indx < vindx.Vsize; indx++ {
				row[indx] = float64(indx)
			}
		case "counting":
			for indx := 0; indx < vindx.Vsize; indx++ {
				row[indx] = float64(indx)
			}
		case "allzeros":
			// nothing to change
		default:
			t.Errorf("TestPrintConstraint: Unexpected test case '%s'\n", elem.testcase)
			continue
		}

		mychan := make(chan string)
		//oldout, w, err := RedirectStdout(mychan)
		oldout, w, err := ms.RedirectModelSpecsLog(mychan)
		if err != nil {
			t.Errorf("RedirectStdout: %s\n", err)
			return // should this be continue?
		}
		ms.printConstraint(row, elem.b)

		//str := RestoreStdout(mychan, oldout, w)
		str := ms.RestoreModelSpecsLog(mychan, oldout, w)
		c := make([]float64, len(row))
		A := make([][]float64, 0)
		A = append(A, row)
		b := make([]float64, 1)
		b[0] = elem.b
		if !checkPrintModelOutput(ms, A, b, c, str) {
			t.Errorf("PrintConstraint Case %d: print model does not match", i)
		}
	}
}

func TestPrintModelRow(t *testing.T) {
	tests := []struct {
		ip              map[string]string
		suppressNewline bool
		expectstr       string
		testcase        string
	}{
		{ // Case 0
			ip: map[string]string{
				"setName":                    "activeParams",
				"filingStatus":               "single",
				"key1":                       "retiree1",
				"key2":                       "retiree2",
				"eT_Age1":                    "60",
				"eT_Age2":                    "",
				"eT_RetireAge1":              "65",
				"eT_RetireAge2":              "",
				"eT_PlanThroughAge1":         "75",
				"eT_PlanThroughAge2":         "",
				"eT_PIA1":                    "20", // 20k
				"eT_PIA2":                    "",
				"eT_SS_Start1":               "70",
				"eT_SS_Start2":               "",
				"eT_TDRA1":                   "10", // 10k
				"eT_TDRA2":                   "",
				"eT_TDRA_Rate1":              "",
				"eT_TDRA_Rate2":              "",
				"eT_TDRA_Contrib1":           "",
				"eT_TDRA_Contrib2":           "",
				"eT_TDRA_ContribStartAge1":   "",
				"eT_TDRA_ContribStartAge2":   "",
				"eT_TDRA_ContribEndAge1":     "",
				"eT_TDRA_ContribEndAge2":     "",
				"eT_Roth1":                   "5", // 5k
				"eT_Roth2":                   "",
				"eT_Roth_Rate1":              "",
				"eT_Roth_Rate2":              "",
				"eT_Roth_Contrib1":           "",
				"eT_Roth_Contrib2":           "",
				"eT_Roth_ContribStartAge1":   "",
				"eT_Roth_ContribStartAge2":   "",
				"eT_Roth_ContribEndAge1":     "",
				"eT_Roth_ContribEndAge2":     "",
				"eT_Aftatax":                 "15", // 15k
				"eT_Aftatax_Rate":            "",
				"eT_Aftatax_Contrib":         "",
				"eT_Aftatax_ContribStartAge": "",
				"eT_Aftatax_ContribEndAge":   "",
				"dollarsInThousands":         "true",
			},
			suppressNewline: false,
			expectstr: `Row: [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]

			`,
			testcase: "allzeros",
		},
		{ // Case 1
			ip: map[string]string{
				"setName":                    "activeParams",
				"filingStatus":               "single",
				"key1":                       "retiree1",
				"key2":                       "retiree2",
				"eT_Age1":                    "60",
				"eT_Age2":                    "",
				"eT_RetireAge1":              "65",
				"eT_RetireAge2":              "",
				"eT_PlanThroughAge1":         "75",
				"eT_PlanThroughAge2":         "",
				"eT_PIA1":                    "20", // 20k
				"eT_PIA2":                    "",
				"eT_SS_Start1":               "70",
				"eT_SS_Start2":               "",
				"eT_TDRA1":                   "10", // 10k
				"eT_TDRA2":                   "",
				"eT_TDRA_Rate1":              "",
				"eT_TDRA_Rate2":              "",
				"eT_TDRA_Contrib1":           "",
				"eT_TDRA_Contrib2":           "",
				"eT_TDRA_ContribStartAge1":   "",
				"eT_TDRA_ContribStartAge2":   "",
				"eT_TDRA_ContribEndAge1":     "",
				"eT_TDRA_ContribEndAge2":     "",
				"eT_Roth1":                   "5", // 5k
				"eT_Roth2":                   "",
				"eT_Roth_Rate1":              "",
				"eT_Roth_Rate2":              "",
				"eT_Roth_Contrib1":           "",
				"eT_Roth_Contrib2":           "",
				"eT_Roth_ContribStartAge1":   "",
				"eT_Roth_ContribStartAge2":   "",
				"eT_Roth_ContribEndAge1":     "",
				"eT_Roth_ContribEndAge2":     "",
				"eT_Aftatax":                 "15", // 15k
				"eT_Aftatax_Rate":            "",
				"eT_Aftatax_Contrib":         "",
				"eT_Aftatax_ContribStartAge": "",
				"eT_Aftatax_ContribEndAge":   "",
				"dollarsInThousands":         "true",
			},
			suppressNewline: true,
			expectstr: `Row: [0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31 32 33 34 35 36 37 38 39 40 41 42 43 44 45 46 47 48 49
			50 51 52 53 54 55 56 57 58 59 60 61 62 63 64 65 66 67 68 69 70 71 72 73 74 75 76 77 78 79 80 81 82 83 84 85 86 87 88 89 90 91 92 93 94 95 96 97 98 99 100 101 102 103 104 105 106 107 108 109 110 111 112 113 114 115 116 117 118 119 120 121 122 123 124 125 126 127 128 129 130 131 132 133 134 135 136 137 138 139 140 141 142 143 144 145 146 147 148 149 150 151 152 153 154 155 156 157 158 159 160 161 162 163 164 165 166 167 168 169 170 171 172 173 174 175 176 177 178 179 180 181 182 183 184 185 186 187 188 189 190 191 192 193 194 195 196 197 198 199 200 201 202 203 204 205 206 207 208 209 210 211 212 213 214 215 216 217 218 219 220 221 222 223 224 225 226 227 228 229 230 231 232 233 234 235 236 237 238 239 240 241 242 243 244 245 246 247 248 249 250 251 252 253 254 255]
							x[0,1]= 1.000, x[0,2]= 2.000, x[0,3]= 3.000, x[0,4]= 4.000, x[0,5]= 5.000, x[0,6]= 6.000, x[1,0]= 7.000, x[1,1]= 8.000, x[1,2]= 9.000, x[1,3]=10.000, x[1,4]=11.000, x[1,5]=12.000, x[1,6]=13.000, x[2,0]=14.000, x[2,1]=15.000, x[2,2]=16.000, x[2,3]=17.000, x[2,4]=18.000, x[2,5]=19.000, x[2,6]=20.000, x[3,0]=21.000, x[3,1]=22.000, x[3,2]=23.000, x[3,3]=24.000, x[3,4]=25.000, x[3,5]=26.000, x[3,6]=27.000, x[4,0]=28.000, x[4,1]=29.000, x[4,2]=30.000, x[4,3]=31.000, x[4,4]=32.000, x[4,5]=33.000, x[4,6]=34.000, x[5,0]=35.000, x[5,1]=36.000, x[5,2]=37.000, x[5,3]=38.000, x[5,4]=39.000, x[5,5]=40.000, x[5,6]=41.000, x[6,0]=42.000, x[6,1]=43.000, x[6,2]=44.000, x[6,3]=45.000, x[6,4]=46.000, x[6,5]=47.000, x[6,6]=48.000, x[7,0]=49.000, x[7,1]=50.000, x[7,2]=51.000, x[7,3]=52.000, x[7,4]=53.000, x[7,5]=54.000, x[7,6]=55.000, x[8,0]=56.000, x[8,1]=57.000, x[8,2]=58.000, x[8,3]=59.000, x[8,4]=60.000, x[8,5]=61.000, x[8,6]=62.000, x[9,0]=63.000, x[9,1]=64.000, x[9,2]=65.000, x[9,3]=66.000, x[9,4]=67.000, x[9,5]=68.000, x[9,6]=69.000, x[10,0]=70.000, x[10,1]=71.000, x[10,2]=72.000, x[10,3]=73.000, x[10,4]=74.000, x[10,5]=75.000, x[10,6]=76.000, sy[0,0]=77.000, sy[0,1]=78.000, sy[0,2]=79.000, sy[1,0]=80.000, sy[1,1]=81.000, sy[1,2]=82.000, sy[2,0]=83.000, sy[2,1]=84.000, sy[2,2]=85.000, sy[3,0]=86.000, sy[3,1]=87.000, sy[3,2]=88.000, sy[4,0]=89.000, sy[4,1]=90.000, sy[4,2]=91.000, sy[5,0]=92.000, sy[5,1]=93.000, sy[5,2]=94.000, sy[6,0]=95.000, sy[6,1]=96.000, sy[6,2]=97.000, sy[7,0]=98.000, sy[7,1]=99.000, sy[7,2]=100.000, sy[8,0]=101.000, sy[8,1]=102.000, sy[8,2]=103.000, sy[9,0]=104.000, sy[9,1]=105.000, sy[9,2]=106.000, sy[10,0]=107.000, sy[10,1]=108.000, sy[10,2]=109.000, y[0,0]=110.000, y[0,1]=111.000, y[0,2]=112.000, y[1,0]=113.000,
			y[1,1]=114.000, y[1,2]=115.000, y[2,0]=116.000, y[2,1]=117.000, y[2,2]=118.000, y[3,0]=119.000, y[3,1]=120.000, y[3,2]=121.000, y[4,0]=122.000, y[4,1]=123.000, y[4,2]=124.000, y[5,0]=125.000, y[5,1]=126.000,
			y[5,2]=127.000, y[6,0]=128.000, y[6,1]=129.000, y[6,2]=130.000, y[7,0]=131.000, y[7,1]=132.000, y[7,2]=133.000, y[8,0]=134.000, y[8,1]=135.000, y[8,2]=136.000, y[9,0]=137.000, y[9,1]=138.000, y[9,2]=139.000,
			y[10,0]=140.000, y[10,1]=141.000, y[10,2]=142.000, w[0,0]=143.000, w[0,1]=144.000, w[0,2]=145.000, w[1,0]=146.000, w[1,1]=147.000, w[1,2]=148.000, w[2,0]=149.000, w[2,1]=150.000, w[2,2]=151.000, w[3,0]=152.000, w[3,1]=153.000, w[3,2]=154.000, w[4,0]=155.000, w[4,1]=156.000, w[4,2]=157.000, w[5,0]=158.000, w[5,1]=159.000, w[5,2]=160.000, w[6,0]=161.000, w[6,1]=162.000, w[6,2]=163.000, w[7,0]=164.000, w[7,1]=165.000, w[7,2]=166.000, w[8,0]=167.000, w[8,1]=168.000, w[8,2]=169.000, w[9,0]=170.000, w[9,1]=171.000, w[9,2]=172.000, w[10,0]=173.000, w[10,1]=174.000, w[10,2]=175.000, b[0,0]=176.000, b[0,1]=177.000, b[0,2]=178.000, b[1,0]=179.000, b[1,1]=180.000, b[1,2]=181.000, b[2,0]=182.000, b[2,1]=183.000, b[2,2]=184.000, b[3,0]=185.000, b[3,1]=186.000, b[3,2]=187.000, b[4,0]=188.000, b[4,1]=189.000, b[4,2]=190.000, b[5,0]=191.000, b[5,1]=192.000, b[5,2]=193.000, b[6,0]=194.000, b[6,1]=195.000, b[6,2]=196.000, b[7,0]=197.000, b[7,1]=198.000, b[7,2]=199.000, b[8,0]=200.000, b[8,1]=201.000, b[8,2]=202.000, b[9,0]=203.000, b[9,1]=204.000, b[9,2]=205.000, b[10,0]=206.000, b[10,1]=207.000, b[10,2]=208.000, b[11,0]=209.000, b[11,1]=210.000, b[11,2]=211.000, s[0]=212.000, s[1]=213.000, s[2]=214.000, s[3]=215.000, s[4]=216.000, s[5]=217.000,
			s[6]=218.000, s[7]=219.000, s[8]=220.000, s[9]=221.000, s[10]=222.000, D[0,0]=223.000, D[0,1]=224.000, D[0,2]=225.000, D[1,0]=226.000, D[1,1]=227.000, D[1,2]=228.000, D[2,0]=229.000, D[2,1]=230.000, D[2,2]=231.000, D[3,0]=232.000, D[3,1]=233.000, D[3,2]=234.000, D[4,0]=235.000, D[4,1]=236.000, D[4,2]=237.000, D[5,0]=238.000, D[5,1]=239.000, D[5,2]=240.000, D[6,0]=241.000, D[6,1]=242.000, D[6,2]=243.000, D[7,0]=244.000, D[7,1]=245.000, D[7,2]=246.000, D[8,0]=247.000, D[8,1]=248.000, D[8,2]=249.000, D[9,0]=250.000, D[9,1]=251.000, D[9,2]=252.000, D[10,0]=253.000, D[10,1]=254.000, D[10,2]=255.000, `,
			testcase: "counting",
		},
	}
	for i, elem := range tests {
		ip, err := NewInputParams(elem.ip, nil)
		if err != nil {
			fmt.Printf("TestNewModelSpecs: %s\n", err)
			continue
		}
		ti := NewTaxInfo(ip.FilingStatus, 2017)
		taxbins := len(*ti.Taxtable)
		cgbins := len(*ti.Capgainstable)
		vindx, err := NewVectorVarIndex(ip.Numyr, taxbins, cgbins, ip.Accmap, os.Stdout)
		if err != nil {
			t.Errorf("PrintModelRow case %d: %s", i, err)
			continue
		}
		ms := ModelSpecs{
			Ip:      *ip,
			Vindx:   vindx,
			Ti:      ti,
			Logfile: os.Stdout,
			Errfile: os.Stderr,
		}

		row := make([]float64, vindx.Vsize)
		switch elem.testcase {
		case "allones":
			for indx := 0; indx < vindx.Vsize; indx++ {
				row[indx] = float64(indx)
			}
		case "counting":
			for indx := 0; indx < vindx.Vsize; indx++ {
				row[indx] = float64(indx)
			}
		case "allzeros":
			// nothing to change
		default:
			t.Errorf("TestPrintModelRow: Unexpected test case '%s'\n", elem.testcase)
			continue
		}
		//fmt.Printf("Vsize: %d\n", vindx.Vsize)

		mychan := make(chan string)
		oldout, w, err := ms.RedirectModelSpecsLog(mychan)
		if err != nil {
			t.Errorf("RedirectStdout: %s\n", err)
			return // should this be continue?
		}
		fmt.Fprintf(ms.Logfile, "Row: %v\n", row)
		ms.printModelRow(row, elem.suppressNewline)

		str := ms.RestoreModelSpecsLog(mychan, oldout, w)
		c := make([]float64, len(row))
		A := make([][]float64, 0)
		A = append(A, row)
		b := make([]float64, 1)
		b[0] = 0.0
		if !checkPrintModelOutput(ms, A, b, c, str) {
			t.Errorf("PrintModelRow Case %d: did not match expected\n", i)
		}
	}
}

/*
	tests := []struct {
	}{
		{},
	}
	for i, elem := range tests {
	}

			ip: map[string]string{
				"setName":                    "",
				"filingStatus":               "",
				"key1":                       "retiree1",
				"key2":                       "retiree2",
				"eT_Age1":            "",
				"eT_Age2":            "",
				"eT_RetireAge1":      "",
				"eT_RetireAge2":      "",
				"eT_PlanThroughAge1": "",
				"eT_PlanThroughAge2": "",
				"eT_PIA1":            "",
				"eT_PIA2":            "",
				"eT_SS_Start1":       "",
				"eT_SS_Start2":       "",
				"eT_TDRA1":           "",
				"eT_TDRA2":           "",
				"eT_TDRA_Rate1":         "",
				"eT_TDRA_Rate2":         "",
				"eT_TDRA_Contrib1":   "",
				"eT_TDRA_Contrib2":   "",
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
			},
*/

func TestRedirectOutput(t *testing.T) {
	mychan := make(chan string)
	oldout, w, err := RedirectStdout(mychan, false)
	if err != nil {
		t.Errorf("RedirectStdout: %s\n", err)
		return
	}
	outstr := "This will be captured for comparisons later\nAnd this too\n"
	fmt.Printf("%s", outstr)
	str := RestoreStdout(mychan, oldout, w, false)
	if str != outstr {
		t.Errorf("Capured output fails: expected '%s', found '%s'", outstr, str)
	}
}

func RedirectStdout(mechan chan string, doNothing bool) (*os.File, *os.File, error) {
	if doNothing {
		return nil, nil, nil
	}
	oldStdout := os.Stdout
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		return os.Stdout, nil, err
	}
	os.Stdout = writePipe
	//mechan := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, readPipe)
		readPipe.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "func() copyPipe: %v\n", err)
			return
		}
		mechan <- buf.String()
	}()
	return oldStdout, writePipe, nil
}

func RestoreStdout(mechan chan string, oldStdout *os.File, writePipe *os.File, doNothing bool) string {
	if doNothing {
		return ""
	}
	// Reset the output again
	writePipe.Close()
	os.Stdout = oldStdout
	str := <-mechan
	return str
}

func (ms *ModelSpecs) RedirectModelSpecsTable(mechan chan string, DoNothing bool) (*os.File, *os.File, error) {
	if DoNothing {
		return nil, nil, nil
	}
	oldtable := ms.Ao.tableFile
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		return oldtable, nil, err
	}
	ms.Ao.tableFile = writePipe
	go func() {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, readPipe)
		readPipe.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "func() copyPipe: %v\n", err)
			return
		}
		mechan <- buf.String()
	}()
	return oldtable, writePipe, nil
}

func (ms *ModelSpecs) RestoreModelSpecsTable(mechan chan string, oldtable *os.File, writePipe *os.File, DoNothing bool) string {
	if DoNothing {
		return ""
	}
	// Reset the output again
	writePipe.Close()
	ms.Ao.tableFile = oldtable
	str := <-mechan
	return str
}

func (ms *ModelSpecs) RedirectModelSpecsLog(mechan chan string) (*os.File, *os.File, error) {
	oldlog := ms.Logfile
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		return oldlog, nil, err
	}
	ms.Logfile = writePipe
	go func() {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, readPipe)
		readPipe.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "func() copyPipe: %v\n", err)
			return
		}
		mechan <- buf.String()
	}()
	return oldlog, writePipe, nil
}

func (ms *ModelSpecs) RestoreModelSpecsLog(mechan chan string, oldlog *os.File, writePipe *os.File) string {
	// Reset the output again
	writePipe.Close()
	ms.Logfile = oldlog
	str := <-mechan
	return str
}

func stripWhitespace(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			// if the character is any white space, drop it
			return -1
		}
		// else keep it in the string
		return r
	}, str)
}
