CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
 id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
 email TEXT UNIQUE NOT NULL,
 name TEXT NOT NULL,
 google_sub TEXT UNIQUE NOT NULL,
 created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE polls (
 id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
 question TEXT NOT NULL,
 creator_id UUID REFERENCES users(id),
 status TEXT NOT NULL DEFAULT 'OPEN',
 created_at TIMESTAMP DEFAULT now(),
 closed_at TIMESTAMP
);

CREATE TABLE options (
   id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
   poll_id UUID REFERENCES polls(id),
   text TEXT NOT NULL,
   vote_count INT DEFAULT 0
);

CREATE TABLE votes (
 id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
 poll_id UUID,
 option_id UUID,
 voter_hash TEXT,
 created_at TIMESTAMP DEFAULT now(),
 UNIQUE (poll_id, voter_hash)
);
