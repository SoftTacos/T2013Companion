package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"

	"github.com/go-yaml/yaml"
)

func ReadCharacterData(data []byte) []*Character {
	delimeter := []byte("---")
	dataSlices := bytes.Split(data, delimeter)
	characters := make([]*Character, len(dataSlices))
	itemsTag := "#!items"
	//characterTag := "#!character"

	for i, charData := range dataSlices { //TODO: Break down into abstracted functions that read each section of character data to allow for flexible parsing of out of order data
		char := Character{}
		subSlices := bytes.Split(charData, []byte(itemsTag))
		//CHARACTER
		err := yaml.Unmarshal(subSlices[0], &char)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		//TODO: Validate character data
		characters[i] = &char
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
			//TODO: validate if the item exists in item list
			char.Items[curItemName] = itemData[curItemName]
		}
		//TODO: Check if character even HAS the current weapon
		char.SetCurrentWeapon(itemsSubset.CurrentWeapon)

	}
	if characters[len(characters)-1].Name == "" {
		characters = characters[0 : len(characters)-1]
	}
	return characters
}

func ReadItemData(data []byte) {
	itemData := make(map[string]map[string]map[string]string) //type->names->data
	err := yaml.Unmarshal(data, &itemData)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	//fmt.Println(itemData)
	SetupRangedWeapons(itemData["rangedWeapons"])
	//SetupMeleeWeapons
	//SetupConsumables
	//SetupMisc
}

//var itemData map[string]Item
func SetupRangedWeapons(weaponData map[string]map[string]string) {
	for wepName := range weaponData {
		currentWeaponData := weaponData[wepName]
		pSpeed, _ := strconv.ParseUint(currentWeaponData["speed"], 10, 64)
		pDamage, _ := strconv.ParseUint(currentWeaponData["damage"], 10, 64)
		pBulk, _ := strconv.ParseUint(currentWeaponData["bulk"], 10, 64)
		pWeight, _ := strconv.ParseFloat(currentWeaponData["weight"], 32)

		wep := RangedWeaponItem{
			Name:   wepName,
			Speed:  uint8(pSpeed),
			Damage: uint8(pDamage),
			Bulk:   uint8(pBulk),
			weight: float32(pWeight),
		}
		fmt.Println(wepName, wep)
		itemData[wepName] = &wep
	}
}

/*
func ReadWeapons(data []byte) {
	delimeter := []byte("---")
	dataSlices := bytes.Split(data, delimeter)
	weapons = make([]RangedWeaponItem, len(dataSlices))

	for i, wepData := range dataSlices {
		weapon := RangedWeaponItem{}
		err := yaml.Unmarshal(wepData, &weapon)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		weapons[i] = weapon
	}
	if weapons[len(weapons)-1].Name == "" {
		weapons = weapons[0 : len(weapons)-1]
	}
}
*/
