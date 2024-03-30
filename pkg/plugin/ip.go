package plugin

var ipShort = "List ip addresses of all pods in the namespace listed"

var ipDescription = ` Prints the known IP addresses of the specified pod(s). if no pod is specified the IP address of
all pods in the current namespace are shown.`

var ipExample = `  # List IP address of pods
  %[1]s ip

  # List IP address of pods output in JSON format
  %[1]s ip -o json

  # List IP address a single pod
  %[1]s ip my-pod-4jh36

  # List IP address of all pods where label app matches web
  %[1]s ip -l app=web

  # List IP address of all pods where the pod label app is either web or mail
  %[1]s ip -l "app in (web,mail)"`

// IP subcommand now points to ports command with useIP bool set to true
