package handlers

import (
	"testing"
)

func TestParseQuestionParams(t *testing.T) {
	// 解析问卷参数
	params, err := ParseQuestionParams("eyJQYWdlVHlwZSI6InF1ZXN0aW9uIiwiUXVlc3Rpb25uYWlyZUlEIjoxLCJRdWVzdGlvbkluZGV4IjowfQ==")
	if err != nil {
		t.Errorf("Failed to parse question params: %s", err)
	}
	t.Log(params)
}
