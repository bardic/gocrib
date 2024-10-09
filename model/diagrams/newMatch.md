```mermaid
 ---
    title: New Match
---
flowchart TD
    CI(Client Input 'n') --> MS(Client Match Service)
    MS --> SR(Server Router)
    SR -- POST --> SM(Server Match)
    SM --> DB[(DB)]
    DB --> SM
    SM -- model.GameMatch --> MS
```
