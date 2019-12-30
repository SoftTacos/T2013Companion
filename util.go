package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"

	"github.com/go-yaml/yaml"
)

func ReadCharacterData(data []byte) {
	delimeter := []byte("---")
	dataSlices := bytes.Split(data, delimeter)
	characters = make(map[string]*Character) //make([]*Character, len(dataSlices))
	itemsTag := "#!items"
	//characterTag := "#!character"

	for _, charData := range dataSlices { //TODO: Break down into abstracted functions that read each section of character data to allow for flexible parsing of out of order data
		char := Character{}
		subSlices := bytes.Split(charData, []byte(itemsTag))
		//CHARACTER
		err := yaml.Unmarshal(subSlices[0], &char)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		//TODO: Validate character data
		characters[char.Name] = &char
		//ITEMS
		itemsSubset := struct {
			Items         []string
			CurrentWeapon string
		}{}

		err = yaml.Unmarshal(subSlices[1], &itemsSubset)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		char.Items = make(map[string]Item)
		for _, curItemName := range itemsSubset.Items {
			_, ok := itemData[curItemName]
			if !ok {
				fmt.Println("Error: Item does not exist in items list: \"", curItemName, "\"")
				continue
			}
			char.Items[curItemName] = itemData[curItemName]
		}
		_, ok := char.Items[itemsSubset.CurrentWeapon]
		if !ok {
			fmt.Println("Error: Character doesn't have item in their inventory: \"", itemsSubset.CurrentWeapon, "\"")
			continue
		}
		char.SetCurrentWeapon(itemsSubset.CurrentWeapon)

	}
	//TODO: REFACTOR THIS FOR CHAR MAP
	/*
		if characters[len(characters)-1].Name == "" {
			characters = characters[0 : len(characters)-1]
		}
	*/

}

func ReadRuleData(data []byte) {
	err := yaml.Unmarshal(data, &rules)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func ReadItemData(data []byte) {
	itemData := make(map[string]map[string]map[string]string) //type->names->data
	err := yaml.Unmarshal(data, &itemData)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	//fmt.Println(itemData)
	SetupRangedWeapons(itemData["rangedWeapons"])
	SetupMeleeWeapons(itemData["meleeWeapons"])
	SetupConsumables(itemData["consumables"])
	SetupMisc(itemData["misc"])
}

//TODO: Consuder switching to type-switch

//var itemData map[string]Item
func SetupRangedWeapons(weaponData map[string]map[string]string) {
	for wepName := range weaponData {
		currentWeaponData := weaponData[wepName]
		pSpeed, _ := strconv.ParseUint(currentWeaponData["speed"], 10, 64)
		pDamage, _ := strconv.ParseUint(currentWeaponData["damage"], 10, 64)
		pBulk, _ := strconv.ParseUint(currentWeaponData["bulk"], 10, 64)
		pWeight, _ := strconv.ParseFloat(currentWeaponData["weight"], 32)
		pDesc, _ := currentWeaponData["description"]

		wep := RangedWeaponItem{
			Name:        wepName,
			Speed:       uint8(pSpeed),
			Damage:      uint8(pDamage),
			Bulk:        uint8(pBulk),
			weight:      float32(pWeight),
			description: pDesc,
		}
		itemData[wepName] = &wep
	}
}

func SetupMeleeWeapons(weaponData map[string]map[string]string) {
	for wepName := range weaponData {
		currentWeaponData := weaponData[wepName]
		pSpeed, _ := strconv.ParseUint(currentWeaponData["speed"], 10, 64)
		pDamage, _ := strconv.ParseUint(currentWeaponData["damage"], 10, 64)
		pWeight, _ := strconv.ParseFloat(currentWeaponData["weight"], 32)
		pDesc, _ := currentWeaponData["description"]

		wep := MeleeWeaponItem{
			Name:        wepName,
			Speed:       uint8(pSpeed),
			Damage:      uint8(pDamage),
			weight:      float32(pWeight),
			description: pDesc,
		}
		itemData[wepName] = &wep
	}
}

func SetupConsumables(data map[string]map[string]string) {
	for name := range data {
		currentItemData := data[name]
		pAmount, _ := strconv.ParseFloat(currentItemData["amount"], 32)
		pWeight, _ := strconv.ParseFloat(currentItemData["weight"], 32)
		pDesc, _ := currentItemData["description"]

		item := ConsumableItem{
			Name:        name,
			Amount:      float32(pAmount),
			weight:      float32(pWeight),
			description: pDesc,
		}
		itemData[name] = &item
	}
}

func SetupMisc(data map[string]map[string]string) {
	for name := range data {
		currentItemData := data[name]
		pWeight, _ := strconv.ParseFloat(currentItemData["weight"], 32)
		pDesc, _ := currentItemData["description"]

		item := MiscItem{
			Name:        name,
			weight:      float32(pWeight),
			description: pDesc,
		}
		itemData[name] = &item
	}
}

func NumberMenu(max uint) uint {
	var validOption uint
	for validNumber := false; !validNumber; {
		input, _ := globals.reader.ReadString('\n')
		input = globals.whitespaceRegex.ReplaceAllString(input, "")
		option, err := strconv.ParseUint(input, 10, 64)
		if err != nil {
			fmt.Println(err)
			print("Input could not be read as a number, please provide a valid number\n")
			continue
		}
		if uint(option) > max {
			print("Input was too high\n")
			continue
		}
		validOption = uint(option)
		validNumber = true
	}
	return validOption
}
