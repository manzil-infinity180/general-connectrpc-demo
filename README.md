## Basic Voting Application

* start the server 
```
$ go run cmd/server/main.go
2025/12/18 19:22:03 Using DATABASE_URL: postgres://
2025/12/18 19:22:03 Voting service running on http://localhost:8080
```
* check the health of the server - http://localhost:8080/health

* run the client demo 
```
$ go run cmd/client/main.go
Voting System Demo Client
STEP 1: Creating Poll
Poll created successfully!
Poll ID: f025f547-9bbb-4d5c-831c-f154be8d1a74

STEP 2: Fetching Options from Database
Found 4 options:
   1. Go (ID: e82c50c3-c7fd-41a9-9634-9e09c1bacb43)
   2. Java (ID: 01658909-b446-4105-9735-6d7ad7951da0)
   3. Python (ID: f0dde967-653b-4cd8-a11e-d2477ff1fb0e)
   4. Rust (ID: a69be411-b0a8-4e09-ab32-87a75db626b4)
STEP 3: Starting Real-time Stream
Stream connected! Waiting for vote updates...

2025/12/18 19:23:13 Stream ended: deadline_exceeded: Post "http://localhost:8080/voting.v1.VotingService/StreamResults": context deadline exceeded (Client.Timeout exceeded while awaiting headers)
STEP 4: Submitting Votes
Vote #1: Voting for 'Go'...
Vote recorded! (Success: true)
Vote #2: Voting for 'Java'...
Vote recorded! (Success: true)
Vote #3: Trying to vote again with same voter (should fail)...
Warning: Duplicate vote was allowed (should be prevented)
Vote #4: Voting for 'Python'...
Vote recorded! (Success: true)

STEP 5: Final Results
────────────────────────
Final Vote Counts:
   • Go: 2 vote(s)
   • Java: 1 vote(s)
   • Python: 1 vote(s)
   • Rust: 0 vote(s)

   Total: 4 vote(s) cast

STEP 6: Closing Poll
Poll closed successfully!
Warning: Vote on closed poll was allowed (should be prevented)

Poll ID: f025f547-9bbb-4d5c-831c-f154be8d1a74

You can view this poll in the database:
psql $DATABASE_URL -c "SELECT * FROM polls WHERE id='f025f547-9bbb-4d5c-831c-f154be8d1a74';"
```
