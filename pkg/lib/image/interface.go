package image

type (
	Provider interface {
		ImagesByIds(ids []uint, quality Quality) ([]ImageFile, error)
	}
)
