package cmd

import (
	"fmt"

	"github.com/mattherman/mhgit/objects"
	"github.com/spf13/cobra"
)

// catFileCmd represents the catFile command
var catFileCmd = &cobra.Command{
	Use:   "cat-file [object]",
	Short: "Provide content or type and size information for repository objects.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		CatFile(args[0], prettyPrint, outputObjectType, outputObjectSize)
	},
}

var prettyPrint bool
var outputObjectType bool
var outputObjectSize bool

func init() {
	rootCmd.AddCommand(catFileCmd)
	catFileCmd.Flags().BoolVarP(&prettyPrint, "pretty", "p", false, "Pretty-print the object based on type.")
	catFileCmd.Flags().BoolVarP(&outputObjectType, "type", "t", false, "Output the type of the object.")
	catFileCmd.Flags().BoolVarP(&outputObjectSize, "size", "s", false, "Output the size of the object.")
}

// CatFile will inspect a stored Git object or return an error if it
// cannot be found.
func CatFile(objectName string, outputObject bool, outputType bool, outputSize bool) {
	obj, err := objects.ReadObject(objectName)

	if err != nil {
		fmt.Println(err)
	} else {
		if outputType {
			fmt.Println(obj.ObjectType)
		} else if outputSize {
			fmt.Printf("%d\n", len(obj.Data))
		} else {
			fmt.Printf("%s", obj.Data)
		}
	}
}
