- Have 2 endpoints (master & slave) that receive the same http webhooks
    Assume hosted at 2 different cloud providers
    Assume 2 splitted databases, direct joins are not possible
    They can talk over http, but also expose database connection (though will be 2 separate d

- Webhooks are stored in a DB for later processing
    Don’t need to worry about processing, just store
    Up to you which DB engine you use

- Master checks every 2 min if it received all webhooks 
    that slave received in interval [-4min,-2min ago]
    Make an efficient comparison algorithm that doesn’t kill the server or the network
    (we receive up to 100 webhooks/second → 12k webhooks/2 minutes)
    Hooks are unique when body is the same + max 1 min diff in received time
    So “forgotten webhooks” = slave[-4min, -2min] - master[-5min, -1min]
    E.g.: slave received body “{a:b}” at 19:45:59 UTC, 
    we should not insert into master when it has a body “{a:b}” between 19:44:59 and 19:46:59 UTC.

- Endpoints that receive json bodies
    [POST] master.probackup.io/v1/webhooks
    [POST] slave.probackup.io/v1/webhooks
    The bodies are always json but do not have a fixed structure
    Response “RECEIVED” is ok 

- Store webhooks async in batch
    Do not store events 1 by 1 (synchronously in the web request) 
    but accumulate them and store in batch

- Database
    Preferably use AWS DynamoDB
