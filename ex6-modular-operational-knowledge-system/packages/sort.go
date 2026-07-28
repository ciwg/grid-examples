package packages

import "slices"

func SortManifests(manifests []Manifest) {
	slices.SortFunc(manifests, func(left, right Manifest) int {
		if left.ID < right.ID {
			return -1
		}
		if left.ID > right.ID {
			return 1
		}
		return 0
	})
}
