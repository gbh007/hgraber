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
