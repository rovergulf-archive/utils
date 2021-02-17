# rovergulf/utils/colors

Go package to generate random colors

```go
package main

import (
	"github.com/rovergulf/utils/colors"
)

func main() {
	var colorInHex string = colors.GetRandomColorInHex()
	var colorInRGB colors.RGBColor = colors.GetRandomColorInRgb()
	var colorInHSV colors.HSVColor = colors.GetRandomColorInHSV()
}
```