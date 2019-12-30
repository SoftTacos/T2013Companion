package main

type Character struct {
	Name          string
	Stats         map[string]uint8 //name->level
	Skills        map[string]uint8 //name->level
	Items         map[string]Item
	CurrentWeapon *RangedWeaponItem //placeholder, will be removed and replaced with an item specific call
}

func (c *Character) Init() {
	c.Name = "BLANK"
	c.Stats = map[string]uint8{
		"AWA":  0,
		"CDN":  0,
		"FIT":  0,
		"MUS":  0,
		"COG":  0,
		"EDU":  0,
		"PER":  0,
		"RES":  0,
		"CUF":  0, //DONT USE
		"OODA": 0, //DONT USE
	}
	//c.Gear = make(map[string]float32)
}

func (c *Character) InitiativeCheck() uint8 {
	encIni := rules.EncMap[c.Encumbrance()]
	roll := advantage(nd20(2))

	//initiative is 2d20 VS OODA, OODA is TN
	checkMargin := (int(c.Stats["AWA"]) - int(roll))
	if checkMargin < 0 {
		return encIni
	}
	return encIni + uint8(checkMargin)*2
}

func (c *Character) Encumbrance() string {
	return ""
}

func (c *Character) SetCurrentWeapon(name string) bool {
	_, ok := c.Items[name]
	if !ok {
		return false
	}
	c.CurrentWeapon = c.Items[name].(*RangedWeaponItem)
	return true
}
