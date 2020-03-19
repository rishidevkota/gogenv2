package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

//--------------------
type User struct {
	ID      uint
	Name    string
	Address string
}

type Order struct {
	ID   uint
	Item string
}

//--------------------

type St struct {
	Name   string
	Fields []string
	Types  []string
	Values []interface{}
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func GenSts(args ...interface{}) []*St {
	sts := make([]*St, 0)
	for _, v := range args {
		st := new(St)
		st.Name = strings.Split(fmt.Sprintf("%T", v), ".")[1]

		e := reflect.ValueOf(v).Elem()

		for i := 0; i < e.NumField(); i++ {
			// varName := e.Type().Field(i).Name
			// varType := e.Type().Field(i).Type
			// varValue := e.Field(i).Interface()
			// fmt.Printf("%v %v %v\n", varName, varType, varValue)

			st.Fields = append(st.Fields, e.Type().Field(i).Name)
			st.Types = append(st.Types, fmt.Sprintf("%v", e.Type().Field(i).Type))
			st.Values = append(st.Values, e.Field(i).Interface())

		}

		sts = append(sts, st)
	}

	return sts
}

func main() {
	//----------------------
	sts := GenSts(&User{ID: 1, Name: "rishi"}, &Order{})
	//----------------------

	for _, v := range sts {
		//create
		fmt.Printf("func (%v *%v) Create() error{\n", ToSnakeCase(v.Name), v.Name)
		fmt.Printf("\terr := db.Exec(\"insert into %v (", ToSnakeCase(v.Name))
		size := len(v.Fields)
		fmt.Printf("%v", ToSnakeCase(v.Fields[1]))
		if size > 2 {
			for i := 2; i < size; i++ {
				fmt.Printf("%v", " ,"+ToSnakeCase(v.Fields[i]))
			}
		}

		fmt.Print(") (?")
		if size > 2 {
			for i := 2; i < size; i++ {
				fmt.Print(" ,?")
			}
		}
		fmt.Print(")\",\n")
		for i := 1; i < size; i++ {
			fmt.Printf("\t\t&%v.%v,\n", ToSnakeCase(v.Name), v.Fields[i])
		}
		fmt.Print("\t)\n")
		fmt.Println("")
		fmt.Print("\tif err != nil {\n\t\treturn err\n\t}\n")
		fmt.Print("\treturn nil\n}\n\n")

		//update update table_name set col_name=value where id=
		fmt.Printf("func (%v *%v) Update() error{\n", ToSnakeCase(v.Name), v.Name)
		fmt.Printf("\terr := db.Exec(`update %v set ", ToSnakeCase(v.Name))
		fmt.Printf("%v", ToSnakeCase(v.Fields[1])+"=?")
		if size > 2 {
			for i := 2; i < size; i++ {
				fmt.Printf("%v", " "+ToSnakeCase(v.Fields[i])+"=?")
			}
		}
		fmt.Print(" where id=?`,\n")
		for i := 1; i < size; i++ {
			fmt.Printf("\t\t&%v.%v,\n", ToSnakeCase(v.Name), v.Fields[i])
		}
		fmt.Printf("\t\t&%v.ID,\n\t)\n\n", ToSnakeCase(v.Name))
		fmt.Print("\tif err != nil {\n\t\treturn err\n\t}\n")
		fmt.Print("\treturn nil\n}\n\n")

		//DELETE FROM table_name WHERE condition;
		fmt.Printf("func (%v *%v) Delete() error{\n", ToSnakeCase(v.Name), v.Name)
		fmt.Printf("\terr := db.Exec(`delete from %v where id=?`,&%v.ID)\n\n", ToSnakeCase(v.Name), ToSnakeCase(v.Name))
		fmt.Print("\tif err != nil {\n\t\treturn err\n\t}\n")
		fmt.Print("\treturn nil\n}\n\n")
	}

}
