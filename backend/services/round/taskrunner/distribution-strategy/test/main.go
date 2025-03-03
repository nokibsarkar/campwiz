package main

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"
)

type EV struct {
	EvaluationID string
	SubmissionID string
	JudgeID      string
}
type EVList []EV

var alreadyEvaluations = EVList{
	{"1", "s8", "j1"},
	{"2", "s7", "j2"},
	{"3", "s1", "j3"},
	{"4", "s6", "j1"},
	{"5", "s2", "j2"},
	{"6", "s2", "j3"},
	{"7", "s3", "j3"},
	{"8", "s4", "j3"},
	{"9", "s5", "j1"},
	{"10", "s5", "j3"},
}

func (e *EVList) Len() int {
	return len(*e)
}
func (e *EVList) Less(i, j int) bool {
	return strings.Compare((*e)[i].SubmissionID, (*e)[j].SubmissionID) < 0
}
func (e *EVList) Swap(i, j int) {
	(*e)[i], (*e)[j] = (*e)[j], (*e)[i]
}

func FindByJuryID(j string) []EV {
	var evs []EV
	for _, ev := range alreadyEvaluations {
		if ev.JudgeID == j {
			evs = append(evs, ev)
		}
	}
	return evs
}
func JurySet(i string) sets.Set[string] {
	evs := sets.Set[string]{}
	for _, ev := range alreadyEvaluations {
		if ev.JudgeID == i {
			evs.Insert(ev.SubmissionID)
		}
	}
	return evs
}
func do() {
	maxRound := 1
	juryIDs := []string{"j1", "j2", "j3"}
	juryMap := map[string]sets.Set[string]{}
	totalExistingWorkload := 0
	for _, jid := range juryIDs {
		juryMap[jid] = JurySet(jid)
		totalExistingWorkload += len(juryMap[jid])
	}
	fmt.Printf("Jury Map %v\n", juryMap)
	candidates := sets.NewString("s1", "s2", "s3", "s4")
	totalNewWorkload := len(candidates) * maxRound
	grandTotal := totalExistingWorkload + totalNewWorkload
	juryCount := len(juryIDs)
	grandAverage := grandTotal / juryCount

	fmt.Printf("Total : %d Avg : %d\n", grandTotal, grandAverage)
	extra := grandTotal % juryCount
	for range maxRound {
		// to be distributed
		c := candidates.Clone()
		for _, j := range juryIDs {
			alreadyAssigned := juryMap[j]
			currentWorkload := len(alreadyAssigned)
			newWorkload := max(0, grandAverage-currentWorkload)
			if extra > 0 {
				if currentWorkload+newWorkload <= grandAverage {
					newWorkload++
					extra--
				}
			}
			fmt.Printf("New workload for Jury %s : %d\n", j, newWorkload)
			for ca := range c {
				if alreadyAssigned.Has(ca) {
					fmt.Printf("Submission  %s is already assigned to %s\n", ca, j)
				} else {
					alreadyAssigned.Insert(ca)
					c.Delete(ca)
				}
			}
		}
		fmt.Println("Non assignable :", c)
	}
	for j, v := range juryMap {
		fmt.Printf("Jury %s has following submissions %v\n", j, v)
	}
}

func main() {
	do()
	// Assignments
	// ExportToCache
}
