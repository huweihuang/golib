package db

import "gorm.io/gorm"

func Paginate(page int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		switch {
		case pageSize > 500:
			pageSize = 500
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func PaginatedFind(db *gorm.DB, page, pageSize int, target interface{}) (count int64, err error) {
	err = db.Count(&count).Error
	if err != nil {
		return 0, err
	}
	err = db.Scopes(Paginate(page, pageSize)).Find(target).Error
	return count, err
}

func ParamsToQuery(queryMap map[string]interface{}) (queryStr string, args []interface{}) {
	for key, value := range queryMap {
		queryStr = AppendQuery(queryStr, key)
		args = append(args, value)
	}
	return queryStr, args
}

func AppendQuery(query, newQuery string) string {
	if query != "" {
		return query + " AND " + newQuery
	}
	return newQuery
}
