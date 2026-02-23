CREATE TABLE click_logs (
    id         BIGSERIAL PRIMARY KEY,
    url_id     BIGINT NOT NULL REFERENCES urls(id),
    ip_address VARCHAR(45),
    user_agent TEXT,
    referer    TEXT,
    clicked_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_click_logs_url_id ON click_logs(url_id);
CREATE INDEX idx_click_logs_clicked_at ON click_logs(clicked_at);
