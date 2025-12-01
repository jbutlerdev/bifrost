# OK Filter Plugin

This plugin filters out unwanted "OK." prefixes from model responses. Some models (like kimi-k2-thinking) have an issue where they prepend "OK." to their responses, which can be undesirable for users.

## Features

- Automatically removes "OK." prefixes from model responses
- Works with chat completions, text completions, and responses
- Configurable to target specific models or all models
- Handles both streaming and non-streaming responses
- Case-insensitive filtering

## Configuration

The plugin can be configured in your Bifrost configuration file:

```json
{
  "plugins": [
    {
      "name": "okfilter",
      "enabled": true,
      "config": {
        "models": ["kimi-k2-thinking"]  // Optional: only filter specific models
      }
    }
  ]
}
```

If no models are specified in the configuration, the plugin will filter "OK." prefixes from all models.

## How It Works

The plugin implements a PostHook that examines responses after they are received from the model. It checks if the response content starts with "OK." (case-insensitive) and removes this prefix along with any leading whitespace.

For streaming responses, it only filters the first chunk if it starts with "OK."

## Installation

The plugin is included with Bifrost by default. Simply enable it in your configuration file.