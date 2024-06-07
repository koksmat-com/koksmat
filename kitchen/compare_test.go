package kitchen

import (
	"encoding/json"
	"fmt"
	"log"
	"path"
	"testing"

	"github.com/spf13/viper"
)

func TestCompareKitchens(t *testing.T) {
	t.Log("Compare")
	root := viper.GetString("KITCHENROOT")
	master := path.Join(root, "magic-people")
	replica := path.Join(root, "magic-files")
	subfolders := []string{".koksmat/web/koksmat", ".koksmat/web/app/magic/components"} // , ".koksmat/web/koksmat/msal"}

	_, err := Compare(master, replica, subfolders, true, *&CompareOptions{
		CopyFunction: Copy,
		//MergeFunction: Merge,
		PrintMergeLink: true,
		PrintResults:   true})
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func printJSON(v any) {
	j, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Println(err)
	}
	s := string(j)
	fmt.Print(s)
	// err = clipboard.Init()
	// if err == nil {
	// 	clipboard.Write(clipboard.FmtText, []byte(s))
	// }

}
func TestCompareKitchens2(t *testing.T) {
	t.Log("Compare")
	root := viper.GetString("KITCHENROOT")
	master := path.Join(root, "magic-master")
	replica := path.Join(root, "magic-files")
	subfolders := []string{".koksmat/web/app/koksmat", ".koksmat/web/app/magic"} // , ".koksmat/web/koksmat/msal"}

	result, err := Compare(master, replica, subfolders, true, *&CompareOptions{PrintResults: true, PrintMergeLink: true})
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	printJSON(result)
}

func TestCompare(t *testing.T) {
	t.Log("Compare")
	subfolders := []string{"case1", "case2", "case3", "case4"} // , ".koksmat/web/koksmat/msal"}

	_, err := Compare("testdata/master", "testdata/replica", subfolders, true, *&CompareOptions{PrintResults: true})
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestCompare2(t *testing.T) {
	t.Log("Compare")
	subfolders := []string{"case5"} // , ".koksmat/web/koksmat/msal"}

	_, err := Compare("testdata/master", "testdata/replica", subfolders, true, *&CompareOptions{CopyFunction: Copy, PrintResults: false})
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}
