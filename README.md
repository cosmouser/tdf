# tdf
package for parsing tdf / fbi formatted data

## example
```
package main

import (
	"fmt"
	"os"

	"github.com/cosmouser/tdf"
)

func main() {
	// scan each of the files in the ./downloadsE directory
	// then print the scanned tdf entries
	dir, err := os.Open("./downloadsE")
	if err != nil {
		panic(err)
	}
	names, err := dir.Readdirnames(0)
	entries := []*tdf.Entry{}
	keys := make(map[string]string)
	for _, v := range names {
		file, err := os.Open(fmt.Sprintf("./downloadsE/%s", v))
		if err != nil {
			panic(err)
		}
		defer file.Close()
		scanned := tdf.Scan(file, keys)
		if scanned != nil {
			entries = append(entries, scanned...)
		}
	}
	for _, v := range entries {
		fmt.Println(v)
	}
}

// Output:
// &{MENUENTRY2 map[UNITMENU:ARMECA MENU:2 BUTTON:11 UNITNAME:ARMVULC]}
// &{MENUENTRY3 map[UNITNAME:ARMVULC UNITMENU:ARMECK MENU:2 BUTTON:11]}
// ...
```
