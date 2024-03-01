package common

import (
	"fmt"
	"strconv"
	"strings"
)

func (config *FyodorovConfig) ParseKeyValuePairs(args []string) {
	pairs := make(map[string]string)
	for _, arg := range args {
		parts := strings.Split(arg, "=")
		if len(parts) == 2 {
			pairs[parts[0]] = parts[1]
		}
	}
	config.parseValues(pairs)
}

func (config *FyodorovConfig) parseValues(args map[string]string) {
	for key, value := range args {
		config.parseKey(key, value)
	}
}

func (config *FyodorovConfig) parseKey(key, value string) {
	keys := strings.Split(key, ".")
	if len(keys) >= 1 {
		firstKey, remainingKeys := parseComplexKey(keys)
		switch firstKey {
		case "version":
			config.Version = value
			fmt.Print("Using fyodorov version: ", config.Version, "\n")
		case "providers":
			config.parseProviderKey(remainingKeys, value)
		case "models":
			config.parseModelKey(remainingKeys, value)
		case "agents":
			config.parseAgentKey(remainingKeys, value)
		case "tools":
			config.parseToolKey(remainingKeys, value)
		default:
			fmt.Printf("Unknown key: %s\n", firstKey)
		}
	}
}

func parseComplexKey(keys []string) (firstKey string, remainingKeys []string) {
	if len(keys) == 0 {
		return "", nil
	}
	firstKey = keys[0]
	remainingKeys = keys[1:]
	if newFirstKey, nextKey, found := strings.Cut(firstKey, "["); found {
		firstKey = newFirstKey
		remainingKeys = append([]string{"[" + nextKey}, remainingKeys...)
	}
	return firstKey, remainingKeys
}

func (config *FyodorovConfig) parseProviderKey(key []string, value string) {
	if len(key) == 0 {
		return
	}
	if config.Providers == nil {
		config.Providers = &[]Provider{}
	}
	index := parseIndex(key[0])
	if index >= len(*config.Providers) {
		fmt.Printf("Invalid index: %d\n", index)
		return
	}
	if len(key) == 1 {
		return
	}
	switch key[1] {
	case "name":
		(*config.Providers)[index].Name = value
	case "url":
		(*config.Providers)[index].URL = value
	case "api_key":
		(*config.Providers)[index].APIKey = value
	default:
		fmt.Printf("Unknown key: %s\n", key[1])
	}
}

func (config *FyodorovConfig) parseModelKey(key []string, value string) {
	if len(key) == 0 {
		return
	}
	if config.Models == nil {
		config.Models = &[]Model{}
	}
	index := parseIndex(key[0])
	if index >= len(*config.Models) {
		fmt.Printf("Invalid index: %d\n", index)
		return
	}
	if len(key) == 1 {
		return
	}
	switch key[1] {
	case "name":
		(*config.Models)[index].Name = value
	case "provider":
		(*config.Models)[index].Provider = value
	case "params":
		(*config.Models)[index].Params = value
	case "model_info":
		(*config.Models)[index].parseModelInfo(key[2:], value)
	default:
		fmt.Printf("Unknown key: %s\n", key[1])
	}
}

func (config *ModelConfig) parseModelInfo(key []string, value string) {
	if len(key) == 0 {
		return
	}
	switch key[1] {
	case "mode":
		config.ModelInfo.Mode = value
	case "input_cost_per_token":
		inputCostPerToken, err := strconv.ParseFloat(value, 64)
		if err != nil {
			fmt.Printf("Invalid input cost per token: %s\n", value)
			return
		}
		config.ModelInfo.InputCostPerToken = &inputCostPerToken
	case "output_cost_per_token":
		outputCostPerToken, err := strconv.ParseFloat(value, 64)
		if err != nil {
			fmt.Printf("Invalid output cost per token: %s\n", value)
			return
		}
		config.ModelInfo.OutputCostPerToken = &outputCostPerToken
	case "max_tokens":
		maxTokens, err := strconv.Atoi(value)
		if err != nil {
			fmt.Printf("Invalid max tokens: %s\n", value)
			return
		}
		config.ModelInfo.MaxTokens = &maxTokens
	case "base_model":
		config.ModelInfo.BaseModel = value
	default:
		fmt.Printf("Unknown key: %s\n", key[1])
	}
}

func (config *FyodorovConfig) parseToolKey(key []string, value string) {
	if len(key) == 0 {
		return
	}
	if config.Tools == nil {
		config.Tools = &[]Tool{}
	}
	index := parseIndex(key[0])
	if index >= len(*config.Tools) {
		fmt.Printf("Invalid index: %d\n", index)
		return
	}
	if len(key) == 1 {
		return
	}
	switch key[1] {
	case "name":
		(*config.Tools)[index].NameForHuman = value
	case "description":
		(*config.Tools)[index].DescriptionForHuman = value
	case "name_for_ai":
		(*config.Tools)[index].NameForAI = value
	case "description_for_ai":
		(*config.Tools)[index].DescriptionForAI = value
	case "api":
		switch key[2] {
		case "type":
			(*config.Tools)[index].API.Type = value
		case "url":
			(*config.Tools)[index].API.URL = value
		default:
			fmt.Printf("Unknown key under api: %s\n", key[2])
		}
	case "logo_url":
		(*config.Tools)[index].LogoURL = value
	case "contact_email":
		(*config.Tools)[index].ContactEmail = value
	case "legal_info_url":
		(*config.Tools)[index].LegalInfoURL = value
	case "auth":
		switch key[2] {
		case "type":
			(*config.Tools)[index].Auth.Type = value
		default:
			fmt.Printf("Unknown key under auth: %s\n", key[2])
		}
	default:
		fmt.Printf("Unknown key: %s\n", key[1])
	}
}

func (config *FyodorovConfig) parseAgentKey(key []string, value string) {
	if len(key) == 0 {
		return
	}
	if config.Agents == nil {
		config.Agents = &[]Agent{}
	}
	index := parseIndex(key[0])
	if index >= len(*config.Agents) {
		fmt.Printf("Agent index out of range: %d\n", index)
		return
	}
	if len(key) == 1 {
		return
	}
	switch key[1] {
	case "model":
		(*config.Agents)[index].Model = value
	case "name":
		(*config.Agents)[index].NameForHuman = value
	case "description":
		(*config.Agents)[index].DescriptionForHuman = value
	case "prompt":
		(*config.Agents)[index].Prompt = value
	case "tools":
		(*config.Agents)[index].Tools = strings.Split(value, ",")
	case "rag":
		(*config.Agents)[index].Rag = strings.Split(value, ",")
	default:
		fmt.Printf("Unknown key: %s\n", key[1])
	}
}

func parseIndex(key string) int {
	if !strings.HasPrefix(key, "[") || !strings.HasSuffix(key, "]") {
		fmt.Printf("Invalid index: %s\n", key)
		return -1
	}
	key = key[1 : len(key)-1]
	index, err := strconv.Atoi(key)
	if err != nil {
		fmt.Printf("Invalid index: %s\n", key)
		return -1
	}
	return index
}
