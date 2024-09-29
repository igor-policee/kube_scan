package main

import (
    "context"
    "flag"
    "fmt"
    "log"
    "os"
    "path/filepath"

    v1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

func main() {
    // Command-line parameters
    namespace := flag.String("namespace", "default", "Namespace to scan")
    kubeconfig := flag.String("kubeconfig", filepath.Join(homeDir(), ".kube", "config"), "Path to the kubeconfig file")
    flag.Parse()

    // Create client configuration
    config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
    if err != nil {
        log.Fatalf("Error building kubeconfig: %v", err)
    }

    // Create Kubernetes client
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        log.Fatalf("Error creating Kubernetes client: %v", err)
    }

    // Get list of pods in the specified namespace
    pods, err := clientset.CoreV1().Pods(*namespace).List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        log.Fatalf("Error retrieving pods: %v", err)
    }

    // Iterate through pods and display required parameters
    for _, pod := range pods.Items {
        fmt.Printf("POD: %s\n", pod.Name)

        // Display pod-level spec parameters
        if pod.Spec.SecurityContext != nil {
            fmt.Printf("  spec.securityContext.fsGroup: %v\n", pod.Spec.SecurityContext.FSGroup)
            fmt.Printf("  spec.securityContext.runAsUser: %v\n", pod.Spec.SecurityContext.RunAsUser)
        } else {
            fmt.Println("  spec.securityContext: nil")
        }

        fmt.Printf("  spec.hostPID: %v\n", pod.Spec.HostPID)
        fmt.Printf("  spec.hostIPC: %v\n", pod.Spec.HostIPC)
        fmt.Printf("  spec.hostNetwork: %v\n", pod.Spec.HostNetwork)
        fmt.Printf("  spec.nodeName: %s\n", pod.Spec.NodeName)
        fmt.Printf("  spec.serviceAccountName: %s\n", pod.Spec.ServiceAccountName)

        // Iterate through volumes
        fmt.Println("  spec.volumes:")
        for _, vol := range pod.Spec.Volumes {
            fmt.Printf("    - %s\n", vol.Name)
        }

        // Function to process containers and initContainers
        processContainers := func(containers []v1.Container, containerType string) {
            for _, container := range containers {
                fmt.Printf("  %s:\n", containerType)
                fmt.Printf("    name: %s\n", container.Name)
                fmt.Printf("    image: %s\n", container.Image)

                // Display readinessProbe
                if container.ReadinessProbe != nil {
                    fmt.Println("    readinessProbe:")
                    if container.ReadinessProbe.HTTPGet != nil {
                        fmt.Println("      httpGet:")
                        fmt.Printf("        path: %s\n", container.ReadinessProbe.HTTPGet.Path)
                        fmt.Printf("        port: %v\n", container.ReadinessProbe.HTTPGet.Port)
                        if container.ReadinessProbe.HTTPGet.Scheme != "" {
                            fmt.Printf("        scheme: %s\n", container.ReadinessProbe.HTTPGet.Scheme)
                        }
                    }
                    if container.ReadinessProbe.Exec != nil {
                        fmt.Println("      exec:")
                        fmt.Printf("        command: %v\n", container.ReadinessProbe.Exec.Command)
                    }
                    if container.ReadinessProbe.TCPSocket != nil {
                        fmt.Println("      tcpSocket:")
                        fmt.Printf("        port: %v\n", container.ReadinessProbe.TCPSocket.Port)
                    }
                } else {
                    fmt.Println("    readinessProbe: nil")
                }

                // Display ports
                if len(container.Ports) > 0 {
                    fmt.Println("    ports:")
                    for _, port := range container.Ports {
                        fmt.Printf("      - containerPort: %d\n", port.ContainerPort)
                        if port.HostPort != 0 {
                            fmt.Printf("        hostPort: %d\n", port.HostPort)
                        }
                        if port.Name != "" {
                            fmt.Printf("        name: %s\n", port.Name)
                        }
                        if port.Protocol != "" {
                            fmt.Printf("        protocol: %s\n", port.Protocol)
                        }
                    }
                } else {
                    fmt.Println("    ports: none")
                }

                // Display securityContext parameters
                if container.SecurityContext != nil {
                    if container.SecurityContext.Capabilities != nil {
                        if len(container.SecurityContext.Capabilities.Add) > 0 {
                            fmt.Printf("    capabilities.add: %v\n", container.SecurityContext.Capabilities.Add)
                        }
                        if len(container.SecurityContext.Capabilities.Drop) > 0 {
                            fmt.Printf("    capabilities.drop: %v\n", container.SecurityContext.Capabilities.Drop)
                        }
                    }
                    if container.SecurityContext.RunAsUser != nil {
                        fmt.Printf("    runAsUser: %v\n", container.SecurityContext.RunAsUser)
                    }
                    if container.SecurityContext.AllowPrivilegeEscalation != nil {
                        fmt.Printf("    allowPrivilegeEscalation: %v\n", container.SecurityContext.AllowPrivilegeEscalation)
                    }
                    if container.SecurityContext.Privileged != nil {
                        fmt.Printf("    privileged: %v\n", container.SecurityContext.Privileged)
                    }
                } else {
                    fmt.Println("    securityContext: nil")
                }
            }
        }

        // Process containers
        processContainers(pod.Spec.Containers, "containers")

        // Process initContainers
        processContainers(pod.Spec.InitContainers, "initContainers")

        fmt.Println("-----------------------------------------------------")
    }
}

// Function to get the home directory
func homeDir() string {
    if h := os.Getenv("HOME"); h != "" {
        return h
    }
    return os.Getenv("USERPROFILE") // For Windows
}
