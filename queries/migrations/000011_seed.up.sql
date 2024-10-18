INSERT INTO
    cards ("value", "suit", "art")
VALUES
    (1, 0, 'meow'),
    (2, 0, 'meow'),
    (3, 0, 'meow'),
    (4, 0, 'meow'),
    (5, 0, 'meow'),
    (6, 0, 'meow'),
    (7, 0, 'meow'),
    (8, 0, 'meow'),
    (9, 0, 'meow'),
    (10, 0, 'meow'),
    (11, 0, 'meow'),
    (12, 0, 'meow'),
    (13, 0, 'meow'),
    (1, 1, 'meow'),
    (2, 1, 'meow'),
    (3, 1, 'meow'),
    (4, 1, 'meow'),
    (5, 1, 'meow'),
    (6, 1, 'meow'),
    (7, 1, 'meow'),
    (8, 1, 'meow'),
    (9, 1, 'meow'),
    (10, 1, 'meow'),
    (11, 1, 'meow'),
    (12, 1, 'meow'),
    (13, 0, 'meow'),
    (1, 2, 'meow'),
    (2, 2, 'meow'),
    (3, 2, 'meow'),
    (4, 2, 'meow'),
    (5, 2, 'meow'),
    (6, 2, 'meow'),
    (7, 2, 'meow'),
    (8, 2, 'meow'),
    (9, 2, 'meow'),
    (10, 2, 'meow'),
    (11, 2, 'meow'),
    (13, 0, 'meow'),
    (1, 3, 'meow'),
    (2, 3, 'meow'),
    (3, 3, 'meow'),
    (4, 3, 'meow'),
    (5, 3, 'meow'),
    (6, 3, 'meow'),
    (7, 3, 'meow'),
    (8, 3, 'meow'),
    (9, 3, 'meow'),
    (10, 3, 'meow'),
    (11, 3, 'meow'),
    (12, 3, 'meow'),
    (13, 0, 'meow');

INSERT INTO
    accounts("name")
VALUES
    ('test1'),
    ('test2');

INSERT INTO
    player(
        "accountid",
        "play",
        "hand",
        "kitty",
        "score",
        "isready",
        "art"
    )
VALUES
    (
        1,
        ARRAY [] :: integer [],
        ARRAY [] :: integer [],
        ARRAY [] :: integer [],
        0,
        false,
        'default.png'
    ),
    (
        2,
        ARRAY [] :: integer [],
        ARRAY [] :: integer [],
        ARRAY [] :: integer [],
        0,
        false,
        'default.png'
    )