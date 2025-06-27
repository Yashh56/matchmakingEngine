package gameorchestrator

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os/exec"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

func CreateGamePod(ctx context.Context, client *kubernetes.Clientset, matchId string) (string, error) {
	podName := fmt.Sprintf("game-%s-%d", matchId[:8], time.Now().Unix())
	namespace := "default"

	// Create Pod
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace,
			Labels: map[string]string{
				"app":   "game-server",
				"match": matchId,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "game-server",
					Image: "jmalloc/echo-server:v0.3.7",
					Ports: []corev1.ContainerPort{
						{
							Name:          "http",
							ContainerPort: 8080,
						},
					},
					ReadinessProbe: &corev1.Probe{
						ProbeHandler: corev1.ProbeHandler{
							HTTPGet: &corev1.HTTPGetAction{
								Path: "/",
								Port: intstr.FromInt(8080),
							},
						},
						InitialDelaySeconds: 5,
						PeriodSeconds:       10,
					},
				},
			},
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}

	if _, err := client.CoreV1().Pods(namespace).Create(ctx, pod, metav1.CreateOptions{}); err != nil {
		return "", fmt.Errorf("❌ failed to create pod: %v", err)
	}

	// Wait for pod to be ready
	if err := waitForPodReady(ctx, client, namespace, podName, 60*time.Second); err != nil {
		return "", fmt.Errorf("❌ pod not ready: %v", err)
	}

	// Find local port for port-forwarding
	localPort, err := findAvailablePort()
	if err != nil {
		return "", fmt.Errorf("❌ failed to find available port: %v", err)
	}

	// Start port-forward in background
	if err := startPortForward(podName, namespace, localPort, 8080); err != nil {
		return "", fmt.Errorf("❌ failed to start port forward: %v", err)
	}

	// ⏳ Cleanup pod after 5 minutes
	go func() {
		time.Sleep(5 * time.Minute)
		log.Printf("🧹 Cleaning up pod %s", podName)
		_ = client.CoreV1().Pods(namespace).Delete(context.Background(), podName, metav1.DeleteOptions{})
	}()

	// Return clean WS URL
	return fmt.Sprintf("ws://localhost:%d/ws", localPort), nil
}

// --- Helpers ---

func waitForPodReady(ctx context.Context, client *kubernetes.Clientset, namespace, podName string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for pod to be ready")
		default:
			pod, err := client.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
			if err != nil {
				return err
			}
			for _, cond := range pod.Status.Conditions {
				if cond.Type == corev1.PodReady && cond.Status == corev1.ConditionTrue {
					return nil
				}
			}
			time.Sleep(2 * time.Second)
		}
	}
}

func findAvailablePort() (int, error) {
	for i := 0; i < 10; i++ {
		port := rand.Intn(10000) + 20000
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			ln.Close()
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available port found")
}

func startPortForward(podName, namespace string, localPort, remotePort int) error {
	cmd := exec.Command("kubectl", "port-forward",
		fmt.Sprintf("pod/%s", podName),
		fmt.Sprintf("%d:%d", localPort, remotePort),
		"-n", namespace)

	return cmd.Start() // silently starts in background
}
