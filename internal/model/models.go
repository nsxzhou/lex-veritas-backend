package model

// AllModels 返回所有需要自动迁移的模型
func AllModels() []any {
	return []any{
		&User{},
		&OAuthAccount{},
		&TokenUsage{},
		&TokenUsageDaily{},
		&ChatSession{},
		&ChatMessage{},
		&MessageCitation{},
		&Document{},
		&DocumentChunk{},
		&KnowledgeVersion{},
		&ProofRecord{},
		&AuditLog{},
		&SystemConfig{},
	}
}
