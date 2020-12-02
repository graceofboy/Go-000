package model

type Persion struct {
	Name string
	Age  string
}

func (p *Persion) TableName() string {
	return "persion"
}
