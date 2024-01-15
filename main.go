package main

import (
	"errors"
	"os"
	"strings"
)

type Register struct {
	Id   string
	Type string      // "M" | "T"
	Ant  []*Register // Previous
	Ants []*Register // My Previous
}

func NewRegister(id string) *Register {
	typ := id[0:1]

	return &Register{Id: id, Type: typ, Ant: make([]*Register, 0), Ants: make([]*Register, 0)}
}

func (r *Register) MarkAsAnt(a *Register) {
	r.Ant = append(r.Ant, a)
}

func (r *Register) SetAnt(a *Register) {
	r.Ants = append(r.Ants, a)
}

func (r *Register) Show() {
	str := "Registo: " + r.Id + "\n"

	if len(r.Ants) == 0 {
		str += "\tNão tem Registro Anterior.\n"
	} else {
		str += "\tTem como Registro Anterior:\n"
		for _, ant := range r.Ants {
			str += "\t" + ant.Id + "\n"
		}
	}

	if len(r.Ant) == 0 {
		str += "\tNão é Registro Anterior.\n"
	} else {
		str += "\tÉ Registro Anterior:\n"
		for _, ant := range r.Ant {
			str += "\t" + ant.Id + "\n"
		}
	}

	println(str)
}

type Database struct {
	Registers map[string]*Register
}

func NewDatabase() *Database {
	return &Database{make(map[string]*Register)}
}

func (d *Database) Add(r *Register) {
	d.Registers[r.Id] = r
}

func (d *Database) Get(id string) (*Register, error) {
	reg := d.Registers[id]
	if reg == nil || reg.Id == "" {
		return nil, errors.New("Register not found")
	}

	return reg, nil
}

func main() {
	content, err := os.ReadFile("test.db")
	if err != nil {
		panic(err)
	}

	db := NewDatabase()

	for _, line := range strings.Split(string(content), "\n") {
		if line == "" {
			continue
		}

		rId := strings.Split(line, ":")[0]
		rAnts := strings.Split(strings.Split(line, ":")[1], ",")

		reg, err := db.Get(rId)
		if err != nil {
			reg = NewRegister(rId)
		}

		for _, ant := range rAnts {
			// Get Ant Reg
			antReg, err := db.Get(ant)
			if err != nil {
				antReg = NewRegister(ant)
				db.Add(antReg)
			}

			antReg.MarkAsAnt(reg)
			reg.SetAnt(antReg)
		}

		db.Add(reg)
	}

	println("Number of registers: ", len(db.Registers))
	for _, reg := range db.Registers {
		reg.Show()
	}
}
