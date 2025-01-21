DROP TABLE IF EXISTS posts, authors;

CREATE TABLE authors (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    author_id INTEGER REFERENCES authors(id) NOT NULL,
    title TEXT  NOT NULL,
    content TEXT NOT NULL,
    created_at BIGINT NOT NULL
);

INSERT INTO authors (id, name) VALUES (0, 'Виктор Пелевин');
INSERT INTO posts (id, author_id, title, content, created_at) 
VALUES (
0, 
0, 
'О Любви', 
`Любовь в сущности, возникает в одиночестве, 
когда рядом нет её объекта, и направлена она 
не столько на того или ту, кого любишь, 
сколько на выстроенный умом образ, 
слабо связанный с оригиналом.`,
0
);