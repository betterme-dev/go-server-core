package image

type (
	// DB Representation
	ImageFile struct {
		ImageID uint32 `db:"image_id"`
		URL     string `db:"url"`
	}
	// DTO
	Image struct {
		ID uint32
	}
	// Other
	Quality int64
)

const (
	ProviderS3      = 5
	StatusProcessed = 1
)

var (
	QualityOriginal Quality = 0
	QualityBasic    Quality = 1
	QualityHD       Quality = 2
	QualityFullHD   Quality = 3
)
