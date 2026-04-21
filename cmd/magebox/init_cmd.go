package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"qoliber/magebox/internal/cli"
	"qoliber/magebox/internal/config"
	"qoliber/magebox/internal/php"
	"qoliber/magebox/internal/project"
)

var (
	initProjectType string
)

var initCmd = &cobra.Command{
	Use:   "init [name]",
	Short: "Initialize a new MageBox project",
	Long:  "Creates a .magebox configuration file in the current directory",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runInit,
}

func init() {
	initCmd.Flags().StringVar(&initProjectType, "type", config.ProjectTypeMagento, "Project type: \"magento\" or \"laravel\"")
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	cwd, err := getCwd()
	if err != nil {
		return err
	}

	// Load global config once for defaults
	homeDir, _ := os.UserHomeDir()
	globalCfg, _ := config.LoadGlobalConfig(homeDir)
	tld := globalCfg.GetTLD()

	reader := bufio.NewReader(os.Stdin)

	// Determine project name
	var projectName string
	if len(args) > 0 {
		projectName = args[0]
	} else {
		projectName = filepath.Base(cwd)
		fmt.Printf("Project name [%s]: ", cli.Highlight(projectName))
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" {
			projectName = input
		}
	}

	// Replace slashes with dots in project name
	projectName = strings.ReplaceAll(projectName, "/", ".")

	// Prompt for PHP version
	defaultPHP := globalCfg.DefaultPHP
	fmt.Printf("PHP version [%s] (%s): ", cli.Highlight(defaultPHP), strings.Join(php.SupportedVersions, ", "))
	phpInput, _ := reader.ReadString('\n')
	phpInput = strings.TrimSpace(phpInput)
	phpVersion := defaultPHP
	if phpInput != "" {
		phpVersion = phpInput
	}

	// Check if .magebox.yaml already exists
	configPath := filepath.Join(cwd, config.ConfigFileName)
	if _, err := os.Stat(configPath); err == nil {
		cli.PrintError("%s file already exists", config.ConfigFileName)
		return nil
	}

	p, err := getPlatform()
	if err != nil {
		return err
	}

	mgr := project.NewManager(p)
	if err := mgr.Init(cwd, projectName, initProjectType, phpVersion); err != nil {
		return err
	}

	cli.PrintSuccess("Created %s for project '%s'", config.ConfigFileName, projectName)
	fmt.Println()
	fmt.Printf("Domain: %s\n", cli.URL(projectName+"."+tld))
	fmt.Println()
	cli.PrintInfo("Next steps:")
	fmt.Println(cli.Bullet("Edit " + config.ConfigFileName + " to customize your configuration"))
	fmt.Println(cli.Bullet("Run " + cli.Command("magebox start") + " to start your project"))

	return nil
}
