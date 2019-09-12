package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/zaikinlv/gk/wrappers"
	"log"
)


var (
	setGoogleProjectOnly bool
	setKubernetesClusterOnly bool
	importClusterFromProject bool
	deleteAllContexts bool
	printInfo bool
	setNamespace string

	rootCmd = &cobra.Command{
		Use:   "gk",
		Short: "Kubectl and Gcloud wrapper",
		Long: `Kubectl and Gcloud wrapper`,
		Args: cobra.NoArgs,
		Run: runCommand,
	}
)

func init() {
	rootCmd.Flags().BoolVarP(&setGoogleProjectOnly,"project","p",false,"only set google project project")
	rootCmd.Flags().BoolVarP(&setKubernetesClusterOnly,"kubecluster","k",false,"only set kubernetes cluster from kubeconfig")
	rootCmd.Flags().BoolVarP(&importClusterFromProject,"config","c",false,"import cluster credentials from google cloud to local kubeconfig")
	rootCmd.Flags().BoolVarP(&deleteAllContexts,"clean","",false,"clean your local kubeconfig - delete all contexts")
	rootCmd.Flags().BoolVarP(&printInfo,"info","i",false,"print info about current state of things")
	rootCmd.Flags().StringVarP(&setNamespace,"ns","n","","set default namespace for the current context")
	rootCmd.Flag("ns").NoOptDefVal = "--view-all-namespaces"
}

// Execute starts root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func runCommand(cmd *cobra.Command, args []string){

	if setGoogleProjectOnly {
		var projects = []string{}
		projects,err := wrappers.GetAllGcloudProjects()
		if err != nil {
			log.Fatal(err)
		}

		var qsProject = []*survey.Question{
			{
				Name: "projects",
				Prompt: &survey.Select{
					Message: "Choose a project:",
					Options: projects,
					PageSize: len (projects),
				},
			},
		}
		answersProject := struct {
			Projects string `survey:"projects"`
		}{}

		err = survey.Ask(qsProject, &answersProject)
		if err != nil {
			log.Fatal(err)
		}
		err = wrappers.SetGcloudProject(answersProject.Projects)
	}

	if setKubernetesClusterOnly {
		kubeContexts, err := wrappers.KubeContexts()
		if err != nil {
			log.Fatal(err)
		}
		var qsContext = []*survey.Question{
			{
				Name: "contexts",
				Prompt: &survey.Select{
					Message:  "Choose a context:",
					Options:  kubeContexts,
					PageSize: len(kubeContexts),
				},
			},
		}
		answersContext := struct {
			Contexts string `survey:"contexts"`
		}{}

		err = survey.Ask(qsContext, &answersContext)
		if err != nil {
			log.Fatal(err)
		}
		err = wrappers.SetKubeContext(answersContext.Contexts,"")
		if err != nil {
			log.Fatal(err)
		}
	}

	if importClusterFromProject {
		fmt.Println("Importing kubernetes cluster credentials")
	}

	if deleteAllContexts {
		fmt.Println("Deleting all kube contexts")
	}

	if printInfo {
		fmt.Println("Print current environment")
	}

	if setNamespace == "--view-all-namespaces" {

		namespaces,err := wrappers.GetNamespaces()
		if err != nil {
			log.Fatal(err)
		}

		var qsNamespace = []*survey.Question{
		{
			Name: "namespaces",
			Prompt: &survey.Select{
				Message: "Choose a namespace:",
				Options: namespaces,
				PageSize: len (namespaces),
			},
		},
	}
	answersNamespace := struct {
		Namespaces string `survey:"namespaces"`
	}{}

	err = survey.Ask(qsNamespace, &answersNamespace)
	if err != nil {
		log.Fatal(err)
	}

	err = wrappers.SetKubeContext("", answersNamespace.Namespaces)
	if err != nil {
		log.Fatal(err)
	}

	} else if setNamespace != "" {

		namespaceExists, err := wrappers.CheckNamespaceExists(setNamespace)
		if err != nil {
			log.Fatal(err)
		}
		if namespaceExists {
			err := wrappers.SetKubeContext("", setNamespace)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatalf("Typed namespace: '%s' doesn't exist in the current cluster", setNamespace)
		}
	}
}

//TODO: variables visibility scope