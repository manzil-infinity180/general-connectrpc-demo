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
---

```md
// Create Cricket Poll

curl -X POST http://localhost:8080/voting.v1.VotingService/CreatePoll \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer fake-dev-token" \
  -d '{
    "question": "Who is the greatest batsman of all time?",
    "option_texts": [
      "Sachin Tendulkar (India)",
      "Virat Kohli (India)",
      "Rohit Sharma (India)",
      "Ricky Ponting (Australia)",
      "Chris Gayle (West Indies)",
      "Pat Cummins (Australia)",
      "AB de Villiers (South Africa)"
    ]
  }'

```
```md
SELECT id, text FROM options WHERE poll_id='efbc5848-fe0f-49a8-94eb-8c1847e0fffd';

voting=> SELECT id, text FROM options WHERE poll_id='efbc5848-fe0f-49a8-94eb-8c1847e0fffd';
id                  |             text              
--------------------------------------+-------------------------------
df7d77b0-32a6-426b-8125-989b03cf1ec2 | Sachin Tendulkar (India)
a8795838-151b-4024-a086-21420b61c4da | Virat Kohli (India)
455584b7-ad1a-49b7-9f97-23a06371e3c1 | Rohit Sharma (India)
63631c5d-cd29-4711-a3d9-ed531c42a79c | Ricky Ponting (Australia)
de0949f8-38d7-4c6e-a02e-1dc6b6ecb7b4 | Chris Gayle (West Indies)
b2bdba1b-c138-433c-894a-2dedbbb7156c | Pat Cummins (Australia)
27b14ca2-17e0-40ae-8048-9cdeceea9ae7 | AB de Villiers (South Africa)
(7 rows)
```

```md
SELECT o.text, o.vote_count FROM options o WHERE o.poll_id='efbc5848-fe0f-49a8-94eb-8c1847e0fffd' ORDER BY o.vote_count DESC;
             text              | vote_count 
-------------------------------+------------
 Sachin Tendulkar (India)      |          3
 Virat Kohli (India)           |          2
 Chris Gayle (West Indies)     |          0
 Rohit Sharma (India)          |          0
 AB de Villiers (South Africa) |          0
 Pat Cummins (Australia)       |          0
 Ricky Ponting (Australia)     |          0
(7 rows)

```

```md
// close the voting 

curl -X POST http://localhost:8080/voting.v1.VotingService/ClosePoll \
-H "Content-Type: application/json" \
-H "Authorization: Bearer fake-dev-token" \
-d "{\"poll_id\": \"$POLL_ID\"}"


voting=> SELECT
ROW_NUMBER() OVER (ORDER BY o.vote_count DESC) as rank,
o.text as player,
o.vote_count as votes,
ROUND(o.vote_count * 100.0 / NULLIF(SUM(o.vote_count) OVER (), 0), 1) as percentage
FROM options o
WHERE o.poll_id='efbc5848-fe0f-49a8-94eb-8c1847e0fffd' ORDER BY o.vote_count DESC;
rank |            player             | votes | percentage
------+-------------------------------+-------+------------
1 | Virat Kohli (India)           |     6 |       54.5
2 | Sachin Tendulkar (India)      |     5 |       45.5
3 | Chris Gayle (West Indies)     |     0 |        0.0
4 | Rohit Sharma (India)          |     0 |        0.0
5 | AB de Villiers (South Africa) |     0 |        0.0
6 | Pat Cummins (Australia)       |     0 |        0.0
7 | Ricky Ponting (Australia)     |     0 |        0.0
(7 rows)

```