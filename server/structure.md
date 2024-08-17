# Gameplay

- cards
  - value
  - suit
  - orig_owner
  - curr_owner
  - state
  - art
- players
  - hand []Card
  - kitty []Card
  - score
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
  - []players
  - art

# Server/DB Structs

- cards
- player
- match
- lobby
- history
- chat
- account
