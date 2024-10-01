CREATE TABLE IF NOT EXISTS championship_teams (
    championship_id UUID NOT NULL REFERENCES championships(id) ON DELETE CASCADE,
    team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    PRIMARY KEY (championship_id, team_id)
);

CREATE INDEX idx_championship_teams_championship_id ON championship_teams(championship_id);
CREATE INDEX idx_championship_teams_team_id ON championship_teams(team_id);