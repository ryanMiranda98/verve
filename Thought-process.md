1. The sample application is a simple Golang API containing one endpoint:
   `GET /api/verve/accept?id=1&endpoint=http://xyz.com`
   Here, the query parameter 'id' is mandatory and the 'endpoint' query parameter is optional.

2. The API returns status 200 OK when there is an integer id provided and returns status 400 BAD REQUEST when an invalid id/no id is provided.

3. The challenge requested that the unique request count has to be either logged (initially) to a file, or pushed to distributed streaming service (extension 3). To faciliate either options, I had decided to create a writer field in the APIServer struct, which is an io.WriterCloser interface.

> The message queue struct MQ, is an interface that consists of 3 methods i.e. Connect(), Close(), Write(). This is done to provide abstraction and extensibility to other message queue implementations in the future.
> Moreover, the Write() and Close() methods here are also a part of the io.WriterCloser interface and thus, can be used in the APIServer struct as a logging/writing mechanism.

> In case, we want to switch back to writing to a file, we can just pass in the \*os.File object. This can be seen in server.go - line 58 (commented out).

4. Initially, I had decided to keep a global variable counter using sync.Map to keep track of all the unique requests. The unique request tracking logic was implemented in the middleware and was concurrent-safe. However, based on the expectation of extension#2 in the challenge, this would not be easy to sync between two or more instances behind a load balancer.
> For that reason, I decided to use Redis (an in-memory and single-threaded) DB for quick access and concurrent-safe updates.
> I created an interface called DB in the db package that abstracts the implementation of a DB, so this can be extended to not only in-memory DB's but also SQL and NOSQL DB's.
> Based on the DB interface, I had created the RedisClientDB which can be switched out for another DB.

5. I added an optional Prometheus Counter as I wasn't too sure on the streaming service requirement, but I decided to keep it as an interface for counting/keeping track of any other metrics.

6. Since we wanted to push the log (either to file/MQ) every minute, I decided to keep a background job goroutine that runs a ticker with 1 minute interval. This goroutine does 2 things every minute - log unique requests and reset the count for the next upcoming minute.

7. I've tried to keep it as extensible as possible by utilizing interfaces wherever possible and avoid fixing hard dependenices on a particular technology.
> As seen in the README.md, the RPS is a bit dependent on the CPU cores available. I tested this using two machines, one with 4 cores and one with 14 cores. I believe on a decent production machine, this API will easily scale past 10K RPS.