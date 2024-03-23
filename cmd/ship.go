package cmd

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/fatih/color"

	"github.com/koksmat-com/koksmat/kitchen"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var shipCmd = &cobra.Command{
	Use:   "ship ",
	Short: "Shipping handling",
	Args:  cobra.MinimumNArgs(0),
	Long:  ``,
}

func Unzip(zipfile string, dst string, rootfoldername string) error {

	archive, err := zip.OpenReader(zipfile)
	if err != nil {
		return err
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(dst, f.Name)
		//		fmt.Println("unzipping file ", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", filePath)

		}
		if f.FileInfo().IsDir() {
			//			fmt.Println("creating directory...")
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return err
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return err
		}

		dstFile.Close()
		fileInArchive.Close()
	}

	destPath := path.Join(dst, rootfoldername)
	err = os.RemoveAll(destPath)
	if err != nil {
		return err
	}
	err = os.Rename(path.Join(dst, archive.File[0].Name), destPath)
	if err != nil {
		return err
	}

	return nil
}

func shipSubCmd(use string, short string, long string, minargs int, run func(cmd *cobra.Command, args []string)) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Args:  cobra.MinimumNArgs(minargs),
		Run:   run,
	}
	// cmd.Flags().StringVarP(&kitchenName, "kitchen", "k", "", "Kitchen (required)")
	// cmd.MarkFlagRequired("kitchen")
	// cmd.Flags().StringVarP(&stationName, "station", "s", "", "Station (required)")
	// cmd.MarkFlagRequired("station")
	return cmd

}

func getPackageNameFromURL(url string) string {
	return strings.Split(path.Base(url), "@")[0]

}

func shipGetCmd(cmd *cobra.Command, args []string) {
	color.White("Downloading mate")
	kitchenRoot := viper.GetString("KITCHENROOT")
	packagePath := path.Join(kitchenRoot, ".koksmat", "packages")

	kitchen.CreateIfNotExists(packagePath, 0755)

	url := "https://koksmat.blob.core.windows.net/packages/koksmat-mate.zip?se=2026-09-03T14%3A30%3A13Z&sp=r&sv=2022-11-02&sr=b&sig=TbggJHwLnJw8cyz3Os4MKSnNirmJw667enYY6p5AJOI%3D"

	dest := path.Join(packagePath, "koksmat-mate.zip")
	kitchen.Download(url, dest)
	err := Unzip(dest, packagePath, "koksmat-mate")
	if err != nil {
		color.Red(err.Error())
		return
	}
	color.White("Installing packages")
	installPnpm := exec.Command("npm", "install", "-g", "pnpm")
	installPnpm.CombinedOutput()

	execCmd := exec.Command("pnpm", "install")
	execCmd.Dir = path.Join(packagePath, "koksmat-mate", ".koksmat", "web")
	result, err := execCmd.CombinedOutput()

	if err != nil {
		color.Red(string(result))
		return

	}
	//	color.White(string(result))

	color.White("Building")
	execCmd2 := exec.Command("pnpm", "build")
	execCmd2.Dir = path.Join(packagePath, "koksmat-mate", ".koksmat", "web")
	result2, err := execCmd2.CombinedOutput()
	if err != nil {
		color.Red(string(result2))
		return

	}
	//color.White(string(result2))
	color.White("Mate ready to launch - run")
	color.Green("koksmat sail")

}
func init() {

	rootCmd.AddCommand(shipCmd)
	shipCmd.AddCommand(shipSubCmd("get [package]", "Get package", "", 1, shipGetCmd))

}
