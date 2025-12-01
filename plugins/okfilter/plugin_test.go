package okfilter

import (
	"context"
	"testing"

	"github.com/maximhq/bifrost/core/schemas"
	"github.com/stretchr/testify/assert"
)

func TestOKFilterPlugin_Init(t *testing.T) {
	// Test with no config
	plugin, err := Init(context.Background(), nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, plugin)
	assert.Equal(t, PluginName, plugin.GetName())

	// Test with config
	config := &Config{
		Models: []string{"kimi-k2-thinking"},
	}
	plugin, err = Init(context.Background(), config, nil)
	assert.NoError(t, err)
	assert.NotNil(t, plugin)
}

func TestOKFilterPlugin_FilterChatResponse(t *testing.T) {
	plugin := &OKFilterPlugin{
		models: map[string]bool{"test-model": true},
	}

	// Test with "OK." prefix
	content := "OK. This is a test response"
	response := &schemas.BifrostChatResponse{
		Choices: []schemas.BifrostResponseChoice{
			{
				ChatNonStreamResponseChoice: &schemas.ChatNonStreamResponseChoice{
					Message: &schemas.ChatMessage{
						Content: &schemas.ChatMessageContent{
							ContentStr: &content,
						},
					},
				},
			},
		},
	}

	plugin.filterChatResponse(response)
	
	// Check that the "OK." prefix was removed
	assert.NotNil(t, response.Choices[0].ChatNonStreamResponseChoice.Message.Content.ContentStr)
	assert.Equal(t, "This is a test response", *response.Choices[0].ChatNonStreamResponseChoice.Message.Content.ContentStr)

	// Test with only "OK."
	content2 := "OK."
	response2 := &schemas.BifrostChatResponse{
		Choices: []schemas.BifrostResponseChoice{
			{
				ChatNonStreamResponseChoice: &schemas.ChatNonStreamResponseChoice{
					Message: &schemas.ChatMessage{
						Content: &schemas.ChatMessageContent{
							ContentStr: &content2,
						},
					},
				},
			},
		},
	}

	plugin.filterChatResponse(response2)
	
	// Check that the content is now nil
	assert.Nil(t, response2.Choices[0].ChatNonStreamResponseChoice.Message.Content.ContentStr)
}

func TestOKFilterPlugin_PostHook(t *testing.T) {
	// Test with specific models configured
	plugin := &OKFilterPlugin{
		models: map[string]bool{"kimi-k2-thinking": true},
	}

	content := "OK. Hello world"
	response := &schemas.BifrostResponse{
		ChatResponse: &schemas.BifrostChatResponse{
			Choices: []schemas.BifrostResponseChoice{
				{
					ChatNonStreamResponseChoice: &schemas.ChatNonStreamResponseChoice{
						Message: &schemas.ChatMessage{
							Content: &schemas.ChatMessageContent{
								ContentStr: &content,
							},
						},
					},
				},
			},
			ExtraFields: schemas.BifrostResponseExtraFields{
				ModelRequested: "kimi-k2-thinking",
			},
		},
	}

	result, err, hookErr := plugin.PostHook(nil, response, nil)
	assert.NoError(t, hookErr)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	// Check that the "OK." prefix was removed
	assert.Equal(t, "Hello world", *result.ChatResponse.Choices[0].ChatNonStreamResponseChoice.Message.Content.ContentStr)

	// Test with a model that shouldn't be filtered
	response2 := &schemas.BifrostResponse{
		ChatResponse: &schemas.BifrostChatResponse{
			Choices: []schemas.BifrostResponseChoice{
				{
					ChatNonStreamResponseChoice: &schemas.ChatNonStreamResponseChoice{
						Message: &schemas.ChatMessage{
							Content: &schemas.ChatMessageContent{
								ContentStr: &content,
							},
						},
					},
				},
			},
			ExtraFields: schemas.BifrostResponseExtraFields{
				ModelRequested: "other-model",
			},
		},
	}

	result2, err2, hookErr2 := plugin.PostHook(nil, response2, nil)
	assert.NoError(t, hookErr2)
	assert.Nil(t, err2)
	assert.NotNil(t, result2)

	// Check that the content was not modified
	assert.Equal(t, "OK. Hello world", *result2.ChatResponse.Choices[0].ChatNonStreamResponseChoice.Message.Content.ContentStr)
}