package kiosk

import (
	"github.com/flanksource/karina/pkg/constants"
	"github.com/flanksource/karina/pkg/platform"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const Namespace = "platform-system"
const Certs = "kiosk-tls"

func Deploy(p *platform.Platform) error {
	if p.Kiosk.IsDisabled() {
		return p.DeleteSpecs(constants.PlatformSystem, "kiosk.yaml")
	}

	if !p.HasSecret(Namespace, Certs) {
		secret := p.NewSelfSigned("kiosk.platform-system.svc").AsTLSSecret()
		secret["ca.crt"] = secret["tls.crt"]
		if err := p.Apply(Namespace, &v1.Secret{
			TypeMeta: metav1.TypeMeta{APIVersion: "v1", Kind: "Secret"},
			ObjectMeta: metav1.ObjectMeta{
				Name:      Certs,
				Namespace: Namespace,
				Annotations: map[string]string{
					"cert-manager.io/allow-direct-injection": "true",
				},
			},
			Data: secret,
			Type: v1.SecretTypeTLS,
		}); err != nil {
			return err
		}
	}

	return p.ApplySpecs(constants.PlatformSystem, "kiosk.yaml")
}
