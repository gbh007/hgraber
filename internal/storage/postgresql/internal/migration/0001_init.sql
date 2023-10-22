CREATE TABLE
    books(
        id SERIAL PRIMARY KEY,
        name TEXT,
        url TEXT UNIQUE,
        page_count INT,
        create_at TIMESTAMPTZ NOT NULL,
        rate INT NOT NULL DEFAULT 0
    );

CREATE TABLE
    pages(
        book_id INT NOT NULL REFERENCES books(id) ON UPDATE CASCADE ON DELETE CASCADE,
        page_number INT NOT NULL,
        ext TEXT NOT NULL,
        url TEXT NOT NULL,
        success BOOL NOT NULL DEFAULT FALSE,
        create_at TIMESTAMPTZ NOT NULL,
        load_at TIMESTAMPTZ,
        rate INT NOT NULL DEFAULT 0,
        PRIMARY KEY(book_id, page_number)
    );

CREATE TABLE attributes( code TEXT PRIMARY KEY );

INSERT INTO attributes(code)
VALUES ('tag'), ('author'), ('character'), ('language'), ('category'), ('parody'), ('group');

CREATE TABLE
    book_attributes(
        book_id INT NOT NULL REFERENCES books(id) ON UPDATE CASCADE ON DELETE CASCADE,
        attr TEXT NOT NULL REFERENCES attributes(code) ON UPDATE CASCADE ON DELETE CASCADE,
        value TEXT NOT NULL
    );

CREATE TABLE
    book_attributes_parsed(
        book_id INT NOT NULL REFERENCES books(id) ON UPDATE CASCADE ON DELETE CASCADE,
        attr TEXT NOT NULL REFERENCES attributes(code) ON UPDATE CASCADE ON DELETE CASCADE,
        parsed BOOL NOT NULL,
        UNIQUE(book_id, attr)
    );