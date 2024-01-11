package rendering

import "app/internal/domain/hgraber"

func attributeDisplayName(code hgraber.Attribute) string {
	switch code {
	case hgraber.AttrAuthor:
		return "Авторы"
	case hgraber.AttrCategory:
		return "Категории"
	case hgraber.AttrCharacter:
		return "Персонажи"
	case hgraber.AttrGroup:
		return "Группы"
	case hgraber.AttrLanguage:
		return "Языки"
	case hgraber.AttrParody:
		return "Пародии"
	case hgraber.AttrTag:
		return "Тэги"
	default:
		return string(code)
	}
}

func attributeOrder(code hgraber.Attribute) int {
	switch code {
	case hgraber.AttrAuthor:
		return 3
	case hgraber.AttrCategory:
		return 2
	case hgraber.AttrCharacter:
		return 4
	case hgraber.AttrGroup:
		return 5
	case hgraber.AttrLanguage:
		return 6
	case hgraber.AttrParody:
		return 7
	case hgraber.AttrTag:
		return 1
	default:
		return 999
	}
}
