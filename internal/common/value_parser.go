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

func (config *FyodorovConfig) parseToolKey(key []string, value string) {
	if len(key) == 0 {
		return
	}
	if config.Tools == nil {
		config.Tools = &[]ToolConfig{}
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
		config.Agents = &[]AgentConfig{}
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
