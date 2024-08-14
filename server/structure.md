# Gameplay

- cards
  - value
  - suit
  - orig_owner
  - curr_owner
  - state
  - art
- hand
  - maxSize
  - []Cards
- players
  - hand
  - kitty
  - score
  - art
- board
  - id
  - []players
  - art
- lobby
  - []accounts
  - id
  - creatationDate
  - private
  - eloRangeMin
  - eloRangeMax
- match
  - lobbyId
  - currentPlayerTurn
  - []TurnPassTimestamps
  - boardId

# Server/DB Structs

- cards
- player
- match
- lobby
- history
- chat
- account
