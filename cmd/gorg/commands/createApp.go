package commands

import (
	"fmt"
	"os"

	"github.com/marefati110/gorg/cmd/gorg/helper"
	"github.com/spf13/cobra"
)

func setupHandlers() error {

	dir := "/handlers"
	path := helper.GetWD()

	err := os.Mkdir(path+dir, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func setupModels() error {

	dir := "/models"
	path := helper.GetWD()

	err := os.Mkdir(path+dir, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func setupConfigs() error {

	dir := "/configs"
	path := helper.GetWD()
	configDir := path + dir

	err := os.Mkdir(configDir, os.ModePerm)

	fmt.Println(err)

	f, err := os.Create(configDir + "/config.go")
	if err != nil {
		return err
	}

	defer f.Close()

	f.Write([]byte("package config"))

	return nil
}

func setupServices() error {

	dir := "/services"
	path := helper.GetWD()

	err := os.Mkdir(path+dir, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

var CreateApp = &cobra.Command{
	Use: "app",
	Run: func(cmd *cobra.Command, args []string) {

		setupModels()
		setupHandlers()
		setupServices()
		setupConfigs()
	},
}
