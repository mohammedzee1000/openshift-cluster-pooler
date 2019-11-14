package utils

import "fmt"

//PrintIfDebug prints if debug flag is set
func PrintIfDebug(debug bool, title string, out string) {
	if debug {
		fmt.Printf("%s : \n%s\n", title, out)
	}
}
