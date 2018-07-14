# front 

Extracts frontmatter from text files with ease.

## Features

* Custom delimiters. You are free to register any delimiter of your choice. Provided its a three character string. e.g `+++`,  `$$$`,  `---`,  `%%%`
* Custom Handlers. Anything that implements `HandlerFunc` can be used to decode the values from the frontmatter text, you can see the `JSONHandler` for how to implement one.
* Support YAML frontmatter
* Support JSON frontmatter.

## Installation

```sh
go get github.com/cirocosta/front
```

## How to use

```go
package main

import (
	"fmt"
	"strings"

	"github.com/cirocosta/front"
)

var txt = `+++
{
    "title":"front"
}
+++

# Body
Over my dead body
`

type MyStruct struct {
        Title `json:"title"`
}

func main() {
        var f = &MyStruct{}

	m := front.NewMatter()
	m.Handle("+++", front.JSONHandler)
	body, err := m.Parse(strings.NewReader(txt), f)
	if err != nil {
		panic(err)
	}

	fmt.Printf("The front matter is:\n%#v\n", f)
	fmt.Printf("The body is:\n%q\n", body)
}
```

Please see the tests formore details

## Licence

This project is under the MIT Licence. See the [LICENCE](LICENCE) file for the full license text.

