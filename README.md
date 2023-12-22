# HGraber

Граббер для манги.

## Возможности

1. Скачивать мангу с определенных сайтов
2. Просматривать скачанную мангу
3. Экспортировать мангу в формат для читалок

## Поддерживаемые БД

1. JSON
2. PostgreSQL

## Планы развития на 4-ю версию

1. Отказаться от рейтина (пока под вопросом)
   - Основная цель граббера закачивать данные, а не оценивать
2. Переработать API
   - Для списка тайтлов уменьшить количество информации
   - Произвести ренейминг
3. Упростить/переработать просмотр
   - На списке тайтлов нет необходимости показывать всю информацию
     - Первично теги и прочие метаданные достаточно отображать только количественно
4. Добавить версионирование и миграции в jdb
5. Рефакторинг фронта, для изоляции страниц между собой
   - Отказаться от разделяемого кода, поскольку после рефакторинга API потребность в нем должна исчезнуть
6. Выгрузку файлов сделать через воркера, с мониторингом
7. Возможно перевести фронт на VueJS 2/3
   - Только в случает отсутствия необходимости сборки его Webpack/Vite

## TODO

1. Рефакторинг экспорта файлов
   - Сделать выгрузку ассинхронной (через воркера)
   - В рамках 4-й версии
2. Покрытие кода тестами

## Альтернативный клиент

Также есть [альтернативный клиент](https://gitlab.com/gbh007/hgraber_ui) для этого сервера на flutter

Этот клиент поддерживает функциональность частично
