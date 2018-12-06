package defaults

import (
	"fmt"
	"testing"

	"github.com/kiegroup/kie-cloud-operator/pkg/controller/kieapp/constants"

	"github.com/kiegroup/kie-cloud-operator/pkg/apis/app/v1"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestLoadTrialEnvironment(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			logrus.Error(err.(error))
		}
	}()

	cr := &v1.KieApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "test-ns",
		},
		Spec: v1.KieAppSpec{
			Environment: "trial",
		},
	}

	env, err := GetEnvironment(cr, fake.NewFakeClient())
	assert.Equal(t, fmt.Sprintf("%s-kieserver-%d", cr.Name, len(env.Servers)-1), env.Servers[len(env.Servers)-1].DeploymentConfigs[0].Name)
	assert.Nil(t, err)
}

func TestLoadUnknownEnvironment(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			logrus.Error(err.(error))
		}
	}()

	cr := &v1.KieApp{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test-ns",
		},
		Spec: v1.KieAppSpec{
			Environment: "unknown",
		},
	}

	_, err := GetEnvironment(cr, fake.NewFakeClient())
	assert.Equal(t, fmt.Sprintf("envs/%s.yaml does not exist, '%s' KieApp not deployed", cr.Spec.Environment, cr.Name), err.Error())
}

func TestInaccessibleConfigMap(t *testing.T) {
	cr := &v1.KieApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "map-test",
			Namespace: "test-ns",
		},
		Spec: v1.KieAppSpec{
			Environment: "na",
		},
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-cm",
			Namespace: "test-ns",
		},
		Data: map[string]string{
			"test-key": "test-value",
		},
	}

	client := fake.NewFakeClient(cm)
	_, err := GetEnvironment(cr, client)
	assert.Equal(t, fmt.Sprintf("%s/%s ConfigMap not yet accessible, '%s' KieApp not deployed. Retrying... ", cr.Namespace, constants.ConfigMapPrefix, cr.Name), err.Error())
}

func TestMultipleServerDeployment(t *testing.T) {
	defer func() {
		err := recover()
		if err != nil {
			logrus.Error(err.(error))
		}
	}()

	cr := &v1.KieApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "test-ns",
		},
		Spec: v1.KieAppSpec{
			Environment:    "trial",
			KieDeployments: 6,
		},
	}

	env, err := GetEnvironment(cr, fake.NewFakeClient())
	assert.Equal(t, cr.Spec.KieDeployments, len(env.Servers))
	assert.Equal(t, fmt.Sprintf("%s-kieserver-%d", cr.Name, cr.Spec.KieDeployments-1), env.Servers[cr.Spec.KieDeployments-1].DeploymentConfigs[0].Name)
	assert.Nil(t, err)
}
