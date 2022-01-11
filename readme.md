kubectl run -it cqlsh --image cassandra:3.11 -- /bin/bash
kubectl exec -it cqlsh -- /bin/bash

cqlsh cassandra-0.cassandra

INSERT INTO test.games (game_id , game_date , game_team1 , game_team2 , stadium, stadium_capacity ) VALUES ( now(), '2011-02-03', 'Niemcy', 'Brazylia', now(), 80000);
