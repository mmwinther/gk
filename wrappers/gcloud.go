package wrappers

import (
	"encoding/json"
	"os/exec"
)

func SetGcloudProject(gproject string) error {

	_,err := exec.Command("gcloud", "config", "set", "project", gproject).Output()
	if err != nil {
		return err
	}
	return nil
}

func GetAllGcloudProjects() ([]string, error) {
	type projectsList struct {
		ProjectId string `json:"ProjectId"`
	}

	rawProjects,err := exec.Command("gcloud", "projects", "list", "--format=json").Output()
	if err != nil {
		return nil,err
	}
	var projects []projectsList
	err = json.Unmarshal(rawProjects, &projects)
	if err != nil {
		return nil,err
	}
	var allProjects []string = nil
	for i:= range projects {
		allProjects = append (allProjects,string(projects[i].ProjectId))
	}

	return allProjects,nil
}

func GetAllGcloudK8s() (map[string]string,error) {
	type k8sclustersList struct {
		Name string `json:"name"`
		Zone string `json:"zone"`
	}
	k8sClusters,err := exec.Command("gcloud", "container", "clusters","list", "--format=json").Output()
	if err != nil {
		return nil,err
	}
	var clusters []k8sclustersList
	err = json.Unmarshal(k8sClusters, &clusters)
	if err != nil {
		return nil,err
	}

	allClusters := make(map[string]string)
	for i:= range clusters {
		allClusters[string(clusters[i].Name)] = string(clusters[i].Zone)
	}
	return allClusters,nil
}
