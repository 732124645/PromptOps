package models

import "time"

// Workspace is a team space that groups prompts, agents and workflows.
type Workspace struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

// Prompt is the core resource: a versioned, environment-scoped prompt.
type Prompt struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	WorkspaceID string    `gorm:"index" json:"workspace_id"`
	Key         string    `gorm:"index" json:"key"`
	Name        string    `json:"name"`
	Content     string    `json:"content"`
	Version     string    `json:"version"`
	Env         string    `gorm:"index" json:"env"`
	Category    string    `gorm:"index" json:"category"`
	Tags        string    `json:"tags"`
	Model       string    `json:"model"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
	WorkspaceID string    `gorm:"index" json:"workspace_id"`
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
	WorkspaceID string    `gorm:"index" json:"workspace_id"`
	Key         string    `gorm:"index" json:"key"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Steps       string    `json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Rollout splits SDK traffic for a prompt key+env between two versions.
type Rollout struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Key       string    `gorm:"index" json:"key"`
	Env       string    `gorm:"index" json:"env"`
	Enabled   bool      `json:"enabled"`
	VariantA  string    `json:"variant_a"`
	VariantB  string    `json:"variant_b"`
	WeightA   int       `json:"weight_a"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AuditLog records a mutation to a prompt, agent or workflow.
type AuditLog struct {
	ID         string    `gorm:"primaryKey" json:"id"`
	Action     string    `json:"action"`
	Resource   string    `json:"resource"`
	ResourceID string    `json:"resource_id"`
	Key        string    `json:"key"`
	Summary    string    `json:"summary"`
	CreatedAt  time.Time `gorm:"index" json:"created_at"`
}

// User is an account with a role: "admin", "editor" or "viewer".
type User struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"uniqueIndex" json:"username"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

// Session maps an opaque bearer token to a user.
type Session struct {
	Token     string    `gorm:"primaryKey" json:"-"`
	UserID    string    `gorm:"index" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

// RunLog records a single model invocation for observability.
type RunLog struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	Source       string    `json:"source"`
	RefKey       string    `json:"ref_key"`
	Provider     string    `json:"provider"`
	Model        string    `json:"model"`
	PromptTokens int       `json:"prompt_tokens"`
	OutputTokens int       `json:"output_tokens"`
	LatencyMs    int64     `json:"latency_ms"`
	Status       string    `json:"status"`
	Error        string    `json:"error"`
	CreatedAt    time.Time `gorm:"index" json:"created_at"`
}
