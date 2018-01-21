package rplanlib

import "fmt"
import "math"
import "strconv"

type account struct {
    bal float64
    estateTax float64
    contributions []float64
    rrate float64
    acctype string
    mykey string
}
type modelSpecs struct {
    InputParams map[string]string
    vindx VectorVarIndex

    ti Taxinfo
    // The following was through 'S'
    numyr int
    maximize string // "Spending" or "PlusEstate"
    accounttable []account
    accmap map[string]int
    retiree []map[string]string

    income []float64
    SS []float64
    expenses []float64
    taxed []float64
    asset_sale []float64
    cg_asset_taxed []float64

    iRate float64
    rRate float64

    min float64
    max float64

    preplanyears int


} 

func intMax(a, b int) int {
    if a>b {
        return a
    } else {
        return b
    }
}
func intMin(a, b int) int {
    if a<b {
        return a
    } else {
        return b
    }
}
func checkStrconvError(e error) {
    if e != nil {
        panic(e)
    }
}

// mergeVectors sums two vectors of equal lenth returning a third vector
func mergeVectors(v1, v2 []float64) ([]float64, error) {
    if len(v1) != len(v2) {
        err := fmt.Errorf("mergeVectors: Can not merge, lengths do not match, %d vs %d", len(v1), len(v2))
        return nil, err
    }
    v3 := make([]float64, len(v1))
    for i:=0 ; i<len(v1) ; i++ {
        v3[i] = v1[i] + v2[i]
    }
    return v3, nil
}

// buildVector creates a vector with a 'rate' adjusted 'yearly' amount in buckets between start and end age 
func buildVector(yearly, startAge, endAge, vecStartAge, vecEndAge int, rate float64, baseAge int) ([]float64, error) {
    //verify that startAge and endAge are within vecStart and end
    if vecStartAge > vecEndAge {
        err := fmt.Errorf("vec start age (%d) is greater than vec end age (%d)", vecStartAge, vecEndAge)
        return nil, err
    }
    if startAge < vecStartAge { 
        startAge = vecStartAge
    }
    if startAge > vecEndAge {
        startAge = vecEndAge
    }
    if endAge < vecStartAge {
        endAge = vecStartAge
    }
    if endAge > vecEndAge {
        endAge = vecEndAge
    }
    if startAge > endAge {
        err := fmt.Errorf("start age (%d) is greater than end age (%d)", startAge, endAge)
        return nil, err
    }
    
    vecSize := vecEndAge - vecStartAge
    vec := make([]float64, vecSize)
    for i:=0; i<vecSize; i++ {
        if i >= startAge - vecStartAge && i <= endAge - vecStartAge {
            to := float64(startAge-baseAge+i)
            adj := math.Pow(rate,to)
            vec[i] = float64(yearly) * adj // or something like this FIXME TODO
        }
    }
    return vec, nil
}
func NewModelSpecs(vindx VectorVarIndex, 
                   ti Taxinfo, 
                   ip map[string]string,
                   verbose bool, 
                   no_TDRA_ROTHRA_DEPOSITS bool) modelSpecs {
                       var err error
    ms := modelSpecs{
        InputParams: ip,
        vindx: vindx,
        ti: ti,
        maximize: "Spending", // or "PlusEstate"
        iRate: 1.025,
        rRate: 1.06,
        min: -1,
        max: -1,
    }
    /*
    ** TODO what if any of the strings returned by ip[] are empty? FIXME
    */
    if ip["filingStatus"] == "joint" {
    age1, err := strconv.Atoi(ip["eT_Age1"])
    checkStrconvError(err) 
    age2, err := strconv.Atoi(ip["eT_Age2"])
    checkStrconvError(err) 
    retireAge1, err := strconv.Atoi(ip["eT_RetireAge1"])
    checkStrconvError(err) 
    if retireAge1 < age1 {
        retireAge1 = age1
    }
    retireAge2, err := strconv.Atoi(ip["eT_RetireAge2"])
    checkStrconvError(err) 
    if retireAge2 < age2 {
        retireAge2 = age2
    }
    planThroughAge1, err := strconv.Atoi(ip["eT_PlanThroughAge1"])
    checkStrconvError(err) 
    planThroughAge2, err := strconv.Atoi(ip["eT_PlanThroughAge2"])
    checkStrconvError(err) 
    yearsToRetire1 := retireAge1 - age1
    yearsToRetire2 := retireAge2 - age2
    yearsToRetire := intMin(yearsToRetire1, yearsToRetire2)
    startPlan := yearsToRetire + age1
    through1 := planThroughAge1 - age1
    through2 := planThroughAge2 - age2
    endPlan := intMax(through1, through2) + 1 + age1
    delta := age1 - age2
    //ms.numyr := endPlan - startPlan
    //accounttable: []map[string]string
    //accmap: map[string]int
    //retiree: []map[string]string
    //income: []float64
    pia1, err := strconv.Atoi(ip["eT_PIA1"])
    checkStrconvError(err) 
    ssStart1, err := strconv.Atoi(ip["eT_SS_Start1"])
    checkStrconvError(err) 
    SS1, err := buildVector(pia1, ssStart1, endPlan, startPlan, endPlan, ms.iRate, age1)
    checkStrconvError(err) 
    pia2, err := strconv.Atoi(ip["eT_PIA2"])
    checkStrconvError(err) 
    ssStart2, err := strconv.Atoi(ip["eT_SS_Start2"])
    checkStrconvError(err) 
    SS2, err2 := buildVector(pia2, ssStart2, endPlan, startPlan, endPlan, ms.iRate, age1)
    checkStrconvError(err) 
    SS, err := mergeVectors(SS1,SS2)
    checkStrconvError(err) 
    //expenses: []float64
    //taxed: []float64 // maybe have a special function that adds vec to tax vec?
    //asset_sale: []float64
    //cg_asset_taxed: []float64
    } else { // single or mseparate
    age1, err := strconv.Atoi(ip["eT_Age1"])
    checkStrconvError(err) 
    retireAge1, err := strconv.Atoi(ip["eT_RetireAge1"])
    checkStrconvError(err) 
    if retireAge1 < age1 {
        retireAge1 = age1
    }
    planThroughAge1, err := strconv.Atoi(ip["eT_PlanThroughAge1"])
    checkStrconvError(err) 
    yearsToRetire1 := retireAge1 - age1
    yearsToRetire := yearsToRetire1
    startPlan := yearsToRetire + age1
    through1 := planThroughAge1 - age1
    endPlan := through1 + 1 + age1
    //delta := age1 - age2
    //ms.numyr := endPlan - startPlan
    //accounttable: []map[string]string
    //accmap: map[string]int
    //retiree: []map[string]string
    //income: []float64
    //SS: []float64
    //expenses: []float64
    //taxed: []float64
    //asset_sale: []float64
    //cg_asset_taxed: []float64

    }
    return ms
}
/*
class lp_constraint_model:
    def __init__(verbose, no_TDRA_ROTHRA_DEPOSITS):
        self.S = S
        self.verbose = verbose
        self.noTdraRothraDeposits = no_TDRA_ROTHRA_DEPOSITS

*/

