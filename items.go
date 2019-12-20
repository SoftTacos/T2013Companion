package main

type Item interface {
	GetType() string
	GetWeight() float32
}

type MiscItem struct {
	Name   string
	weight float32
}

func (m *MiscItem) GetType() string    { return "misc" }
func (m *MiscItem) GetWeight() float32 { return m.weight }

type ConsumableItem struct {
	Name   string
	Amount float32
	weight float32
}

func (c *ConsumableItem) GetType() string    { return "consumable" }
func (c *ConsumableItem) GetWeight() float32 { return c.weight }

/*
type Weapon interface {
	Item

}
*/
//TODO: Make melee and ranged inherit from a weapon
type RangedWeaponItem struct {
	Name   string
	Speed  uint8
	Damage uint8
	Bulk   uint8
	weight float32
}

func (w *RangedWeaponItem) GetType() string    { return "rangedWeapon" }
func (w *RangedWeaponItem) GetWeight() float32 { return w.weight }

type MeleeWeaponItem struct {
	Name   string
	Speed  uint8
	Damage uint8
	weight float32
}

func (w *MeleeWeaponItem) GetType() string    { return "meleeWeapon" }
func (w *MeleeWeaponItem) GetWeight() float32 { return w.weight }
