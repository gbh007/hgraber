package base

import (
	"errors"
)

var (
	ErrNotAuth        = errors.New("необходима авторизация")
	ErrForbidden      = errors.New("нет прав доступа")
	ErrNotFound       = errors.New("данные не найдены")
	ErrParseData      = errors.New("некорректный формат данных")
	ErrDeleteData     = errors.New("невозможно удалить данные")
	ErrCreateData     = errors.New("невозможно создать данные")
	ErrAppendData     = errors.New("невозможно добавить данные")
	ErrUpdateData     = errors.New("невозможно изменить данные")
	ErrProcessingData = errors.New("невозможно обработать данные")
	ErrPanicDetected  = errors.New("нарушение потока выполнения запроса")
)
