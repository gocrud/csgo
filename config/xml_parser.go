package config

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

// XmlConfigurationSource represents an XML file configuration source.
type XmlConfigurationSource struct {
	Path           string
	Optional       bool
	ReloadOnChange bool
}

// Build builds an IConfigurationProvider from the XML source.
func (s *XmlConfigurationSource) Build(builder IConfigurationBuilder) IConfigurationProvider {
	return &XmlConfigurationProvider{
		source: s,
	}
}

// XmlConfigurationProvider is a provider for XML configuration files.
type XmlConfigurationProvider struct {
	*ConfigurationProvider
	source  *XmlConfigurationSource
	watcher *FileWatcher
}

// Load loads configuration from XML file.
func (p *XmlConfigurationProvider) Load() map[string]string {
	s := p.source
	if p.ConfigurationProvider == nil {
		p.ConfigurationProvider = NewConfigurationProvider()
	}

	data := make(map[string]string)

	file, err := os.Open(s.Path)
	if err != nil {
		if s.Optional && os.IsNotExist(err) {
			return data
		}
		if os.IsNotExist(err) {
			return data
		}
		panic(fmt.Sprintf("failed to read XML config file %s: %v", s.Path, err))
	}
	defer file.Close()

	// Parse XML file
	xmlData := parseXmlFile(file)
	
	// Flatten to key:value format
	flattenXmlMap("", xmlData, data)

	// Store data in base provider
	p.SetData(data)
	
	// Setup file watching if enabled
	if s.ReloadOnChange && p.watcher == nil {
		p.watcher = NewFileWatcher(s.Path, func() {
			newData := p.Load()
			p.SetData(newData)
		})
	}

	return data
}

// xmlNode represents a node in the XML structure.
type xmlNode struct {
	XMLName  xml.Name
	Attrs    []xml.Attr `xml:",any,attr"`
	Content  string     `xml:",chardata"`
	Children []xmlNode  `xml:",any"`
}

// parseXmlFile parses an XML file into a nested map structure.
func parseXmlFile(r io.Reader) map[string]interface{} {
	decoder := xml.NewDecoder(r)
	
	var root xmlNode
	if err := decoder.Decode(&root); err != nil {
		return make(map[string]interface{})
	}

	result := make(map[string]interface{})
	processXmlNode(&root, result)
	
	return result
}

// processXmlNode processes an XML node and adds it to the result map.
func processXmlNode(node *xmlNode, result map[string]interface{}) {
	// If node has no children, it's a leaf node
	if len(node.Children) == 0 {
		content := strings.TrimSpace(node.Content)
		if content != "" {
			result[node.XMLName.Local] = content
		}
		return
	}

	// Process child nodes
	childMap := make(map[string]interface{})
	childArrays := make(map[string][]interface{})
	
	for _, child := range node.Children {
		childName := child.XMLName.Local
		
		// Check if this is a repeated element (array)
		if existing, exists := childMap[childName]; exists {
			// Convert to array if not already
			if _, isArray := childArrays[childName]; !isArray {
				childArrays[childName] = []interface{}{existing}
				delete(childMap, childName)
			}
			
			// Process the child
			if len(child.Children) == 0 {
				content := strings.TrimSpace(child.Content)
				if content != "" {
					childArrays[childName] = append(childArrays[childName], content)
				}
			} else {
				nestedMap := make(map[string]interface{})
				processXmlNode(&child, nestedMap)
				childArrays[childName] = append(childArrays[childName], nestedMap)
			}
		} else {
			// First occurrence
			if len(child.Children) == 0 {
				content := strings.TrimSpace(child.Content)
				if content != "" {
					childMap[childName] = content
				}
			} else {
				nestedMap := make(map[string]interface{})
				processXmlNode(&child, nestedMap)
				childMap[childName] = nestedMap
			}
		}
	}

	// Merge childMap and childArrays into result
	for k, v := range childMap {
		result[k] = v
	}
	for k, v := range childArrays {
		result[k] = v
	}
}

// flattenXmlMap flattens the XML map to configuration key:value format.
func flattenXmlMap(prefix string, src map[string]interface{}, dst map[string]string) {
	for key, value := range src {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + ":" + key
		}

		switch v := value.(type) {
		case map[string]interface{}:
			// Recursively process nested objects
			flattenXmlMap(fullKey, v, dst)
		case []interface{}:
			// Process arrays
			for i, item := range v {
				arrayKey := fmt.Sprintf("%s:%d", fullKey, i)
				if nested, ok := item.(map[string]interface{}); ok {
					flattenXmlMap(arrayKey, nested, dst)
				} else {
					dst[arrayKey] = fmt.Sprintf("%v", item)
				}
			}
		case string:
			dst[fullKey] = v
		default:
			// Convert other types to string
			dst[fullKey] = fmt.Sprintf("%v", v)
		}
	}
}

