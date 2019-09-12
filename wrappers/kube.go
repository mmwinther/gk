package wrappers

import (
	"encoding/json"
	"os/exec"
	"strings"
)

type kubeContextDetails struct {
	Name string `json:"name"`
	Context struct {
		Cluster string
		User string
	}
}

type kubeConfig struct {
	Contexts []kubeContextDetails
}

// KubeContexts returns all kubeconfig contexts
// present in ~/.kube/config
func KubeContexts() ([]string, error) {

	var kConfigs kubeConfig
	var kContexts []string

	kconfig, err := exec.Command("kubectl", "config", "view", "-o=json","--raw=true").Output()
	if err != nil {
		return nil,err
	}

	err = json.Unmarshal(kconfig,&kConfigs)
	if err != nil {
		return nil,err
	}

	for k:= range kConfigs.Contexts {
		kContexts = append(kContexts,string(kConfigs.Contexts[k].Name))
	}

	return kContexts, nil
}

// SetKubeContext sets context as current in kubeconfig
// set default namespace for the context if namespace is defined
func SetKubeContext(context string, namespace string) error  {

	if namespace != "" {
		_, err := exec.Command("kubectl", "config", "set-context", "--current",  "--namespace="+namespace).Output()
		if err != nil {
			return err
		}
	} else {
		_,err := exec.Command("kubectl", "config", "use-context", context).Output()
		if err != nil {
			return err
		}
	}
	return nil
}

func CurrentKubeContext() ([]byte, error) {
	currentcontext, err := exec.Command("kubectl", "config", "current-context").Output()
	if err != nil {
		return nil, err
	}
	return currentcontext,nil
}

func GetNamespaces() ([]string, error) {

	namespaces, err := exec.Command("kubectl", "get", "namespaces","-ojsonpath={.items[*].metadata.name}").Output()
	if err != nil {
		return nil,err
	}
	nsSlice := strings.Split(string(namespaces)," ")
	return nsSlice,nil
}

func CheckNamespaceExists(namespace string) (bool,error) {
	namespaces, err := GetNamespaces()
	if err != nil {
		return false, err
	}

	for _,n := range namespaces {
		if namespace == n{
			return true, nil
		}
	}
	return false, nil

}

//TODO: check types and functions naming( small,capital,camelcase : public,private)