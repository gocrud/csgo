package openapi

import (
	"reflect"
	"strconv"
	"strings"
)

// FieldInfo contains OpenAPI information for a struct field.
type FieldInfo struct {
	Name             string
	Description      string
	Example          interface{}
	Required         bool
	Format           string
	ContentMediaType string // For images: image/png, image/jpeg, etc.
	Minimum          *float64
	Maximum          *float64
	MinLength        *int
	MaxLength        *int
	Pattern          string
	Enum             []interface{}
}

// ParseStructTags parses struct tags to generate field information.
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

	// Parse individual tags
	if desc := field.Tag.Get("desc"); desc != "" {
		info.Description = desc
	}
	if example := field.Tag.Get("example"); example != "" {
		info.Example = example
	}

	// Parse file tag (convenience for file upload fields)
	// file:"true" -> format: binary (for multipart/form-data)
	if fileTag := field.Tag.Get("file"); fileTag == "true" {
		info.Format = "binary"
	}

	// Parse image tag (convenience for Base64 images)
	// image:"png" -> format: byte, media: image/png
	if imageTag := field.Tag.Get("image"); imageTag != "" {
		info.Format = "byte"
		info.ContentMediaType = parseImageMediaType(imageTag)
	}

	// Parse format and media tags (can override file and image tags)
	if format := field.Tag.Get("format"); format != "" {
		info.Format = format
	}
	if media := field.Tag.Get("media"); media != "" {
		info.ContentMediaType = media
	}

	if min := field.Tag.Get("min"); min != "" {
		if val, err := strconv.ParseFloat(min, 64); err == nil {
			info.Minimum = &val
		}
	}
	if max := field.Tag.Get("max"); max != "" {
		if val, err := strconv.ParseFloat(max, 64); err == nil {
			info.Maximum = &val
		}
	}
	if minLen := field.Tag.Get("minLen"); minLen != "" {
		if val, err := strconv.Atoi(minLen); err == nil {
			info.MinLength = &val
		}
	}
	if maxLen := field.Tag.Get("maxLen"); maxLen != "" {
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
	if required := field.Tag.Get("required"); required == "true" {
		info.Required = true
	}

	return info
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

// parseImageMediaType converts short image type to full media type.
// png -> image/png, jpg/jpeg -> image/jpeg, etc.
func parseImageMediaType(imageType string) string {
	switch strings.ToLower(imageType) {
	case "png":
		return "image/png"
	case "jpg", "jpeg":
		return "image/jpeg"
	case "gif":
		return "image/gif"
	case "webp":
		return "image/webp"
	case "svg":
		return "image/svg+xml"
	case "bmp":
		return "image/bmp"
	case "ico":
		return "image/x-icon"
	default:
		// If already in full format or unknown, return as-is
		if strings.HasPrefix(imageType, "image/") {
			return imageType
		}
		return "image/" + imageType
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

