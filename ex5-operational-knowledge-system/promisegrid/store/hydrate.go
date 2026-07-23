package store

import "fmt"

// Intent: Keep authoritative frozen-envelope hydration substrate-owned so apps
// can reuse one durable replay rule while still choosing their own projections
// and workflows above it. Source: DI-lemor
func HydrateAuthoritativeEnvelopes[T any](
	cas *CASStore,
	records []T,
	envelopeCID func(T) string,
	manifestBase64 func(T) string,
	describe func(T) string,
	setEnvelopeBase64 func(*T, string),
) ([]T, error) {
	out := append([]T(nil), records...)
	for i := range out {
		base64Envelope, err := cas.AuthoritativeEnvelopeBase64(envelopeCID(out[i]), manifestBase64(out[i]))
		if err != nil {
			return nil, fmt.Errorf("load authoritative %s: %w", describe(out[i]), err)
		}
		setEnvelopeBase64(&out[i], base64Envelope)
	}
	return out, nil
}
