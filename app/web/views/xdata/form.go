package xdata

import (
	"encoding/json"
	"fmt"
)

const houseKeeping = "reset() {this.form = JSON.parse(JSON.stringify(this.initial))}, isModified() {return JSON.stringify(this.form) !== JSON.stringify(this.initial)}"

func ToFormData(in any) string {
	return fmt.Sprintf(`{form: %[1]s, initial: %[1]s, %s }`, ToJsonData(in), houseKeeping)
}

func ToJsonData(in any) string {
	out, _ := json.Marshal(in)
	return string(out)
}
