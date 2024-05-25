package kitchen

import (
	"path"
	"testing"

	"github.com/spf13/viper"
)

func TestCompareKitchens(t *testing.T) {
	t.Log("Compare")
	root := viper.GetString("KITCHENROOT")
	master := path.Join(root, "magic-people")
	replica := path.Join(root, "magic-mix")
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
func TestCompareKitchens2(t *testing.T) {
	t.Log("Compare")
	root := viper.GetString("KITCHENROOT")
	master := path.Join(root, "magic-people")
	replica := path.Join(root, "magic-mix")
	subfolders := []string{".koksmat/web/koksmat"} // , ".koksmat/web/koksmat/msal"}

	_, err := Compare(master, replica, subfolders, true, *&CompareOptions{PrintResults: true})
	if err != nil {
		t.Errorf("Error: %v", err)
	}
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
