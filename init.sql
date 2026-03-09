CREATE TABLE IF NOT EXISTS battles (
    id SERIAL PRIMARY KEY,
    boss_id INTEGER NOT NULL,
    result VARCHAR(20) NOT NULL,
    duration INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS rankings (
    id SERIAL PRIMARY KEY,
    nickname VARCHAR(100) NOT NULL,
    class VARCHAR(50),
    total_damage INTEGER,
    incidents_solved INTEGER,
    battle_id INTEGER NOT NULL,
    CONSTRAINT fk_battle
        FOREIGN KEY (battle_id)
        REFERENCES battles(id)
        ON DELETE CASCADE
);