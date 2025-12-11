-- Таблица для хранения информации о пользователях
CREATE TABLE IF NOT EXISTS users (
                                     user_id BIGINT PRIMARY KEY,
                                     username VARCHAR(255),
                                     first_name VARCHAR(255),
                                     first_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица для логирования использования техник
CREATE TABLE IF NOT EXISTS technique_usage (
                                               id SERIAL PRIMARY KEY,
                                               user_id BIGINT NOT NULL,
                                               technique_id VARCHAR(100) NOT NULL,
                                               category VARCHAR(50) NOT NULL,
                                               started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                               completed BOOLEAN DEFAULT FALSE,
                                               completed_at TIMESTAMP,
                                               FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- Индексы для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_technique_usage_user_id ON technique_usage(user_id);
CREATE INDEX IF NOT EXISTS idx_technique_usage_technique_id ON technique_usage(technique_id);
CREATE INDEX IF NOT EXISTS idx_technique_usage_started_at ON technique_usage(started_at);
CREATE INDEX IF NOT EXISTS idx_technique_usage_category ON technique_usage(category);

-- Представление для статистики по пользователям
CREATE OR REPLACE VIEW user_statistics AS
SELECT
    u.user_id,
    u.username,
    COUNT(tu.id) as total_techniques_used,
    COUNT(CASE WHEN tu.completed = true THEN 1 END) as completed_techniques,
    MAX(tu.started_at) as last_activity,
    COUNT(CASE WHEN tu.category = 'breathing' THEN 1 END) as breathing_count,
    COUNT(CASE WHEN tu.category = 'grounding' THEN 1 END) as grounding_count
FROM users u
         LEFT JOIN technique_usage tu ON u.user_id = tu.user_id
GROUP BY u.user_id, u.username;
