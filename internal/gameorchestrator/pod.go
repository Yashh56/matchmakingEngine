package gameorchestrator

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreateGamePod(ctx context.Context, client *kubernetes.Clientset, matchId string) (string, error) {

	podName := "game-" + matchId[:8]

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: podName,
			Labels: map[string]string{
				"match": matchId,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "game",
					Image: "nginx",
					Ports: []corev1.ContainerPort{
						{
							ContainerPort: 8080,
						},
					},
				},
			},
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}

	_, err := client.CoreV1().Pods("Matchmaking-Engine").Create(ctx, pod, metav1.CreateOptions{})

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s.default.svc.cluster.local:8080", podName), nil

}
