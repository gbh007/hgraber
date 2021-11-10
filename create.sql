-- таблица тайтлов (основная информация о тайтле)
CREATE TABLE IF NOT EXISTS titles(
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  url TEXT NOT NULL UNIQUE,
  page_count INT NOT NULL DEFAULT 0,
  creation_time TIMESTAMP,
  loaded BOOL DEFAULT FALSE,
  parsed_pages BOOL DEFAULT FALSE,
  parsed_tags BOOL DEFAULT FALSE,
  parsed_authors BOOL DEFAULT FALSE,
  parsed_characters BOOL DEFAULT FALSE,
  parsed_languages BOOL DEFAULT FALSE,
  parsed_categories BOOL DEFAULT FALSE,
  parsed_parodies BOOL DEFAULT FALSE,
  parsed_groups BOOL DEFAULT FALSE
);
-- таблица страниц произведения
CREATE TABLE IF NOT EXISTS pages(
  title_id INTEGER NOT NULL REFERENCES titles(id) ON UPDATE CASCADE ON DELETE CASCADE,
  ext TEXT NOT NULL,
  url TEXT NOT NULL,
  page_number INT NOT NULL,
  success BOOL DEFAULT FALSE,
  PRIMARY KEY(title_id, page_number)
);
-- таблица мета информации
CREATE TABLE IF NOT EXISTS meta(
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  type TEXT NOT NULL,
  UNIQUE(name, type)
);
-- линковка таблицы мета информации на тайтлы
CREATE TABLE IF NOT EXISTS link_meta_titles(
  meta_id INTEGER REFERENCES meta(id) ON UPDATE CASCADE ON DELETE CASCADE,
  title_id INTEGER REFERENCES titles(id) ON UPDATE CASCADE ON DELETE CASCADE,
  type TEXT NOT NULL
);