CREATE TYPE match_status AS ENUM ('scheduled', 'in_progress', 'finished');

CREATE TABLE IF NOT EXISTS matches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    championship_id UUID NOT NULL REFERENCES championships(id) ON DELETE CASCADE,
    home_team_id UUID NOT NULL REFERENCES teams(id),
    away_team_id UUID NOT NULL REFERENCES teams(id),
    match_date TIMESTAMP WITH TIME ZONE,
    status match_status NOT NULL DEFAULT 'scheduled',
    score_home INTEGER DEFAULT 0,
    score_away INTEGER DEFAULT 0,
    has_extra_time BOOLEAN DEFAULT FALSE,
    score_home_extra_time INTEGER DEFAULT 0,
    score_away_extra_time INTEGER DEFAULT 0,
    has_penalties BOOLEAN DEFAULT FALSE,
    score_home_penalties INTEGER DEFAULT 0,
    score_away_penalties INTEGER DEFAULT 0,
    winner_team_id UUID REFERENCES teams(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    phase INTEGER NOT NULL DEFAULT 1,
    parent_match_id UUID REFERENCES matches(id),
    left_child_match_id UUID REFERENCES matches(id),
    right_child_match_id UUID REFERENCES matches(id)
);

CREATE INDEX idx_matches_championship_id ON matches(championship_id);
CREATE INDEX idx_matches_home_team_id ON matches(home_team_id);
CREATE INDEX idx_matches_away_team_id ON matches(away_team_id);
CREATE INDEX idx_matches_status ON matches(status);