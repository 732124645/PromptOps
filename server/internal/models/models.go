package models

import "time"

// Prompt is the core resource: a versioned, environment-scoped prompt.
type Prompt struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Key       string    `gorm:"index" json:"key"`
	Name      string    `json:"name"`
	Content   string    `json:"content"`
	Version   string    `json:"version"`
	Env       string    `gorm:"index" json:"env"`
	Category  string    `gorm:"index" json:"category"`
	Tags      string    `json:"tags"`
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PromptVersion is an immutable snapshot of a prompt's content, created on publish.
type PromptVersion struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	PromptID  string    `gorm:"index" json:"prompt_id"`
	Key       string    `gorm:"index" json:"key"`
	Version   string    `json:"version"`
	Env       string    `json:"env"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// Agent is a reusable, config-driven prompt + model setup.
type Agent struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	Key         string    `gorm:"index" json:"key"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Prompt      string    `json:"prompt"`
	Provider    string    `json:"provider"`
	Model       string    `json:"model"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Workflow is an ordered pipeline of steps. Steps holds a JSON-encoded array
// of workflow steps; handlers expose it as a structured "steps" field.
type Workflow struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	Key         string    `gorm:"index" json:"key"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Steps       string    `json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
