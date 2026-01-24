package docker

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func GenerateOverriddenCompose(sourceComposePath string, outTmpComposePath string) error {
	composeData, err := AddAlfredNetwork(sourceComposePath)
	if err == nil {
		return SaveComposeFile(composeData, outTmpComposePath)
	}
	return err
}

func AddAlfredNetwork(composePath string) (map[string]interface{}, error) {
	composeData, _ := Load(composePath)
	networks, ok := composeData["networks"].(map[string]interface{})
	if !ok || networks == nil {
		networks = make(map[string]interface{})
		composeData["networks"] = networks
	}

	// Add minimal network definition
	networks["alfred-dev"] = map[string]interface{}{
		"external": true, // reuse existing Docker network
	}
	services, ok := composeData["services"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Could not find any services in docker-compose.")
	}

	for _, s := range services {
		svc, ok := s.(map[string]interface{})
		if !ok {
			continue
		}

		networks, ok := svc["networks"].(map[string]interface{})
		if !ok || networks == nil {
			networks = make(map[string]interface{})
			svc["networks"] = networks
		}

		networks["alfred-dev"] = map[string]interface{}{} // minimal attach
	}
	return composeData, nil
}

func SaveComposeFile(compose map[string]interface{}, path string) error {
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return err
	}
	data, err := yaml.Marshal(compose)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func Up(composePath string, projectName string) error {
	// Run docker-compose
	cmd := exec.Command(
		"docker", "compose",
		"-p", projectName,
		"-f", composePath, "up", "-d",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Running:", cmd.String())
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func CreateDevNetwork(networkName string) error {
	cmd := exec.Command(
		"docker", "network",
		"create", networkName,
	)
	cmd.Stdout = os.Stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	fmt.Println("Created Network alfred-dev")
	err := cmd.Run()
	if err != nil {
		msg := stderr.String()
		if strings.Contains(strings.ToLower(msg), "already exists") {
			fmt.Println("Network already exists, skipping creation.")
			return nil
		} else {
			return err
		}
	}
	return nil
}

func Load(composePath string) (map[string]interface{}, error) {
	data, err := os.ReadFile(composePath)
	if err != nil {
		return nil, err
	}

	var compose map[string]interface{}
	if err := yaml.Unmarshal(data, &compose); err != nil {
		return nil, err
	}
	return compose, err
}
