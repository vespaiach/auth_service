package mysql

import "github.com/vespaiach/auth_service/pkg/share"

func getOrderDirection(i share.Direction) string {
	switch i {
	case share.Ascendant:
		return "ASC"
	case share.Descendant:
		return "DESC"
	default:
		return ""
	}
}
