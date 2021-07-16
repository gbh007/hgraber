-- таблица тайтлов (основная информация о тайтле)
CREATE TABLE IF NOT EXISTS titles(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    page_count INT NOT NULL DEFAULT 0,
    creation_time TIMESTAMP,
    loaded BOOL DEFAULT FALSE,
    parsed_pages BOOL DEFAULT FALSE,
    parsed_tags BOOL DEFAULT FALSE,
    parsed_authors BOOL DEFAULT FALSE,
    parsed_characters BOOL DEFAULT FALSE
);

-- таблица страниц произведения
CREATE TABLE IF NOT EXISTS pages(
    title_id INTEGER NOT NULL REFERENCES titles(id) ON UPDATE CASCADE ON DELETE CASCADE,
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    page_number INT NOT NULL,
    success BOOL DEFAULT FALSE
);

-- таблица тегов
CREATE TABLE IF NOT EXISTS tags(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE
);

-- таблица авторов
CREATE TABLE IF NOT EXISTS authors(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE
);

-- таблица персонажей
CREATE TABLE IF NOT EXISTS characters(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE
);

-- линковка таблицы тегов на тайтлы
CREATE TABLE IF NOT EXISTS link_tags_titles(
    tag_id INTEGER REFERENCES tags(id) ON UPDATE CASCADE ON DELETE CASCADE,
    title_id INTEGER REFERENCES titles(id) ON UPDATE CASCADE ON DELETE CASCADE
);

-- линковка таблицы авторов на тайтлы
CREATE TABLE IF NOT EXISTS link_authors_titles(
    author_id INTEGER REFERENCES authors(id) ON UPDATE CASCADE ON DELETE CASCADE,
    title_id INTEGER REFERENCES titles(id) ON UPDATE CASCADE ON DELETE CASCADE
);

-- линковка таблицы персонажей на тайтлы
CREATE TABLE IF NOT EXISTS link_characters_titles(
    character_id INTEGER REFERENCES characters(id) ON UPDATE CASCADE ON DELETE CASCADE,
    title_id INTEGER REFERENCES titles(id) ON UPDATE CASCADE ON DELETE CASCADE
);

