package plugin

import (
	"bufio"
	"log"
	"os"

	a1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

func (b *RowBuilder) loadYaml(filename string) ([]v1.Pod, error) {
	var pods []v1.Pod
	var content string
	var scanner *bufio.Scanner

	if b.StdinChanged {
		// read yaml from stdin
		file := bufio.NewReader(os.Stdin)
		scanner = bufio.NewScanner(file)
	} else {
		// load yaml file
		file, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
	}

	for scanner.Scan() {
		line := scanner.Text()
		if line == "---" {
			pod, err := b.convertFromYaml([]byte(content))
			if err != nil {
				return []v1.Pod{}, err
			}
			pods = append(pods, pod)
			content = ""
		} else {
			content += line + "\n"
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	pod, err := b.convertFromYaml([]byte(content))
	if err != nil {
		return []v1.Pod{}, err
	}
	pods = append(pods, pod)

	return pods, nil
}

func (b *RowBuilder) convertFromYaml(input []byte) (v1.Pod, error) {
	// var allPods []v1.Pod
	var err error
	var pod v1.Pod
	var newPod v1.Pod

	// Happy accident, it looks like pod unmarshalling sets the kind field so we dont have to guess,
	//  not sure if this is intended so it might break in the future
	err = yaml.Unmarshal(input, &pod)
	if err == nil {
		if pod.Kind == "Pod" {
			newPod = pod
		}
	} else {
		return v1.Pod{}, err
	}

	switch pod.Kind {
	case "Deployment":
		var deploySpec a1.Deployment
		err = yaml.Unmarshal(input, &deploySpec)
		if err == nil {
			podTemplate := deploySpec.Spec.Template
			newPod = v1.Pod{
				Spec: podTemplate.Spec,
			}
			newPod.SetName(deploySpec.Name)
		} else {
			return v1.Pod{}, err
		}

	case "ReplicaSet":
		var replicaSpec a1.ReplicaSet
		err = yaml.Unmarshal(input, &replicaSpec)
		if err == nil {
			podTemplate := replicaSpec.Spec.Template
			newPod = v1.Pod{
				Spec: podTemplate.Spec,
			}
			newPod.SetName(replicaSpec.Name)
		} else {
			return v1.Pod{}, err
		}

	case "StatefulSet":
		var statefulSpec a1.StatefulSet
		err = yaml.Unmarshal(input, &statefulSpec)
		if err == nil {
			podTemplate := statefulSpec.Spec.Template
			newPod = v1.Pod{
				Spec: podTemplate.Spec,
			}
			newPod.SetName(statefulSpec.Name)
		} else {
			return v1.Pod{}, err
		}

	case "DaemonSet":
		var daemonSpec a1.DaemonSet
		err = yaml.Unmarshal(input, &daemonSpec)
		if err == nil {
			podTemplate := daemonSpec.Spec.Template
			newPod = v1.Pod{
				Spec: podTemplate.Spec,
			}
			newPod.SetName(daemonSpec.Name)
		} else {
			return v1.Pod{}, err
		}

	case "Job":
		var jobSpec batchv1.Job
		err = yaml.Unmarshal(input, &jobSpec)
		if err == nil {
			podTemplate := jobSpec.Spec.Template
			newPod = v1.Pod{
				Spec: podTemplate.Spec,
			}
			newPod.SetName(jobSpec.Name)
		} else {
			return v1.Pod{}, err
		}

	case "CronJob":
		var cronJobSpec batchv1.CronJob
		err = yaml.Unmarshal(input, &cronJobSpec)
		if err == nil {
			podTemplate := cronJobSpec.Spec.JobTemplate
			newPod = v1.Pod{
				Spec: podTemplate.Spec.Template.Spec,
			}
			newPod.SetName(cronJobSpec.Name)
		} else {
			return v1.Pod{}, err
		}
	}

	return newPod, nil
}
