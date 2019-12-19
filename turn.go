package main

import "fmt"

type Turn struct {
	Init int
	Char *Character
}

func TakeTurn(turn *Turn) {
	//fmt.Printf("\n\t%s's turn:%d\n", turn.Char.Name, turn.Init)
	validRange := uint(5) //lazy!
	fmt.Printf("\n\t%s's turn:%d\n1:Attack\n2:Move\n3:Change Stance\n4:Communicate\n5:Reload\n", turn.Char.Name, turn.Init)
	choice := NumberMenu(validRange)
	rules.TurnActions[choice](turn)
}

func Turn_Attack(turn *Turn) {
	fmt.Println("ATTACKING")
	if turn.Char.CurrentWeapon == nil {
		fmt.Println("You can't attack when you don't have a weapon!")
		return
	}
	turn.Init -= int(turn.Char.CurrentWeapon.Speed)
}

func Turn_ChangeStance(turn *Turn) {
	fmt.Println("STANCE")
	turn.Init -= 2 //stance changes are static cost of 2
}

func Turn_Communicate(turn *Turn) {
	fmt.Println("COMMUNICATE")
	turn.Init -= int(NumberMenu(20))
}

func Turn_Reload(turn *Turn) {
	fmt.Println("RELOAD")
	if turn.Char.CurrentWeapon == nil {
		fmt.Println("You can't reload when you don't have a weapon!")
		return
	}
	turn.Init -= int(turn.Char.CurrentWeapon.Bulk)
}

func Turn_Move(turn *Turn) {
	fmt.Println("MOVE")

}