// BuildModel for:
// Minimize: c^T * x
// Subject to: A_ub * x <= b_ub
// all vars positive
func (ms modelSpecs) BuildModel() ([]float64, [][]float64, []float64) {
    
        // TODO integrate the following assignments into the code and remove them
        //S = ms.S
    
        nvars := ms.vindx.Vsize
        A := make([][]float64, 0)
        b := make([]float64, 0)
        c := make([]float64, nvars)
    
        //
        // Add objective function (S1') becomes (R1') if PlusEstate is added
        //
        for year := 0; year<ms.numyr; year++ {
            c[ms.vindx.S(year)] = -1
        }
        //
        // Add objective function tax bracket forcing function (EXPERIMENTAL)
        //
        for year := 0; year<ms.numyr; year++ {
            for k := 0; k<len(*ms.ti.Taxtable); k++ {
                // multiplies the impact of higher brackets opposite to
                // optimization the intent here is to pressure higher 
                // brackets more and pack the lower brackets
                c[ms.vindx.X(year,k)] = float64(k)/10 
            }
        }
        //
        // Adder objective function (R1') when PlusEstate is added
        //
        if ms.maximize == "PlusEstate" {
            for j := 0; j<len(ms.accounttable); j++ {
                c[ms.vindx.B(ms.numyr,j)] = -1*ms.accounttable[j].estateTax // account discount rate
            }
            print("\nConstructing Spending + Estate Model:\n")
        } else {
            print("\nConstructing Spending Model:\n")
            startamount := 0.0
            for j := 0; j<len(ms.accounttable); j++ {
                startamount += ms.accounttable[j].bal
            }
            balancer := 1.0/(startamount) 
            for j := 0; j<len(ms.accounttable); j++ {
                c[ms.vindx.B(ms.numyr,j)] = -1*balancer *ms.accounttable[j].estateTax // balance and discount rate
            }
        }
        //
        // Add constraint (2')
        //
        for year := 0; year<ms.numyr; year++ {
            row := make([]float64, nvars)
            for j := 0; j<len(ms.accounttable); j++ {
                p := 1.0
                if ms.accounttable[j].acctype != "aftertax" {
                    if ms.ti.apply_early_penalty(year,ms.accounttable[j].mykey) {
                        p = 1-ms.ti.Penalty
                    }
                }
                row[ms.vindx.W(year,j)] = -1*p 
            }
            for k:=0;  k<len(*ms.ti.Taxtable);k++ {
                row[ms.vindx.X(year,k)] = (*ms.ti.Taxtable)[k][2] // income tax
            } 
            if ms.accmap["aftertax"] > 0 {
                for l:=0; l<len(*ms.ti.Capgainstable); l++ {
                    row[ms.vindx.Y(year,l)] = (*ms.ti.Capgainstable)[l][2] // cap gains tax
                } 
                for j := 0; j<len(ms.accounttable); j++ {
                    row[ms.vindx.D(year,j)] = 1
                }
            }
            row[ms.vindx.S(year)] = 1
            A = append(A, row)
            b= append(b, ms.income[year] + ms.SS[year] - ms.expenses[year])
        }
        //
        // Add constraint (3a')
        //
        for year := 0; year<ms.numyr-1; year++ {
            row := make([]float64, nvars)
            row[ms.vindx.S(year+1)] = 1
            row[ms.vindx.S(year)] = -1*ms.iRate
            A = append(A, row)
            b = append(b, 0)
        }
        //
        // Add constraint (3b')
        //
        for year := 0; year<ms.numyr-1; year++ {
            row := make([]float64, nvars)
            row[ms.vindx.S(year)] = ms.iRate
            row[ms.vindx.S(year+1)] = -1
            A = append(A, row)
            b = append(b, 0)
        }
        //
        // Add constrant (4') rows - not needed if [desired.income] is not defined in input
        //
        if ms.min != 0{
            for year:=0; year<1; year++ { // Only needs setting at the beginning
                row := make([]float64, nvars)
                row[ms.vindx.S(year)] = -1
                A = append(A, row)
                b = append(b, - ms.min )     // [- d_i]
            }
        }
    
        //
        // Add constraints for (5') rows - not added if [max.income] is 
        // not defined in input
        //
        if ms.max != 0{
            for year:=0; year<1; year++ { // Only needs to be set at the beginning
                row := make([]float64, nvars)
                row[ms.vindx.S(year)] = 1
                A = append(A, row)
                b = append(b, ms.max )     // [ dm_i]
            }
        }
    
        //
        // Add constaints for (6') rows
        //
        for year := 0; year<ms.numyr; year++ {
            row := make([]float64, nvars)
            for j := 0; j<len(ms.accounttable); j++ {
                if ms.accounttable[j].acctype != "aftertax"{
                    row[ms.vindx.D(year,j)] = 1
                } 
            }
            A = append(A, row)
            //b+=[min(ms.income[year],ms.maxContribution(year,None))] 
            // using ms.taxed rather than ms.income because income could
            // include non-taxed anueities that don't count.
            None := ""
            b = append(b, math.Min(ms.taxed[year], ms.ti.maxContribution(year,None)))
        }
        //
        // Add constaints for (7') rows
        //
        for year := 0; year<ms.numyr; year++ {
            // TODO this is not needed when there is only one retiree
            for _, v := range ms.retiree {
                row := make([]float64, nvars)
                for j := 0; j<len(ms.accounttable); j++ {
                    if v["mykey"] == ms.accounttable[j].mykey {
                        // ["acctype"] != "aftertax": no "mykey" in aftertax
                        // (this will either break or just not match - we 
                        // will see)
                        row[ms.vindx.D(year,j)] = 1
                    } 
                }
                A = append(A, row)
                b = append(b, S.maxContribution(year,v["mykey"]))
            }
        } 
        //
        // Add constaints for (8') rows
        //
        for year := 0; year<ms.numyr; year++ {
            for j := 0; j<len(ms.accounttable); j++ {
                v := ms.accounttable[j].contributions
                if v != nil {
                    if v[year] > 0{
                        row := make([]float64, nvars)
                        row[ms.vindx.D(year,j)] = -1
                        A = append(A, row)
                        b = append(b, -1*v[year])
                    }
                } 
            }
        }
        //
        // Add constaints for (9') rows
        //
        for year := 0; year<ms.numyr; year++ {
            for j := 0; j<intMin(2, len(ms.accounttable)); j++ {
                // at most the first two accounts are type IRA w/ 
                // RMD requirement
                if ms.accounttable[j].acctype == "IRA"{
                    ownerage := S.account_owner_age(year, ms.accounttable[j])
                    if ownerage >= 70 {
                        row := make([]float64, nvars)
                        row[ms.vindx.D(year,j)] = 1
                        A = append(A, row)
                        b = append(b, 0)
                    }
                }
            } 
        }
        //
        // Add constaints for (N') rows
        //
        if ms.noTdraRothraDeposits{
            for year:=0; year<ms.numyr; year++ {
                for j := 0; j<len(ms.accounttable); j++ {
                    v := ms.accounttable[j].contributions
                    max := 0.0
                    if v != nil {
                        max = v[year]
                    } 
                    if ms.accounttable[j].acctype != "aftertax"{
                        row := make([]float64, nvars)
                        row[ms.vindx.D(year,j)] = 1
                        A = append(A, row)
                        b = append(b, max)
                    }
                }
            }
        }
        //
        // Add constaints for (10') rows
        //
        for year := 0; year<ms.numyr; year++ {
            for j := 0; j<intMin(2, len(ms.accounttable)); j++ {
                // at most the first two accounts are type IRA 
                // w/ RMD requirement
                if ms.accounttable[j].acctype == "IRA"{
                    rmd := S.rmd_needed(year,ms.accounttable[j].mykey)
                    if rmd > 0{
                        row := make([]float64, nvars)
                        row[ms.vindx.B(year,j)] = 1/rmd 
                        row[ms.vindx.W(year,j)] = -1
                        A = append(A, row)
                        b = append(b, 0)
                    }
                }
            } 
        }
    
        //
        // Add constraints for (11')
        //
        for year := 0; year<ms.numyr; year++ {
            adj_inf := ms.i_rate**(ms.preplanyears+year)
            row := make([]float64, nvars)
            for j := 0; j<intMin(2, len(ms.accounttable)); j++ {
                // IRA can only be in the first two accounts
                if ms.accounttable[j].acctype == "IRA"{
                    row[ms.vindx.W(year,j)] = 1 // Account 0 is TDRA
                    row[ms.vindx.D(year,j)] = -1 // Account 0 is TDRA
                }
            } 
            for k:=0;  k<len(ms.ti.Taxtable);k++ {
                row[ms.vindx.X(year,k)] = -1
            }
            A = append(A, row)
            b = append(b, ms.ti.stded*adj_inf-ms.taxed[year]-ms.ti.SStaxable*ms.SS[year])
        }
        //
        // Add constraints for (12')
        //
        for year := 0; year<ms.numyr; year++ {
            for k:=0;  k<len(ms.ti.Taxtable)-1;k++ {
                row := make([]float64, nvars)
                row[ms.vindx.X(year,k)] = 1
                A = append(A, row)
                b = append(b, (ms.ti.Taxtable[k][1])*(ms.i_rate**(ms.preplanyears+year))) // inflation adjusted
            }
        }
        //
        // Add constraints for (13a')
        //
        if ms.accmap["aftertax"] > 0{
            for year:=0; year<ms.numyr; year++ {
                f = ms.cg_taxable_fraction(year)
                row := make([]float64, nvars)
                for l:=0; l<len(ms.ti.Capgainstable);l++{
                    row[ms.vindx.Y(year,l)] = 1
                }
                // Awful Hack! If year of asset sale, assume w(i,j)-D(i,j) is 
                // negative so taxable from this is zero
                if ms.cg_asset_taxed[year] <= 0{ // i.e., no sale
                    j = len(ms.accounttable)-1 // last Acc is investment / stocks
                    row[ms.vindx.W(year,j)] = -1*f 
                    row[ms.vindx.D(year,j)] = f 
                }
                A = append(A, row)
                b = append(b, ms.cg_asset_taxed[year])
            }
        }
        //
        // Add constraints for (13b')
        //
        if ms.accmap["aftertax"] > 0 {
            for year:=0; year<ms.numyr; year++ {
                f = ms.cg_taxable_fraction(year)
                row := make([]float64, nvars)
                ////// Awful Hack! If year of asset sale, assume w(i,j)-D(i,j) is 
                ////// negative so taxable from this is zero
                if ms.cg_asset_taxed[year] <= 0{ // i.e., no sale
                    j = len(ms.accounttable)-1 // last Acc is investment / stocks
                    row[ms.vindx.W(year,j)] = f 
                    row[ms.vindx.D(year,j)] = -f 
                }
                for l:=0; l<len(ms.ti.Capgainstable);l++{
                    row[ms.vindx.Y(year,l)] = -1
                }
                A = append(A, row)
                b = append(b, -ms.cg_asset_taxed[year])
            }
        }
        //
        // Add constraints for (14')
        //
        if ms.accmap["aftertax"] > 0{
            for year:=0; year<ms.numyr; year++ {
                adj_inf = ms.i_rate**(ms.preplanyears+year)
                for l:=0; l<len(ms.ti.Capgainstable)-1;l++{
                    row := make([]float64, nvars)
                    row[ms.vindx.Y(year,l)] = 1
                    for k:=0;  k<len(ms.ti.Taxtable)-1;k++ {
                        if ms.ti.Taxtable[k][0] >= (*ms.ti.Capgainstable)[l][0] && ms.ti.Taxtable[k][0] < (*ms.ti.Capgainstable)[l+1][0] {
                            row[ms.vindx.X(year,k)] = 1
                        }
                    }
                    A = append(A, row)
                    b = append(b, ms.ti.Capgainstable[l][1]*adj_inf) // mcg[i,l] inflation adjusted
                    //print_constraint( row, ms.ti.Capgainstable[l][1]*adj_inf)
                }
            }
        }
        //
        // Add constraints for (15a')
        //
        for year := 0; year<ms.numyr; year++ {
            for j := 0; j<len(ms.accounttable); j++ {
                //j = len(ms.accounttable)-1 // nl the last account, the investment account
                row := make([]float64, nvars)
                row[ms.vindx.B(year+1,j)] = 1 // b[i,j] supports an extra year
                row[ms.vindx.B(year,j)] = -1*ms.accounttable[j].rate
                row[ms.vindx.W(year,j)] = ms.accounttable[j].rate
                row[ms.vindx.D(year,j)] = -1*ms.accounttable[j].rate
                A = append(A, row)
                // In the event of a sell of an asset for the year 
                temp = 0
                if ms.accounttable[j].acctype == "aftertax"{
                    temp= ms.asset_sale[year] * ms.accounttable[j].rate //TODO test
                }
                b = append(b, temp)
                //print("temp_a: ", temp, "rate", ms.accounttable[j].rate , "asset sell price: ", ms.asset_sale[year]  )
            }
        } 
        //
        // Add constraints for (15b')
        //
        for year := 0; year<ms.numyr; year++ {
            for j := 0; j<len(ms.accounttable); j++ {
                //j = len(ms.accounttable)-1 // nl the last account, the investment account
                row := make([]float64, nvars)
                row[ms.vindx.B(year,j)] = ms.accounttable[j].rate
                row[ms.vindx.W(year,j)] = -1*ms.accounttable[j].rate
                row[ms.vindx.D(year,j)] = ms.accounttable[j].rate
                row[ms.vindx.B(year+1,j)] = -1  ////// b[i,j] supports an extra year
                A = append(A, row)
                temp = 0
                if ms.accounttable[j].acctype == "aftertax"{
                    temp= -1 * ms.asset_sale[year] * ms.accounttable[j].rate//TODO test
                }
                b = append(b, temp)
                //print("temp_b: ", temp, "rate", ms.accounttable[j].rate , "asset sell price: ", ms.asset_sale[year]  )
            }
        }
        //
        // Constraint for (16a')
        //   Set the begining b[1,j] balances
        //
        for j := 0; j<len(ms.accounttable); j++ {
            row := make([]float64, nvars)
            row[ms.vindx.B(0,j)] = 1
            A = append(A, row)
            b = append(b, ms.accounttable[j].bal)
        }
        //
        // Constraint for (16b')
        //   Set the begining b[1,j] balances
        //
        for j := 0; j<len(ms.accounttable); j++ {
            row := make([]float64, nvars)
            row[ms.vindx.B(0,j)] = -1
            A = append(A, row)
            b = append(b, -1*ms.accounttable[j].bal)
        }
        //
        // Constrant for (17') is default for sycpy so no code is needed
        //
        if ms.verbose{
            print("Num vars: ", len(c))
            print("Num contraints: ", len(b))
            print()
        }
    
        return c, A, b
}
    
    
func (ms modelSpecs) cg_taxable_fraction(year int) float64 {
    f = 1
    if ms.accmap["aftertax"] > 0{
        for _, v := range ms.accounttable {
            if v["acctype"] == "aftertax" {
                if v["bal"] > 0 {
                    f = 1 - (v["basis"]/(v["bal"]*v["rate"]**year))
                }
                break // should be the last entry anyway but...
            }
        }
    }
    return f
}
    
    
func (ms modelSpecs) print_model_matrix(c []float64, A [][]float64, b []float64, s, non_binding_only bool){
        if ! non_binding_only {
            fmt.Printf("c: ")
            ms.print_model_row(c)
            fmt.Printf("\n")
            fmt.Printf("B? i: A_ub[i]: b[i]")
            for constraint:=0; constraint<len(A); constraint++ {
                if s == nil || s[constraint] > 0 {
                    fmt.Printf("  ")
                } else{
                    fmt.Printf("B ")
                }
                fmt.Printf("%d: ", constraint)
                ms.print_constraint( A[constraint], b[constraint])
            }
        } else {
            fmt.Printf(" i: A_ub[i]: b[i]")
            j = 0
            for constraint:=0;  contraint<len(A); constraint++ {
                if s[constraint] >0 {
                    j+=1
                    fmt.Printf("%d: ", constraint)
                    ms.print_constraint( A[constraint], b[constraint])
                }
            }
            fmt.Printf("\n\n%d non-binding constrains printed\n", j)
        }
        fmt.Printf("\n")
}
    
    
func (ms modelSpecs) print_constraint(row []float64, b float64) {
        ms.print_model_row(row, true)
        fmt.Printf("<= b[]: %6.2f", b)
}
    
