<?php

for($i = 0 ; $i <10; $i++){
    $cmd = "ab -n 5000000 -c 10 -p a.json 'http://localhost:8000/push?topic_id=QUEUE_{$i}' >> bench.log 2>&1 &";
    shell_exec($cmd);
    $cmd = "ab -n 5000000 -c 10 -p a.json 'http://localhost:8000/pop?topic_id=QUEUE_{$i}&target_count=10000&target_size=-1' >> bench.log 2>&1 &";
    shell_exec($cmd);
}
