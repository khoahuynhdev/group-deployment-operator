package controller

import (
	"context"
	"fmt"

	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	multideploymentv1 "xplatform.api/group_deployment/api/v1"
)

func (r *GroupDeploymentReconciler) EnsureDeployment(ctx context.Context, gd *multideploymentv1.GroupDeployment) error {
	// check blue deployment
	log := log.FromContext(ctx)
	log.Info("Ensure blueDeployment", "blueDeployment", gd.Spec.BlueDeployment)
	blueDeployment := &appv1.Deployment{}
	blueDeploymentName := fmt.Sprintf("%s-blue-deployment", gd.Name)
	err := r.Get(ctx, client.ObjectKey{Name: blueDeploymentName}, blueDeployment)
	if err != nil {
		// if not found create one
		if apierrors.IsNotFound(err) {
			blueDeployment = &appv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      blueDeploymentName,
					Namespace: gd.Namespace,
				},
				Spec: appv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app": blueDeploymentName,
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app": blueDeploymentName,
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Image: gd.Spec.BlueDeployment,
									Name:  "blue-deployment",
								},
							},
						},
					},
				},
			}
			if err = r.Create(ctx, blueDeployment); err != nil {
				log.Error(err, "failed to create bluedeployment")
				return err
			}
		}
		log.Info("blueDeployment is available")

	}

	log.Info("Ensure blueDeployment", "greenDeployment", gd.Spec.GreenDeployment)
	greenDeployment := &appv1.Deployment{}
	greenDeploymentName := fmt.Sprintf("%s-green-deployment", gd.Name)
	err = r.Get(ctx, client.ObjectKey{Name: greenDeploymentName}, greenDeployment)
	if err != nil {
		// if not found create one
		if apierrors.IsNotFound(err) {
			greenDeployment = &appv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      greenDeploymentName,
					Namespace: gd.Namespace,
				},
				Spec: appv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app": greenDeploymentName,
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app": greenDeploymentName,
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Image: gd.Spec.GreenDeployment,
									Name:  "green-deployment",
								},
							},
						},
					},
				},
			}
			if err = r.Create(ctx, greenDeployment); err != nil {
				log.Error(err, "failed to create greenDeployment")
				return err
			}
		}
		log.Info("greenDeployment is available")
	}
	return nil
}
