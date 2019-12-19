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

//TODO: Make melee and ranged interit from a weapon
type RangedWeaponItem struct {
	Name   string
	Speed  uint8
	Damage uint8
	Bulk   uint8
	weight float32
}

func (w *RangedWeaponItem) GetType() string    { return "weapon" }
func (w *RangedWeaponItem) GetWeight() float32 { return w.weight }
