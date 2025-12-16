package openapi

// SchemaBuilder provides a fluent API for building JSON Schema.
type SchemaBuilder struct {
	schema Schema
}

// NewSchema creates a new Schema builder.
// Default type is "object" with an empty properties map.
func NewSchema() *SchemaBuilder {
	return &SchemaBuilder{
		schema: Schema{
			Type:       "object",
			Properties: make(map[string]Schema),
		},
	}
}

// Type sets the Schema type (object, array, string, integer, number, boolean).
func (b *SchemaBuilder) Type(t string) *SchemaBuilder {
	b.schema.Type = t
	return b
}

// Description sets the description.
func (b *SchemaBuilder) Description(desc string) *SchemaBuilder {
	b.schema.Description = desc
	return b
}

// Example sets the example value.
func (b *SchemaBuilder) Example(ex interface{}) *SchemaBuilder {
	b.schema.Example = ex
	return b
}

// Property adds an object property.
func (b *SchemaBuilder) Property(name string, schema Schema) *SchemaBuilder {
	if b.schema.Properties == nil {
		b.schema.Properties = make(map[string]Schema)
	}
	b.schema.Properties[name] = schema
	return b
}

// StringProperty adds a string property (convenience method).
func (b *SchemaBuilder) StringProperty(name, description string, example ...string) *SchemaBuilder {
	s := Schema{Type: "string", Description: description}
	if len(example) > 0 {
		s.Example = example[0]
	}
	return b.Property(name, s)
}

// IntProperty adds an integer property (convenience method).
func (b *SchemaBuilder) IntProperty(name, description string, example ...int) *SchemaBuilder {
	s := Schema{Type: "integer", Description: description}
	if len(example) > 0 {
		s.Example = example[0]
	}
	return b.Property(name, s)
}

// NumberProperty adds a number property (convenience method).
func (b *SchemaBuilder) NumberProperty(name, description string, example ...float64) *SchemaBuilder {
	s := Schema{Type: "number", Description: description}
	if len(example) > 0 {
		s.Example = example[0]
	}
	return b.Property(name, s)
}

// BoolProperty adds a boolean property (convenience method).
func (b *SchemaBuilder) BoolProperty(name, description string, example ...bool) *SchemaBuilder {
	s := Schema{Type: "boolean", Description: description}
	if len(example) > 0 {
		s.Example = example[0]
	}
	return b.Property(name, s)
}

// ObjectProperty adds an object property (convenience method).
// The builder function allows nested schema definition.
func (b *SchemaBuilder) ObjectProperty(name, description string, builder func(*SchemaBuilder)) *SchemaBuilder {
	nested := NewSchema()
	builder(nested)
	schema := nested.Build()
	schema.Description = description
	return b.Property(name, schema)
}

// ArrayProperty adds an array property (convenience method).
func (b *SchemaBuilder) ArrayProperty(name, description string, itemSchema Schema) *SchemaBuilder {
	s := Schema{
		Type:        "array",
		Description: description,
		Items:       &itemSchema,
	}
	return b.Property(name, s)
}

// ImageProperty adds a Base64 encoded image property (convenience method).
// This creates a string property with format: byte and contentMediaType for proper Swagger UI display.
// Supported mediaType values: image/png, image/jpeg, image/gif, image/webp, image/svg+xml.
// Example: b.ImageProperty("avatar", "用户头像", "image/png", "iVBORw0KGgo...")
func (b *SchemaBuilder) ImageProperty(name, description, mediaType string, example ...string) *SchemaBuilder {
	s := Schema{
		Type:             "string",
		Format:           "byte",
		ContentMediaType: mediaType,
		Description:      description,
	}
	if len(example) > 0 {
		s.Example = example[0]
	}
	return b.Property(name, s)
}

// Required sets required fields.
func (b *SchemaBuilder) Required(fields ...string) *SchemaBuilder {
	b.schema.Required = append(b.schema.Required, fields...)
	return b
}

// Format sets the format (email, date-time, uri, byte, binary, etc.).
func (b *SchemaBuilder) Format(format string) *SchemaBuilder {
	b.schema.Format = format
	return b
}

// MinLength sets the minimum length (for strings).
func (b *SchemaBuilder) MinLength(min int) *SchemaBuilder {
	b.schema.MinLength = &min
	return b
}

// MaxLength sets the maximum length (for strings).
func (b *SchemaBuilder) MaxLength(max int) *SchemaBuilder {
	b.schema.MaxLength = &max
	return b
}

// Min sets the minimum value (for numbers).
func (b *SchemaBuilder) Min(min float64) *SchemaBuilder {
	b.schema.Minimum = &min
	return b
}

// Max sets the maximum value (for numbers).
func (b *SchemaBuilder) Max(max float64) *SchemaBuilder {
	b.schema.Maximum = &max
	return b
}

// Pattern sets the regex pattern.
func (b *SchemaBuilder) Pattern(pattern string) *SchemaBuilder {
	b.schema.Pattern = pattern
	return b
}

// Enum sets the enum values.
func (b *SchemaBuilder) Enum(values ...interface{}) *SchemaBuilder {
	b.schema.Enum = values
	return b
}

// Nullable marks the schema as nullable.
func (b *SchemaBuilder) Nullable(nullable bool) *SchemaBuilder {
	b.schema.Nullable = nullable
	return b
}

// Build constructs the final Schema.
func (b *SchemaBuilder) Build() Schema {
	return b.schema
}

// ========== Convenience Functions ==========

// StringSchema creates a string Schema.
func StringSchema(description string, example ...string) Schema {
	s := Schema{Type: "string", Description: description}
	if len(example) > 0 {
		s.Example = example[0]
	}
	return s
}

// IntSchema creates an integer Schema.
func IntSchema(description string, example ...int) Schema {
	s := Schema{Type: "integer", Description: description}
	if len(example) > 0 {
		s.Example = example[0]
	}
	return s
}

// NumberSchema creates a number Schema.
func NumberSchema(description string, example ...float64) Schema {
	s := Schema{Type: "number", Description: description}
	if len(example) > 0 {
		s.Example = example[0]
	}
	return s
}

// BoolSchema creates a boolean Schema.
func BoolSchema(description string, example ...bool) Schema {
	s := Schema{Type: "boolean", Description: description}
	if len(example) > 0 {
		s.Example = example[0]
	}
	return s
}

// ArraySchema creates an array Schema.
func ArraySchema(description string, itemSchema Schema) Schema {
	return Schema{
		Type:        "array",
		Description: description,
		Items:       &itemSchema,
	}
}

// ObjectSchema creates an object Schema with properties.
func ObjectSchema(description string, builder func(*SchemaBuilder)) Schema {
	b := NewSchema()
	b.Description(description)
	if builder != nil {
		builder(b)
	}
	return b.Build()
}
