INSERT INTO
    cards ("value", "suit", "art")
VALUES
    ('Ace', 'Spades', 'meow'),
    ('Two', 'Spades', 'meow'),
    ('Three', 'Spades', 'meow'),
    ('Four', 'Spades', 'meow'),
    ('Five', 'Spades', 'meow'),
    ('Six', 'Spades', 'meow'),
    ('Seven', 'Spades', 'meow'),
    ('Eight', 'Spades', 'meow'),
    ('Nine', 'Spades', 'meow'),
    ('Ten', 'Spades', 'meow'),
    ('Jack', 'Spades', 'meow'),
    ('Queen', 'Spades', 'meow'),
    ('King', 'Spades', 'meow'),
    ('Ace', 'Clubs', 'meow'),
    ('Two', 'Clubs', 'meow'),
    ('Three', 'Clubs', 'meow'),
    ('Four', 'Clubs', 'meow'),
    ('Five', 'Clubs', 'meow'),
    ('Six', 'Clubs', 'meow'),
    ('Seven', 'Clubs', 'meow'),
    ('Eight', 'Clubs', 'meow'),
    ('Nine', 'Clubs', 'meow'),
    ('Ten', 'Clubs', 'meow'),
    ('Jack', 'Clubs', 'meow'),
    ('Queen', 'Clubs', 'meow'),
    ('King', 'Hearts', 'meow'),
    ('Ace', 'Hearts', 'meow'),
    ('Two', 'Hearts', 'meow'),
    ('Three', 'Hearts', 'meow'),
    ('Four', 'Hearts', 'meow'),
    ('Five', 'Hearts', 'meow'),
    ('Six', 'Hearts', 'meow'),
    ('Seven', 'Hearts', 'meow'),
    ('Eight', 'Hearts', 'meow'),
    ('Nine', 'Hearts', 'meow'),
    ('Ten', 'Hearts', 'meow'),
    ('Jack', 'Hearts', 'meow'),
    ('Queen', 'Hearts', 'meow'),
    ('King', 'Hearts', 'meow'),
    ('Ace', 'Diamonds', 'meow'),
    ('Two', 'Diamonds', 'meow'),
    ('Three', 'Diamonds', 'meow'),
    ('Four', 'Diamonds', 'meow'),
    ('Five', 'Diamonds', 'meow'),
    ('Six', 'Diamonds', 'meow'),
    ('Seven', 'Diamonds', 'meow'),
    ('Eight', 'Diamonds', 'meow'),
    ('Nine', 'Diamonds', 'meow'),
    ('Ten', 'Diamonds', 'meow'),
    ('Jack', 'Diamonds', 'meow'),
    ('Queen', 'Diamonds', 'meow'),
    ('King', 'Diamonds', 'meow');

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