# tdf
package for parsing tdf / fbi formatted data

## example
```
import (
	"strings"
	"testing"
)

func TestDecode(t *testing.T) {
	reader := strings.NewReader(sample)
	nodes, err := tdf.Decode(reader)
	if err != nil {
		t.Error(err)
	}
	for i := range nodes {
		t.Log(nodes[i])
	}
}


```
