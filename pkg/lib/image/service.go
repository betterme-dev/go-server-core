package image

import (
	log "github.com/sirupsen/logrus"
)

type Service struct {
	Repository Repository
}

func NewService() Service {
	return Service{Repository: NewRepository()}
}

func (s Service) ImagesByIds(ids []uint64, quality Quality) ([]ImageFile, error) {
	res, err := s.Repository.ImagesByIds(ids, ProviderS3, quality)
	if err != nil {
		return nil, err
	}

	order := []Quality{
		QualityFullHD,
		QualityHD,
		QualityBasic,
		QualityOriginal,
	}

	skipQuality := true
	for _, orderQuality := range order {
		if skipQuality && quality == orderQuality {
			skipQuality = false
		}
		if skipQuality {
			continue
		}

		notFoundIDs := notFoundImages(ids, res)
		if len(notFoundIDs) == 0 {
			return res, nil
		}

		fallbackImages, err := s.Repository.ImagesByIds(notFoundIDs, ProviderS3, orderQuality)
		if err != nil {
			return nil, err
		}

		fallbackImagesCount := len(fallbackImages)
		if fallbackImagesCount != 0 {
			log.Debugf("Some images not found (%d) with quality %d. Fallback to %d quality", fallbackImagesCount, quality, orderQuality)
		}
		res = append(res, fallbackImages...)
	}

	return res, nil
}

func notFoundImages(ids []uint64, images []ImageFile) []uint64 {
	result := make([]uint64, 0)
	if len(images) < len(ids) {
		for _, id := range ids {
			found := false
			for _, image := range images {
				if id == uint64(image.ImageID) {
					found = true
				}
			}
			if !found {
				result = append(result, id)
			}
		}
	}

	return result
}
