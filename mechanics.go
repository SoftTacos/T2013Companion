package main

import (
	"container/list"
	"fmt"
	"sort"
)

func SkillCheck(statName string, skillName string, char *Character, difficulty int) int {
	mos := 0 //MOS = margin of success

	return mos
}

func StatCheck(statName string, char *Character, difficulty int) int {
	//MOS = margin of success
	//2d20 VS TN+Difficulty; TN=target number; you want to be below the TN+DIFF
	roll := advantage(nd20(2)) //all stat checks are 2d20L
	//roll <= TN+diff
	TN := int(char.Stats[statName]) + difficulty
	//mos := TN - roll
	return TN - int(roll)
}

func advantage(rolls []uint8) uint8 {
	lowest := rolls[0]
	for i, _ := range rolls {
		if lowest > rolls[i] {
			lowest = rolls[i]
		}
	}
	return lowest
}

func nd20(n uint8) []uint8 {
	rolls := make([]uint8, n)
	for i, _ := range rolls {
		rolls[i] = d20()
	}
	return rolls
}

func d20() uint8 {
	return uint8((randy.Uint64() % 20) + 1)
}

func Reorder(inits *list.List) {
	init := inits.Front().Value.(*Turn).Init
	for e := inits.Back(); e != nil; e = e.Prev() {
		turn := e.Value.(*Turn)

		if init < turn.Init {
			fmt.Println("REORDERING")
			inits.InsertAfter(inits.Remove(inits.Front()), e)
			break
		}
	}
}

//combat is rounds of EoF, EoF is a series of turns that happen until initiatives all go to 0
//Combat() just does one round of EoF
func Combat(inits *list.List) {
	turn := inits.Front().Value.(*Turn)
	for { //roundOver := false; !roundOver; {
		//turn is a collection of ticks of the same number
		//TAKE TURN
		TakeTurn(turn)
		Reorder(inits)
		//decrement Turn.Init value, move the item within the linked list
		//look at the front of the LL, is it the same as turn? Then do that
		turn = inits.Front().Value.(*Turn)
		if turn.Init <= 0 {
			break
		}
	}
}

//todo: will want to just generate the list and sort it later, mostly as an exercise in LinkedList sorting
func GenerateInitiatives(chars map[string]*Character) *list.List {
	initiatives := []*Turn{}
	for _, value := range chars {
		initiatives = append(initiatives,
			&Turn{
				Init: int(value.InitiativeCheck()),
				Char: value,
			})
	}
	sort.Slice(initiatives, func(i, j int) bool { return initiatives[i].Init > initiatives[j].Init })
	initList := list.New()
	for i, _ := range initiatives {
		initList.PushBack(initiatives[i])
	}
	return initList
}
