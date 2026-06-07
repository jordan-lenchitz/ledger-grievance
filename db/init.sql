CREATE TABLE IF NOT EXISTS incidents (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  reporter_id VARCHAR(128) NOT NULL,
  occurred_at DATETIME NULL,
  recorded_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  subject VARCHAR(255) NOT NULL,
  category VARCHAR(128) NOT NULL DEFAULT 'unspecified',
  severity TINYINT UNSIGNED NOT NULL DEFAULT 1,
  description TEXT NOT NULL,
  evidence_uri TEXT NULL,
  requires_accommodation BOOLEAN NOT NULL DEFAULT FALSE,
  status ENUM('reported','reviewing','resolved','dismissed','archived') NOT NULL DEFAULT 'reported',
  notes TEXT NULL,
  CONSTRAINT chk_severity CHECK (severity BETWEEN 1 AND 5),
  INDEX idx_reporter_recorded (reporter_id, recorded_at),
  INDEX idx_subject (subject),
  INDEX idx_status (status),
  INDEX idx_category (category)
);
