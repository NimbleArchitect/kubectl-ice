package plugin

import (
	"io/ioutil"

	a1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

var yamlfile = `
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: myappdeploy
  name: myapp
spec:
  replicas: 2
  selector:
    matchLabels:
      app: myappdeploy
  template:
    metadata:
      labels:
            app: myappdeploy
    spec:
      containers:
      - name: frontend
        image: python:latest
        command: ['python', '/myapp/randomcpuapp.py']
        ports:
          - containerPort: 8080
        resources:
          requests:
            cpu: "125m"
            memory: "1M"
          limits:
            cpu: "1"
            memory: 256M
        volumeMounts:
          - name: app
            mountPath: /myapp/
        livenessProbe:
          exec:
            command:
            - /bin/true
          initialDelaySeconds: 10
          periodSeconds: 5
        
      - name: nginx
        image: nginx:1.7.9
        ports:
        - containerPort: 80
        resources:
          requests:
            cpu: "1m"
            memory: "1M"
          limits:
            cpu: "1"
            memory: 256M
        livenessProbe:
          httpGet:
            path: /
            port: 80
          initialDelaySeconds: 60
          failureThreshold: 8
          periodSeconds: 60

      volumes:
      - name: app
        configMap:
          name: app.py
          defaultMode: 0777
          items:
          # - key: mainapp
          - key: randomcpu
            path: randomcpuapp.py
`

func loadYaml(filename string) ([]v1.Pod, error) {

	// load yaml file
	//TODO: doesn't read multi part yaml files :(
	content, err := ioutil.ReadFile(filename)

	if err != nil {
		return []v1.Pod{}, err
	}

	return convertFromYaml([]byte(content))
}

func convertFromYaml(input []byte) ([]v1.Pod, error) {
	var allPods []v1.Pod
	var err error

	var pod v1.Pod
	err = yaml.Unmarshal(input, &pod)
	if err == nil {
		allPods = append(allPods, pod)
	} else {
		return []v1.Pod{}, err
	}

	var deploySpec a1.Deployment
	err = yaml.Unmarshal(input, &deploySpec)
	if err == nil {
		podTemplate := deploySpec.Spec.Template
		pod := v1.Pod{
			Spec: podTemplate.Spec,
		}
		pod.SetName(deploySpec.Name)
		allPods = append(allPods, pod)
	} else {
		return []v1.Pod{}, err
	}

	var replicaSpec a1.ReplicaSet
	err = yaml.Unmarshal(input, &replicaSpec)
	if err == nil {
		podTemplate := replicaSpec.Spec.Template
		pod := v1.Pod{
			Spec: podTemplate.Spec,
		}
		pod.SetName(replicaSpec.Name)
		allPods = append(allPods, pod)
	} else {
		return []v1.Pod{}, err
	}

	var statefulSpec a1.StatefulSet
	err = yaml.Unmarshal(input, &statefulSpec)
	if err == nil {
		podTemplate := statefulSpec.Spec.Template
		pod := v1.Pod{
			Spec: podTemplate.Spec,
		}
		pod.SetName(statefulSpec.Name)
		allPods = append(allPods, pod)
	} else {
		return []v1.Pod{}, err
	}

	var daemonSpec a1.DaemonSet
	err = yaml.Unmarshal(input, &daemonSpec)
	if err == nil {
		podTemplate := daemonSpec.Spec.Template
		pod := v1.Pod{
			Spec: podTemplate.Spec,
		}
		pod.SetName(daemonSpec.Name)
		allPods = append(allPods, pod)
	} else {
		return []v1.Pod{}, err
	}

	return allPods, nil
}
