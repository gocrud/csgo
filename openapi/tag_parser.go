package openapi

import (
	"reflect"
	"strconv"
	"strings"
)

// FieldInfo contains OpenAPI information for a struct field.
type FieldInfo struct {
	Name        string
	Description string
	Example     interface{}
	Required    bool
	Format      string
	Minimum     *float64
	Maximum     *float64
	MinLength   *int
	MaxLength   *int
	Pattern     string
	Enum        []interface{}
}

// ParseStructTags parses struct tags to generate field information.
// It supports both the custom "doc" tag and individual tags for compatibility.
func ParseStructTags(t reflect.Type) map[string]FieldInfo {
	fields := make(map[string]FieldInfo)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		info := parseFieldTags(field)
		if info.Name != "" && info.Name != "-" {
			fields[info.Name] = info
		}
	}

	return fields
}

func parseFieldTags(field reflect.StructField) FieldInfo {
	info := FieldInfo{
		Name: field.Name,
	}

	// Parse json tag for field name
	if jsonTag := field.Tag.Get("json"); jsonTag != "" {
		parts := strings.Split(jsonTag, ",")
		if parts[0] == "-" {
			return FieldInfo{} // Skip this field
		}
		if parts[0] != "" {
			info.Name = parts[0]
		}
		// If json tag doesn't have omitempty, field is required by default
		if !contains(parts, "omitempty") {
			info.Required = true
		}
	}

	// Parse custom "doc" tag - format: "描述,required,example:value,min:0,max:100"
	if docTag := field.Tag.Get("doc"); docTag != "" {
		parseDocTag(&info, docTag)
	}

	// Also support individual tags for backward compatibility or explicit overrides
	if desc := field.Tag.Get("description"); desc != "" {
		info.Description = desc
	}
	if example := field.Tag.Get("example"); example != "" {
		info.Example = example
	}
	if format := field.Tag.Get("format"); format != "" {
		info.Format = format
	}
	if min := field.Tag.Get("minimum"); min != "" {
		if val, err := strconv.ParseFloat(min, 64); err == nil {
			info.Minimum = &val
		}
	}
	if max := field.Tag.Get("maximum"); max != "" {
		if val, err := strconv.ParseFloat(max, 64); err == nil {
			info.Maximum = &val
		}
	}
	if minLen := field.Tag.Get("minLength"); minLen != "" {
		if val, err := strconv.Atoi(minLen); err == nil {
			info.MinLength = &val
		}
	}
	if maxLen := field.Tag.Get("maxLength"); maxLen != "" {
		if val, err := strconv.Atoi(maxLen); err == nil {
			info.MaxLength = &val
		}
	}
	if pattern := field.Tag.Get("pattern"); pattern != "" {
		info.Pattern = pattern
	}
	if enumTag := field.Tag.Get("enum"); enumTag != "" {
		parseEnumTag(&info, enumTag)
	}

	return info
}

// parseDocTag parses the custom doc tag.
// Format: "描述,required,example:value,minLength:2,maxLength:50,min:0,max:100,format:email,enum:a|b|c"
func parseDocTag(info *FieldInfo, docTag string) {
	parts := strings.Split(docTag, ",")
	if len(parts) == 0 {
		return
	}

	// First part is always the description
	info.Description = strings.TrimSpace(parts[0])

	// Parse remaining parts as key:value or flags
	for i := 1; i < len(parts); i++ {
		part := strings.TrimSpace(parts[i])
		if part == "" {
			continue
		}

		// Check if it's a key:value pair
		if strings.Contains(part, ":") {
			kv := strings.SplitN(part, ":", 2)
			if len(kv) != 2 {
				continue
			}
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])

			switch key {
			case "example":
				info.Example = value
			case "format":
				info.Format = value
			case "min":
				// For numbers, this is minimum
				if val, err := strconv.ParseFloat(value, 64); err == nil {
					info.Minimum = &val
				}
			case "max":
				// For numbers, this is maximum
				if val, err := strconv.ParseFloat(value, 64); err == nil {
					info.Maximum = &val
				}
			case "minLength":
				if val, err := strconv.Atoi(value); err == nil {
					info.MinLength = &val
				}
			case "maxLength":
				if val, err := strconv.Atoi(value); err == nil {
					info.MaxLength = &val
				}
			case "minimum":
				if val, err := strconv.ParseFloat(value, 64); err == nil {
					info.Minimum = &val
				}
			case "maximum":
				if val, err := strconv.ParseFloat(value, 64); err == nil {
					info.Maximum = &val
				}
			case "pattern":
				info.Pattern = value
			case "enum":
				// Enum values separated by |
				enumValues := strings.Split(value, "|")
				info.Enum = make([]interface{}, len(enumValues))
				for j, v := range enumValues {
					info.Enum[j] = strings.TrimSpace(v)
				}
			}
		} else {
			// It's a flag
			switch part {
			case "required":
				info.Required = true
			}
		}
	}
}

// parseEnumTag parses enum tag (comma or pipe separated).
func parseEnumTag(info *FieldInfo, enumTag string) {
	// Support both comma and pipe as separators
	var enumValues []string
	if strings.Contains(enumTag, "|") {
		enumValues = strings.Split(enumTag, "|")
	} else {
		enumValues = strings.Split(enumTag, ",")
	}

	info.Enum = make([]interface{}, len(enumValues))
	for i, v := range enumValues {
		info.Enum[i] = strings.TrimSpace(v)
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

