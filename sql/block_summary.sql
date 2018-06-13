CREATE TABLE IF NOT EXISTS zerochain.block_summary (
hash text,
prev_hash text,
round bigint,
creation_date bigint,
version text,
round_random_seed bigint,
PRIMARY KEY(hash, round)
);

CREATE INDEX IF NOT EXISTS block_summary_nu1_creation_date ON zerochain.block_summary (creation_date);
CREATE INDEX IF NOT EXISTS block_summary_u1_round ON zerochain.block_summary (round);
