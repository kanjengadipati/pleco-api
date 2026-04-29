CREATE TABLE error_analyses (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    error_message TEXT NOT NULL,
    error_type VARCHAR(255) NOT NULL,
    root_cause TEXT,
    affected_components TEXT,
    recommended_action TEXT,
    severity VARCHAR(50)
);