func (ms modelSpecs) print_model_row(row []float64, suppress_newline bool) {
        S = ms.S
    
        for i:=0; i<ms.numyr; i++ {
            for k:=0;  k<len(ms.ti.Taxtable);k++ {
                if row[ms.vindx.X(i, k)] != 0{
                    fmt.Printf("x[%d,%d]: %6.3f", i, k, row[ms.vindx.X(i, k)])
                }
            }
        }
        if ms.accmap["aftertax"] > 0 {
            for i:=0; i<ms.numyr; i++ {
                for l:=0; l<len(ms.ti.Capgainstable);l++{
                    if row[ms.vindx.Y(i, l)] != 0{
                        fmt.Printf("y[%d,%d]: %6.3f ", i, l, row[ms.vindx.Y(i, l)])
                    }
                }
            }
        }
        for i:=0; i<ms.numyr; i++ {
            for j := 0; j<len(ms.accounttable); j++ {
                if row[ms.vindx.W(i, j)] != 0{
                    fmt.Printf("w[%d,%d]: %6.3f ", i, j, row[ms.vindx.W(i, j)])
                }
            }
        }
        for i:=0; i<ms.numyr+1; i++ { // b[] has an extra year
            for j := 0; j<len(ms.accounttable); j++ {
                if row[ms.vindx.B(i, j)] != 0{
                    fmt.Printf("b[%d,%d]: %6.3f ", i, j, row[ms.vindx.B(i, j)])
                }
            }
        }
        for i:=0; i<ms.numyr; i++ {
            if row[ms.vindx.S(i)] !=0{
                fmt.Printf("s[%d]: %6.3f ", i, row[ms.vindx.S(i)])
            }
        }
        if ms.accmap["aftertax"] > 0 {
            for i:=0; i<ms.numyr; i++ {
                for j := 0; j<len(ms.accounttable); j++ {
                    if row[ms.vindx.D(i,j)] !=0{
                        fmt.Printf("D[%d,%d]: %6.3f ", i, j, row[ms.vindx.D(i,j)])
                    }
                }
            }
        }
        if ! suppress_newline {
            fmt.Printf("\n")
        }
}
    