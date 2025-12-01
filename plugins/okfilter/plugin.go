// Package okfilter provides a Bifrost plugin that filters out "OK." responses
// from models that incorrectly prepend this text to their responses.
package okfilter

import (
	"context"
	"strings"

	"github.com/maximhq/bifrost/core/schemas"
)

const (
	PluginName = "okfilter"
)

// Config represents the configuration for the OK filter plugin
type Config struct {
	// Models is a list of model names that should have the "OK." prefix filtered
	Models []string `json:"models,omitempty"`
}

// OKFilterPlugin implements the schemas.Plugin interface
type OKFilterPlugin struct {
	models map[string]bool
}

// Init creates a new OK filter plugin with the given configuration
func Init(_ context.Context, config *Config, _ schemas.Logger) (*OKFilterPlugin, error) {
	plugin := &OKFilterPlugin{
		models: make(map[string]bool),
	}

	// If no models are specified, the plugin will filter all models
	if config != nil && len(config.Models) > 0 {
		for _, model := range config.Models {
			plugin.models[model] = true
		}
	}

	return plugin, nil
}

// GetName returns the name of the plugin
func (p *OKFilterPlugin) GetName() string {
	return PluginName
}

// TransportInterceptor is not used for this plugin
func (p *OKFilterPlugin) TransportInterceptor(ctx *context.Context, url string, headers map[string]string, body map[string]any) (map[string]string, map[string]any, error) {
	return headers, body, nil
}

// PreHook is called before a request is processed - no modifications needed
func (p *OKFilterPlugin) PreHook(ctx *context.Context, req *schemas.BifrostRequest) (*schemas.BifrostRequest, *schemas.PluginShortCircuit, error) {
	return req, nil, nil
}

// PostHook is called after a response is received and filters out "OK." prefixes
func (p *OKFilterPlugin) PostHook(ctx *context.Context, result *schemas.BifrostResponse, bifrostErr *schemas.BifrostError) (*schemas.BifrostResponse, *schemas.BifrostError, error) {
	// If there's an error, we don't need to filter anything
	if bifrostErr != nil || result == nil {
		return result, bifrostErr, nil
	}

	// Check if we should filter this model
	shouldFilter := len(p.models) == 0 // If no models specified, filter all
	if !shouldFilter && result != nil {
		// Get the model from the response
		extraFields := result.GetExtraFields()
		if extraFields.ModelRequested != "" {
			_, shouldFilter = p.models[extraFields.ModelRequested]
		}
	}

	// If we shouldn't filter this model, return as-is
	if !shouldFilter {
		return result, bifrostErr, nil
	}

	// Filter the response based on its type
	switch {
	case result.ChatResponse != nil:
		p.filterChatResponse(result.ChatResponse)
	case result.TextCompletionResponse != nil:
		p.filterTextCompletionResponse(result.TextCompletionResponse)
	case result.ResponsesResponse != nil:
		p.filterResponsesResponse(result.ResponsesResponse)
	}

	return result, bifrostErr, nil
}

// filterChatResponse filters "OK." prefixes from chat responses
func (p *OKFilterPlugin) filterChatResponse(response *schemas.BifrostChatResponse) {
	if response == nil || len(response.Choices) == 0 {
		return
	}

	for i := range response.Choices {
		choice := &response.Choices[i]
		
		// Handle non-streaming responses
		if choice.ChatNonStreamResponseChoice != nil && choice.ChatNonStreamResponseChoice.Message != nil {
			message := choice.ChatNonStreamResponseChoice.Message
			
			// Check if content exists and starts with "OK."
			if message.Content != nil && message.Content.ContentStr != nil {
				content := *message.Content.ContentStr
				
				// Check if content starts with "OK." (case-insensitive)
				if strings.HasPrefix(strings.ToLower(strings.TrimSpace(content)), "ok.") {
					// Remove the "OK." prefix and any leading whitespace
					newContent := strings.TrimLeft(strings.TrimPrefix(content, content[:3]), " \t\n\r")
					
					// If the content is empty after removing "OK.", set it to nil
					if newContent == "" {
						message.Content.ContentStr = nil
					} else {
						message.Content.ContentStr = &newContent
					}
				}
			}
		}
		
		// Handle streaming responses (if needed)
		if choice.ChatStreamResponseChoice != nil && choice.ChatStreamResponseChoice.Delta != nil {
			delta := choice.ChatStreamResponseChoice.Delta
			
			// For streaming, we only filter if it's the first chunk and starts with "OK."
			if delta.Content != nil && strings.HasPrefix(strings.ToLower(strings.TrimSpace(*delta.Content)), "ok.") {
				// Remove the "OK." prefix and any leading whitespace
				newContent := strings.TrimLeft(strings.TrimPrefix(*delta.Content, (*delta.Content)[:3]), " \t\n\r")
				delta.Content = &newContent
			}
		}
	}
}

// filterTextCompletionResponse filters "OK." prefixes from text completion responses
func (p *OKFilterPlugin) filterTextCompletionResponse(response *schemas.BifrostTextCompletionResponse) {
	if response == nil || len(response.Choices) == 0 {
		return
	}

	for i := range response.Choices {
		choice := &response.Choices[i]
		
		// Handle text completion responses
		if choice.TextCompletionResponseChoice != nil && choice.TextCompletionResponseChoice.Text != nil {
			content := *choice.TextCompletionResponseChoice.Text
			
			// Check if content starts with "OK." (case-insensitive)
			if strings.HasPrefix(strings.ToLower(strings.TrimSpace(content)), "ok.") {
				// Remove the "OK." prefix and any leading whitespace
				newContent := strings.TrimLeft(strings.TrimPrefix(content, content[:3]), " \t\n\r")
				
				// If the content is empty after removing "OK.", set it to nil
				if newContent == "" {
					choice.TextCompletionResponseChoice.Text = nil
				} else {
					choice.TextCompletionResponseChoice.Text = &newContent
				}
			}
		}
	}
}

// filterResponsesResponse filters "OK." prefixes from responses responses
func (p *OKFilterPlugin) filterResponsesResponse(response *schemas.BifrostResponsesResponse) {
	if response == nil || len(response.Output) == 0 {
		return
	}

	for i := range response.Output {
		output := &response.Output[i]
		
		// Check if content exists and starts with "OK."
		if output.Content != nil && output.Content.ContentStr != nil {
			content := *output.Content.ContentStr
			
			// Check if content starts with "OK." (case-insensitive)
			if strings.HasPrefix(strings.ToLower(strings.TrimSpace(content)), "ok.") {
				// Remove the "OK." prefix and any leading whitespace
				newContent := strings.TrimLeft(strings.TrimPrefix(content, content[:3]), " \t\n\r")
				
				// If the content is empty after removing "OK.", set it to nil
				if newContent == "" {
					output.Content.ContentStr = nil
				} else {
					output.Content.ContentStr = &newContent
				}
			}
		}
	}
}

// Cleanup is called when the plugin is being shut down
func (p *OKFilterPlugin) Cleanup() error {
	// No cleanup needed for this plugin
	return nil
}