# Tests 

Steps to create and ready match : 

- P1 creates new match with their account : 
  - [x] {{host}}/match/{{p1accountId}}
  - Set state: Create
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
  - Set state: Discard {{host}}/match/{{matchId}}/player/{{playerId}}/discard
- Await: Players to submit cards for kitty
  - Set state: P1 Player {{host}}/match/{{matchId}}/player/{{playerId}}/play
- Await: P1 submit match card id 
  - Set state: P2 Play {{host}}/match/{{matchId}}/player/{{playerId}}/play
- Await: P2 submit match card id
  - Set state: Count 
- Server: Update players score
  - Set state: Kitty 
- Server: Update dealer with kitty score 
//cawa-man 

