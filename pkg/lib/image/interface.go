package image

type (
	Provider interface {
		ImagesByIds(ids []uint64, quality Quality) ([]ImageFile, error)
	}
)
