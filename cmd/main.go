package cmd

import "github.com/GoSimplicity/template/internal/pkg/di"

func main() {
	if err := di.InitViper(); err != nil {
		panic(err)
	}
}
