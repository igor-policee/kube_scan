import argparse
from kubernetes import client, config


def get_pod_security_contexts_and_parameters(namespace, kube_config):
    # Загружаем конфигурацию
    config.load_kube_config(kube_config)

    # Создаем клиент для работы с API
    v1 = client.CoreV1Api()

    # Получаем список подов в заданном namespace
    pods = v1.list_namespaced_pod(namespace)

    for pod in pods.items:
        print(f"Pod name: {pod.metadata.name}")

        # 1. spec.securityContext.fsGroup
        fs_group = pod.spec.security_context.fs_group if pod.spec.security_context else None
        print(f"  fsGroup: {fs_group}")

        # 2. spec.securityContext.runAsUser
        run_as_user = pod.spec.security_context.run_as_user if pod.spec.security_context else None
        print(f"  runAsUser: {run_as_user}")

        # 3. spec.hostPID
        host_pid = pod.spec.host_pid
        print(f"  hostPID: {host_pid}")

        # 4. spec.hostIPC
        host_ipc = pod.spec.host_ipc
        print(f"  hostIPC: {host_ipc}")

        # 5. spec.hostNetwork
        host_network = pod.spec.host_network
        print(f"  hostNetwork: {host_network}")

        # 6. spec.volumes
        volumes = [volume.name for volume in pod.spec.volumes]
        print(f"  Volumes: {volumes}")

        # 7. spec.serviceAccount
        service_account = pod.spec.service_account
        print(f"  Service Account: {service_account}")

        # 8. spec.serviceAccountName
        service_account_name = pod.spec.service_account_name
        print(f"  Service Account Name: {service_account_name}")

        # 9. spec.nodeName
        node_name = pod.spec.node_name
        print(f"  Node Name: {node_name}")

        for container in pod.spec.containers:
            print(f"  Container name: {container.name}")

            # 10. spec.containers.securityContext.capabilities.drop
            capabilities = container.security_context.capabilities if container.security_context else None
            drop_capabilities = capabilities.drop if capabilities else []
            print(f"    drop capabilities: {drop_capabilities}")

            # 11. spec.containers.securityContext.capabilities.add
            add_capabilities = capabilities.add if capabilities else []
            print(f"    add capabilities: {add_capabilities}")

            # 12. spec.containers.securityContext.runAsUser
            container_run_as_user = container.security_context.run_as_user if container.security_context else None
            print(f"    runAsUser: {container_run_as_user}")

            # 13. spec.containers.securityContext.allowPrivilegeEscalation
            allow_privilege_escalation = container.security_context.allow_privilege_escalation if container.security_context else None
            print(f"    Allow Privilege Escalation: {allow_privilege_escalation}")

            # 14. spec.containers.securityContext.privileged
            privileged = container.security_context.privileged if container.security_context else None
            print(f"    Privileged: {privileged}")

            # 15. spec.containers.ports.hostPort
            host_ports = [port.host_port for port in container.ports if port.host_port]
            print(f"    Host Ports: {host_ports}")

            # 16. spec.containers.readinessProbe
            readiness_probe = container.readiness_probe is not None
            print(f"    Readiness Probe: {readiness_probe}")

            # 17. spec.containers.image
            image = container.image
            print(f"    Image: {image}")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Kubernetes Pod Security Context and Parameters")
    parser.add_argument("--namespace", default="default", help="The namespace to inspect.")
    parser.add_argument("--kubeconfig", default="~/.kube/config", help="Path to the kubeconfig file.")
    args = parser.parse_args()

    get_pod_security_contexts_and_parameters(args.namespace, args.kubeconfig)