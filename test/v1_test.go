package test

import (
	"fmt"
	"github.com/worldiety/gdoc/internal/api"
)

func ExampleTestFunc() {
	fmt.Print(api.TestFunc("worldiety"))
	// Output: Hello worldiety
}
