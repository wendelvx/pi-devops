CREATE TABLE IF NOT EXISTS battles (
    id SERIAL PRIMARY KEY,
    boss_id INTEGER NOT NULL,
    result VARCHAR(20) NOT NULL,
    duration INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS rankings (
    id SERIAL PRIMARY KEY,
    nickname VARCHAR(100) NOT NULL UNIQUE,
    class VARCHAR(50),
    total_damage INTEGER DEFAULT 0,
    incidents_solved INTEGER DEFAULT 0,
    battle_id INTEGER NOT NULL,
    CONSTRAINT fk_battle
        FOREIGN KEY (battle_id)
        REFERENCES battles(id)
        ON DELETE CASCADE
);

-- PULO DO GATO: Inserir a batalha número 1 para a Chave Estrangeira do Go não falhar!
INSERT INTO battles (id, boss_id, result, duration) 
VALUES (1, 99, 'IN_PROGRESS', 0) 
ON CONFLICT (id) DO NOTHING;