package main

type Item interface {
	GetName() string
	GetDescription() string
	GetType() string
	GetWeight() float32
}

type MiscItem struct {
	Name        string
	weight      float32
	description string
}

func (m *MiscItem) GetName() string        { return m.Name }
func (m *MiscItem) GetType() string        { return "misc" }
func (m *MiscItem) GetDescription() string { return m.description }
func (m *MiscItem) GetWeight() float32     { return m.weight }

type ConsumableItem struct {
	Name        string
	Amount      float32
	weight      float32
	description string
}

func (c *ConsumableItem) GetName() string        { return c.Name }
func (c *ConsumableItem) GetType() string        { return "consumable" }
func (c *ConsumableItem) GetDescription() string { return c.description }
func (c *ConsumableItem) GetWeight() float32     { return c.weight }

/*
type Weapon interface {
	Item

}
*/
//TODO: Make melee and ranged inherit from a weapon
type RangedWeaponItem struct {
	Name        string
	Speed       uint8
	Damage      uint8
	Bulk        uint8
	weight      float32
	description string
}

func (w *RangedWeaponItem) GetName() string        { return w.Name }
func (w *RangedWeaponItem) GetType() string        { return "rangedWeapon" }
func (w *RangedWeaponItem) GetDescription() string { return w.description }
func (w *RangedWeaponItem) GetWeight() float32     { return w.weight }

type MeleeWeaponItem struct {
	Name        string
	Speed       uint8
	Damage      uint8
	weight      float32
	description string
}

func (w *MeleeWeaponItem) GetName() string        { return w.Name }
func (w *MeleeWeaponItem) GetType() string        { return "meleeWeapon" }
func (w *MeleeWeaponItem) GetDescription() string { return w.description }
func (w *MeleeWeaponItem) GetWeight() float32     { return w.weight }
