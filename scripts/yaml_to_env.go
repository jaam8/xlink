package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func main() {
	data, err := os.ReadFile("configs/config.yaml")
	if err != nil {
		log.Fatalf("failed to read YAML: %v", err)
	}

	var root yaml.Node
	err = yaml.Unmarshal(data, &root)
	if err != nil {
		log.Fatalf("failed to unmarshal YAML: %v", err)
	}

	if len(root.Content) == 0 {
		log.Fatal("empty YAML content")
	}

	env := make([]string, 0)

	content := root.Content[0]

	for i := 0; i < len(content.Content); i += 2 {
		keyNode := content.Content[i]
		valNode := content.Content[i+1]
		prefix := strings.ToUpper(keyNode.Value) + "_"

		flattenYAML(valNode, prefix, &env)

		// добавляем пустую строку между секциями, кроме последней
		if i+2 < len(content.Content) {
			env = append(env, "")
		}
	}

	err = os.WriteFile("configs/.env", []byte(strings.Join(env, "\n")), 0644)
	if err != nil {
		log.Fatalf("failed to write .env file: %v", err)
	}

	fmt.Println(".env generated 🚀")
}

func flattenYAML(node *yaml.Node, prefix string, env *[]string) {
	if node.Kind != yaml.MappingNode {
		return
	}

	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valNode := node.Content[i+1]
		key := strings.ToUpper(prefix + keyNode.Value)

		switch valNode.Kind {
		case yaml.MappingNode:
			flattenYAML(valNode, key+"_", env)
		case yaml.SequenceNode:
			items := make([]string, 0, len(valNode.Content))
			for _, item := range valNode.Content {
				items = append(items, item.Value)
			}
			*env = append(*env, fmt.Sprintf("%s=%s", key, strings.Join(items, ",")))
		case yaml.ScalarNode:
			*env = append(*env, fmt.Sprintf("%s=%s", key, valNode.Value))
		default:
			log.Printf("unsupported YAML kind: %v", valNode.Kind)
		}
	}
}
