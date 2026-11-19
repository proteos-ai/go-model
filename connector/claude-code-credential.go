package connectormodel

import "strings"

// ConnectorKeyClaudeCode is the pre-built connector that stores the
// credential Claude Code runs on in a person's agent sandboxes.
const ConnectorKeyClaudeCode = "claude-code"

// The environment variables Claude Code reads its credential from.
const (
	// ClaudeCodeEnvOAuthToken carries a subscription token from `claude setup-token`.
	ClaudeCodeEnvOAuthToken = "CLAUDE_CODE_OAUTH_TOKEN"
	// ClaudeCodeEnvApiKey carries an Anthropic API key.
	ClaudeCodeEnvApiKey = "ANTHROPIC_API_KEY"
)

const (
	claudeCodeOAuthTokenPrefix = "sk-ant-oat"
	claudeCodeApiKeyPrefix     = "sk-ant-api"
)

// ClaudeCodeEnvironmentVariable names the variable Claude Code reads a stored
// credential from: a subscription token and an API key are not
// interchangeable, and Claude Code rejects one passed as the other.
// isRecognized is false for a value that is neither.
func ClaudeCodeEnvironmentVariable(credential string) (name string, isRecognized bool) {
	switch {
	case strings.HasPrefix(credential, claudeCodeOAuthTokenPrefix):
		return ClaudeCodeEnvOAuthToken, true
	case strings.HasPrefix(credential, claudeCodeApiKeyPrefix):
		return ClaudeCodeEnvApiKey, true
	default:
		return "", false
	}
}
