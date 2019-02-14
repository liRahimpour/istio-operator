/*
Copyright 2019 Banzai Cloud.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sidecarinjector

import (
	istiov1beta1 "github.com/banzaicloud/istio-operator/pkg/apis/operator/v1beta1"
	"github.com/banzaicloud/istio-operator/pkg/resources/templates"
	apiv1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (r *Reconciler) serviceAccount(owner *istiov1beta1.Config) runtime.Object {
	return &apiv1.ServiceAccount{
		ObjectMeta: templates.ObjectMeta(serviceAccountName, sidecarInjectorLabels, owner),
	}
}

func (r *Reconciler) clusterRole(owner *istiov1beta1.Config) runtime.Object {
	return &rbacv1.ClusterRole{
		ObjectMeta: templates.ObjectMetaClusterScope(clusterRoleName, sidecarInjectorLabels, owner),
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"configmaps"},
				Resources: []string{"*"},
				Verbs:     []string{"get", "watch", "list"},
			},
			{
				APIGroups: []string{"admissionregistration.k8s.io"},
				Resources: []string{"mutatingwebhookconfigurations"},
				Verbs:     []string{"get", "list", "watch", "patch"},
			},
		},
	}
}

func (r *Reconciler) clusterRoleBinding(owner *istiov1beta1.Config) runtime.Object {
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: templates.ObjectMetaClusterScope(clusterRoleBindingName, sidecarInjectorLabels, owner),
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			APIGroup: "rbac.authorization.k8s.io",
			Name:     clusterRoleName,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      serviceAccountName,
				Namespace: owner.Namespace,
			},
		},
	}
}
