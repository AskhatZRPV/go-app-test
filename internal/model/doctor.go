package model

type Doctor struct {
	ID   string
	Name string
	Spec string
	Slot []SlotModel
}

type SlotModel struct {
	Id   string
	Time string
}
