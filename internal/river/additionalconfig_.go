//go:build !trustcenter

package river

// AdditionalWorkers is an empty struct when trustcenter build tag is not present
type AdditionalWorkers struct {
}
