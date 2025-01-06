# Tests 

Steps to create and ready match : 

- P1 creates new match with their account : 
  - [x] {{host}}/match/{{p1accountId}}
  - Set state: New
    - [x] Creates a match player
    - [x] Creates deck
    - [x] Creates 52 match cards
    - [x] Associate match cards with deck
  - Set state: Waiting
- [x] Get open matches for select match to join: {{host}}/open
- P2 joins :
  - [x] {{host}}/match/{{matchId}}/join/{{p2accountId}}
  - Set state: Determine
    - Determine first player
    - Set state: Deal
- Server: Deal hand to each player
  - [x] {{host}}/match/{{matchId}}/player/{{playerId}}/kitty
  - [x] Set state: Discard
- Await: Players to submit cards for kitty
  - [x]  {{host}}/match/{{matchId}}/player/{{playerId}}/play
  - [x] Set state: PlayOwn
- Await: P1 submit match card id 
  - [x] {{host}}/match/{{matchId}}/player/{{playerId}}/play
  - [x] Set state: PlayOpponent
- Await: P2 submit match card id
  - Set state: Count 
- Server: Update players score
  - Set state: Kitty 
- Server: Update dealer with kitty score 
//cawa-man 

