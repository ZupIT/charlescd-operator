package event

import "sigs.k8s.io/controller-runtime/pkg/client"

func diff(old, new client.Object) string {
	patch := client.MergeFrom(old)
	data, _ := patch.Data(new)
	return string(data)
}
