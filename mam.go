package mamgo

import (
	"fmt"

	"github.com/gookit/ini/v2"
)

func main() {
	fmt.Println("Loading configuration file")
	err := ini.LoadExists("mamgo.ini")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load config file: %v\n", err)
	}
	fmt.Fprintf("Main ToDo: %s\n", ini.String("main.todo"))
}
