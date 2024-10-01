CREATE TYPE championship_type AS ENUM ('league', 'cup');
CREATE TYPE tiebreaker_method AS ENUM ('penalties', 'extra_time');
CREATE TYPE progression_type AS ENUM ('fixed', 'random_draw');

CREATE TABLE IF NOT EXISTS championships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    type championship_type NOT NULL,
    tiebreaker_method tiebreaker_method NOT NULL DEFAULT 'penalties',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    phases INTEGER NOT NULL,
    progression_type progression_type NOT NULL DEFAULT 'fixed'
);

CREATE INDEX idx_championships_type ON championships(type);
