<?php

for($i = 0 ; $i <50; $i++){
    $cmd = "ab -n 500000 -c 10 -p a.json 'http://localhost:8000/push?queue_id=QUEUE_{$i}' >> bench.log 2>&1 &";
    echo $cmd;
    shell_exec($cmd);
}
