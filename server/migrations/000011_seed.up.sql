INSERT INTO
    cards ("value", "suit", "art")
VALUES
    (0, 0, 'meow'),
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
    (0, 1, 'meow'),
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
    (0, 2, 'meow'),
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
    (12, 2, 'meow'),
    (0, 3, 'meow'),
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
    (12, 3, 'meow');

INSERT INTO
    gameplaycards ("cardid", "origowner", "currowner", "state")
VALUES
    (0, 0, 0, 0),
    (1, 0, 0, 0),
    (2, 0, 0, 0),
    (3, 0, 0, 0),
    (4, 0, 0, 0),
    (5, 0, 0, 0),
    (6, 0, 0, 0),
    (7, 0, 0, 0),
    (8, 0, 0, 0),
    (9, 0, 0, 0),
    (10, 0, 0, 0),
    (11, 0, 0, 0),
    (12, 0, 0, 0),
    (13, 0, 0, 0),
    (14, 0, 0, 0),
    (15, 0, 0, 0),
    (16, 0, 0, 0),
    (17, 0, 0, 0),
    (18, 0, 0, 0),
    (19, 0, 0, 0),
    (20, 0, 0, 0),
    (21, 0, 0, 0),
    (22, 0, 0, 0),
    (23, 0, 0, 0),
    (24, 0, 0, 0),
    (25, 0, 0, 0),
    (26, 0, 0, 0),
    (27, 0, 0, 0),
    (28, 0, 0, 0),
    (29, 0, 0, 0),
    (30, 0, 0, 0),
    (31, 0, 0, 0),
    (32, 0, 0, 0),
    (33, 0, 0, 0),
    (34, 0, 0, 0),
    (35, 0, 0, 0),
    (36, 0, 0, 0),
    (37, 0, 0, 0),
    (38, 0, 0, 0),
    (39, 0, 0, 0),
    (40, 0, 0, 0),
    (41, 0, 0, 0),
    (42, 0, 0, 0),
    (43, 0, 0, 0),
    (44, 0, 0, 0),
    (45, 0, 0, 0),
    (46, 0, 0, 0),
    (47, 0, 0, 0),
    (48, 0, 0, 0),
    (49, 0, 0, 0),
    (50, 0, 0, 0),
    (51, 0, 0, 0);

INSERT INTO
    deck ("matchid", "cards")
VALUES
    (
        0,
        ARRAY [ 0,
        1,
        2,
        3,
        4,
        5,
        6,
        7,
        8,
        9,
        10,
        11,
        12,
        13,
        14,
        15,
        16,
        17,
        18,
        19,
        20,
        21,
        22,
        23,
        24,
        25,
        26,
        27,
        28,
        29,
        30,
        31,
        32,
        33,
        34,
        35,
        36,
        37,
        38,
        39,
        40,
        41,
        42,
        43,
        44,
        45,
        46,
        47,
        48,
        49,
        50,
        51
        ]
    );

INSERT INTO
    lobby(
        "accountids",
        "privatematch",
        "elorangemin",
        "elorangemax"
    )
VALUES
    (ARRAY [ 0, 1 ], false, 0, 10);

INSERT INTO
    accounts("name")
VALUES
    ('test1'),
    ('test2');

INSERT INTO
    match(
        "lobbyid",
        "deckid",
        "playerids",
        "cardsinplay",
        "cutgamecardid",
        "currentplayerturn",
        "turnpasstimestamps",
        "art"
    )
VALUES
    (
        1,
        1,
        ARRAY [ 1, 2 ],
        ARRAY [] :: integer [],
        12,
        1,
        ARRAY [] :: varchar(255) [],
        'meow'
    );

INSERT INTO
    player("hand", "kitty", "score", "art")
VALUES
    (ARRAY [0,1,2,3], ARRAY [0,1,2,3], 0, 'meow'),
    (ARRAY [0,1,2,3], ARRAY [0,1,2,3], 0, 'meow');