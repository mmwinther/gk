package main

import (
    "os/exec"
    "os"
    "log"
    "fmt"
    "flag"
    "encoding/json"
    "github.com/AlecAivazis/survey/v2"
    "errors"
    "runtime"
    "strings"
	"github.com/fatih/color"
)

const (
	appVersion = "v0.0.5"
)

func main() {
  //// TODO: checks for kubectl and gcloud are accessible on path and gcloud is configured
  // TODO: Add default value for projects list (gcloud config list)


  fullSetupPtr := flag.Bool("c",false,"set current gcp project and kubernetes cluster, context copied to kubeconfig")
  copyTokenClipboard := flag.Bool("t",false,"copy access token of current context to clipboard")
  projectonly := flag.Bool("p",false, "choose current gcp project only")
  lightinfo := flag.Bool("i",false, "print current gcp project and current selected kubecontext")
  namespace := flag.Bool("n",false,"set default kubectl namespace for current context")
  clean := flag.Bool("clean",false,"delete all kubeconfig clusters and contexts, except minikube")
  version := flag.Bool("v",false,"prints app version")
  flag.Parse()

  if *version {
  	fmt.Println(appVersion)
  	os.Exit(0)
  }
  if *clean {
  	cleanKubeconfig()
  	os.Exit(0)

  }
  if *namespace && len(flag.Args()) == 0 {
  	setNamespace("all-namespaces")
  	os.Exit(0)

  } else if *namespace && flag.Arg(0) != "" {
  	setNamespace(flag.Arg(0))
  	os.Exit(0)
  }



  if *lightinfo {
  	printLightInfo()
  	os.Exit(0)
  }

  if *projectonly {
	setCurrentProjectOnly()
	os.Exit(0)
  }


  if *copyTokenClipboard {
    mcontext,err := getCurrentContext()
    checkErr(err)
    err = tokenToClipboard(mcontext)
    os.Exit(0)
  }

  if *fullSetupPtr == false {
    err := setKubeContextOnly()
    checkErr(err)
    os.Exit(0)
  }


  //Getting projects list
  var allprojects = []string{}
  allprojects = getAllProjects()
  // Choose project
  var qs_project = []*survey.Question{
    {
        Name: "projects",
        Prompt: &survey.Select{
            Message: "Choose a project:",
            Options: allprojects,
            PageSize: len (allprojects),
            //Default: "red",
        },
    },
  }
  answers_project := struct {
    Projects string `survey:"projects"`
  }{}

  err := survey.Ask(qs_project, &answers_project)
  checkErr(err)

  // Set chosen project as active
  _,err = exec.Command("gcloud", "config", "set", "project",string(answers_project.Projects)).Output()
  checkErr(err)

  // Choose cluster
  allclusters := getAllK8s()
  if len(allclusters) == 0 {
    log.Fatal ("ERROR: Chosen project has no kubernetes clusters")
  }


  cluster_names := make([]string,0,len(allclusters))
  for k := range allclusters {
    cluster_names = append(cluster_names, k)
  }
  var qs_cluster = []*survey.Question{
    {
        Name: "clusters",
        Prompt: &survey.Select{
            Message: "Choose a cluster:",
            Options: cluster_names,
            PageSize: len(cluster_names),
            //Default: "red",
        },
    },
  }
  answers_cluster := struct {
    Clusters string `survey:"clusters"`
  }{}

  err = survey.Ask(qs_cluster, &answers_cluster)
  checkErr(err)
  // Activate cluster (get kubeconfig credentials for chosen cluster)
  _,err = exec.Command("gcloud", "container", "clusters", "get-credentials",answers_cluster.Clusters,"--zone",allclusters[answers_cluster.Clusters]).Output()

}
func cleanKubeconfig() {
	byteclusters, err := exec.Command("kubectl","config","get-clusters").Output()
	clusters := strings.Split(string(byteclusters),"\n")
	if err != nil {
		log.Fatal(err)
	}
	for _,cluster := range clusters[1:len(clusters)-1] {
		if cluster != "minikube" {
			_,err := exec.Command("kubectl","config","delete-cluster",cluster).Output()
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	contexts := getAllContexts()
	for _,context := range contexts{
		if context != "minikube" {
			_,err := exec.Command("kubectl", "config", "delete-context", context ).Output()
			if err != nil {
				log.Fatal(err)
			}
		}
	}



}
func setNamespace(ns string) {


	currentcontext, err := exec.Command("kubectl", "config", "current-context").Output()
	if err != nil {
		log.Fatal(err)
	}

	if ns != "all-namespaces" {
		_, err = exec.Command("kubectl", "config", "set-context", strings.Trim(string(currentcontext),"\n"),  "--namespace="+ns).Output()
		if err != nil {
			log.Fatal(err)
		}
		return

	}

	namespaces, err := exec.Command("kubectl", "get", "namespaces","-ojsonpath={.items[*].metadata.name}").Output()
	if err != nil {
		log.Fatal(err)
	}

	namespaces_slice := strings.Split(string(namespaces)," ")
	var qs_namespace = []*survey.Question{
		{
			Name: "namespaces",
			Prompt: &survey.Select{
				Message: "Choose a namespace:",
				Options: namespaces_slice,
				PageSize: len (namespaces_slice),
			},
		},
	}
	answers_namespace := struct {
		Namespaces string `survey:"namespaces"`
	}{}

	err = survey.Ask(qs_namespace, &answers_namespace)
	if err != nil {
		log.Fatal(err)
	}

	// Set chosen namespace as default
	_, err = exec.Command("kubectl", "config", "set-context", strings.Trim(string(currentcontext),"\n"), "--namespace="+string(answers_namespace.Namespaces)).Output()
	if err != nil {
		log.Fatal(err)
	}

}

func getAllK8s() map[string]string {
  type ClustersList struct {
    Name string `json:"name"`
    Zone string `json:"zone"`
  }
  gclusters,err := exec.Command("gcloud", "container", "clusters","list", "--format=json").Output()
  checkErr(err)
  var clusters []ClustersList
  err = json.Unmarshal(gclusters, &clusters)
  checkErr(err)

  allclusters := make(map[string]string)
  for i:= range clusters {
    allclusters[string(clusters[i].Name)] = string(clusters[i].Zone)
  }
  return allclusters
}

func printLightInfo() {


	var namespace string
	currentcontext, err := exec.Command("kubectl", "config", "current-context").Output()
	if err != nil {
		log.Fatal(err)
	}

	jsonpath := "-ojsonpath={.contexts [?(@.name=='"+strings.Trim( string(currentcontext),"\n")+"')].context.namespace}"
	currentNamespace, err := exec.Command("kubectl", "config", "view", jsonpath).Output()
	if err != nil {
		log.Fatal(err)
	}

	if string(currentNamespace) == "" {
		namespace = "default\n"
	} else {
		namespace = string(currentNamespace)+"\n"
	}


	context_slice := strings.Split(string(currentcontext),"_")

	currentproject, err := exec.Command("gcloud", "config", "get-value","core/project").Output()
	if err != nil {
		log.Fatal(err)
	}
	red := color.New(color.FgRed,color.Bold).SprintFunc()
	cyan := color.New(color.FgCyan, color.Bold).SprintfFunc()

	fmt.Printf("Gcloud set to project: %s",red(string(currentproject)))
	fmt.Printf("Kubeconfig set to project: %s, cluster: %s", red(context_slice[1]), cyan(context_slice[len(context_slice)-1]))
	fmt.Printf("Current namespace: %s", cyan(namespace))
}

func getAllContexts() []string {
  type ContextDetail struct {
    Name string `json:"name"`
    Context struct {
      Cluster string
      User string
    }
  }
  type Structure struct {
    Contexts []ContextDetail
  }

  //Getting kubeconfig contexts
  kconfig, err := exec.Command("kubectl", "config", "view", "-o=json","--raw=true").Output()
  checkErr(err)
  var contextnames Structure

  err = json.Unmarshal(kconfig,&contextnames)
  checkErr(err)

  var kcontexts = []string{}
  for k:= range contextnames.Contexts {
    kcontexts = append(kcontexts,string(contextnames.Contexts[k].Name))
  }
  return kcontexts
}

func getAllProjects() []string {

  type ProjectsList struct {
    ProjectId string `json:"ProjectId"`
  }

  gprojects,err := exec.Command("gcloud", "projects", "list", "--format=json").Output()
  checkErr(err)
  var projects []ProjectsList
  err = json.Unmarshal(gprojects, &projects)
  checkErr(err)
  var allprojects = []string{}
  for i:= range projects {
    allprojects = append (allprojects,string(projects[i].ProjectId))
  }

  return allprojects
}


func setCurrentProjectOnly() {
	var projects = []string{}
	projects = getAllProjects()
	// Choose project
	var qs_project = []*survey.Question{
		{
			Name: "projects",
			Prompt: &survey.Select{
				Message: "Choose a project:",
				Options: projects,
				PageSize: len (projects),
			},
		},
	}
	answers_project := struct {
		Projects string `survey:"projects"`
	}{}

	err := survey.Ask(qs_project, &answers_project)
	if err != nil {
		log.Fatal(err)
	}

	// Set chosen project as active
	_,err = exec.Command("gcloud", "config", "set", "project",string(answers_project.Projects)).Output()
	if err != nil {
		log.Fatal(err)
	}
}


func setKubeContextOnly() (error) {

  var allcontexts = []string{}
  allcontexts = getAllContexts()

  var qs_context = []*survey.Question{
    {
        Name: "contexts",
        Prompt: &survey.Select{
            Message: "Choose a context:",
            Options: allcontexts,
            PageSize: len(allcontexts),
            //Default: "red",
        },
    },
  }
  answers_context := struct {
    Contexts string `survey:"contexts"`
  }{}

  err := survey.Ask(qs_context, &answers_context)
  if err != nil {
    return errors.New("error in context survey")
  }
  _,err = exec.Command("kubectl", "config", "use-context", answers_context.Contexts).Output()
  if err != nil {
    return errors.New("error in kubectl config setup")
  }
  return nil
}

func getCurrentContext() (string, error) {
  mcontext,err := exec.Command("kubectl", "config", "current-context").Output()
  if err != nil {
    return "", errors.New("cannot get current context")
  }
  return strings.TrimSuffix(string(mcontext), "\n"), nil
}

func tokenToClipboard(context string) error {
  arch := runtime.GOOS
  var jsonpath string = "-o=jsonpath={.users[?(@.name=='" + context + "')].user.auth-provider.config.access-token}"
  out, err := exec.Command("kubectl", "config", "view", jsonpath).Output()
  if  err != nil {
      return errors.New("cannot get token")
  }
  toClipboard(out,arch)
  return nil
}

func toClipboard(output []byte, arch string) {
    var copyCmd *exec.Cmd

    // Mac "OS"
    if arch == "darwin" {
        copyCmd = exec.Command("pbcopy")
    }
    // Linux
    if arch == "linux" {
        copyCmd = exec.Command("xclip")
    }

    in, err := copyCmd.StdinPipe()

    if err != nil {
        log.Fatal(err)
    }

    if err := copyCmd.Start(); err != nil {
        log.Fatal(err)
    }

    if _, err := in.Write([]byte(output)); err != nil {
        log.Fatal(err)
    }

    if err := in.Close(); err != nil {
        log.Fatal(err)
    }

    copyCmd.Wait()
}

func checkErr(err error) {
    if err != nil {
        log.Fatal("ERROR:", err)
    }
}
