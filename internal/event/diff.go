package event

import "sigs.k8s.io/controller-runtime/pkg/client"

func diff(oldest, newest client.Object) string {
	patch := client.MergeFrom(oldest)
	data, _ := patch.Data(newest)
	return string(data)
}
