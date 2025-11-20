package util

import (
	"strings"
)

// createNormalizedLabels returns a map with multiple variants for each label key.
// Priority rules:
// 1) keep any explicit keys from the original labels (don't overwrite them)
// 2) add variants only if they do not exist yet
//
// Variants generated for a key like "kubernetes.pod_name":
// - "kubernetes.pod_name" (original)
// - "kubernetes.pod.name" (underscore -> dot)
// - "kubernetes_pod_name" (dot -> underscore)
// - "kubernetes_pod.name" (mix variants)
//
// This covers cases like:
// - kubernetes.pod.name
// - kubernetes.pod_name
// - kubernetes_pod.name
// - kubernetes_pod_name
func createNormalizedLabels(labels map[string]string) map[string]string {
	out := make(map[string]string, len(labels)*4)

	// copy originals first (highest priority)
	for k, v := range labels {
		out[k] = v
	}

	// helper to add variant only if key not already present
	addIfMissing := func(key, val string) {
		if key == "" {
			return
		}
		if _, exists := out[key]; !exists {
			out[key] = val
		}
	}

	for k, v := range labels {
		// variant 1: dots -> underscores
		dotToUnderscore := strings.ReplaceAll(k, ".", "_")
		addIfMissing(dotToUnderscore, v)

		// variant 2: underscores -> dots
		underscoreToDot := strings.ReplaceAll(k, "_", ".")
		addIfMissing(underscoreToDot, v)

		// variant 3: first replace dots->underscores then underscores->dots (mix), to cover mixed cases
		// e.g. kubernetes.pod_name -> kubernetes_pod.name
		mix1 := strings.ReplaceAll(dotToUnderscore, "_", ".")
		addIfMissing(mix1, v)

		// variant 4: first underscores->dots then dots->underscores
		// e.g. kubernetes_pod.name -> kubernetes.pod_name
		mix2 := strings.ReplaceAll(underscoreToDot, ".", "_")
		addIfMissing(mix2, v)
	}

	return out
}
