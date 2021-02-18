package image

import (
	"github.com/betterme-dev/go-server-core/pkg/db"
	"github.com/doug-martin/goqu/v9"
)

const TableName = "image_file"

type Repository struct {
	table string
	db    *goqu.Database
}

func NewRepository() Repository {
	return Repository{
		table: TableName,
		db:    db.Goqu(),
	}
}

func (r Repository) ImagesByIds(ids []uint64, provider uint, quality Quality) ([]ImageFile, error) {
	var images []ImageFile
	if len(ids) == 0 {
		return images, nil
	}

	err := r.db.
		From(r.table).
		Select(
			goqu.C("image_id").Distinct(),
			"url",
		).
		Where(
			goqu.C("image_id").In(ids),
			goqu.C("provider").Eq(provider),
			goqu.C("quality").Eq(quality),
			goqu.C("status").Eq(StatusProcessed),
		).ScanStructs(&images)
	if err != nil {
		return nil, err
	}

	return images, nil
}
