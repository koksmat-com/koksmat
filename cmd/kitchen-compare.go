package cmd

import (
	"log"
	"path"

	"github.com/koksmat-com/koksmat/kitchen"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var subfolders = []string{
	".koksmat/web/koksmat",
	".koksmat/web/lib",
	".koksmat/web/app/magic/components"} // , ".koksmat/web/koksmat/msal"}

func KitchenCompareCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "compare ",
		Short: "Compare kitchens",
		Long:  ``,
	}

	web := &cobra.Command{
		Use:   "web master replica",
		Short: "Web",
		Long:  ``,
		Example: `
Compare the web folder of the master kitchen with the web folder of the replica kitchen

koksmat kitchen compare web magic-people magic-files 
		`,
		Args: cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			masterKitchen := args[0]
			replicaKitchen := args[1]

			root := viper.GetString("KITCHENROOT")
			master := path.Join(root, masterKitchen)
			replica := path.Join(root, replicaKitchen)

			result, err := kitchen.Compare(master, replica, subfolders, true, *&kitchen.CompareOptions{
				//CopyFunction: kitchen.Copy,
				//MergeFunction: Merge,
				//PrintMergeLink: true,
				//PrintResults:   true
			})
			if err != nil {
				log.Fatal(err)
			}
			printJSON(result)
			// kitchen := args[0]

		},
	}

	cmd.AddCommand(web)
	return cmd
}

func KitchenUpdateCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "update ",
		Short: "Update kitchen",
		Long:  ``,
	}

	web := &cobra.Command{
		Use:   "web replica",
		Short: "Web",
		Long:  ``,
		Example: `
Update the web folder of the replica kitchen with the web folder of the master kitchen

koksmat kitchen update web magic-people magic-files 
		`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			kitchenName := args[0]
			kitchenRoot := viper.GetString("KITCHENROOT")

			//source, err := repos.DownloadRepo("magicbutton", "magic-master")
			master := path.Join(kitchenRoot, "magic-master") //path.Join(*source)
			replica := path.Join(kitchenRoot, kitchenName)

			var subfolders = []string{
				".koksmat/web/app/koksmat", ".koksmat/web/app/magic"}

			j, err := kitchen.Compare(master, replica, subfolders, true, *&kitchen.CompareOptions{
				//CopyFunction: kitchen.Copy,
				//MergeFunction: Merge,
				PrintMergeLink: false,
				PrintResults:   true})
			if err != nil {
				log.Fatal(err)
			}
			printJSON(j)

		},
	}

	cmd.AddCommand(web)
	return cmd
}
