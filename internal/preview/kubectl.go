package preview

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

func EnsureNamespace(name string) error {
	create := exec.Command("kubectl", "create", "ns", name, "--dry-run=client", "-o", "yaml")
	var yamlOutput bytes.Buffer
	var createErr bytes.Buffer
	create.Stdout = &yamlOutput
	create.Stderr = &createErr

	if err := create.Run(); err != nil {
		if createErr.Len() > 0 {
			log.Printf("kubectl create stderr: %s", createErr.String())
		}
		return fmt.Errorf("generate ns yaml failed: %w", err)
	}

	apply := exec.Command("kubectl", "apply", "-f", "-")
	apply.Stdin = &yamlOutput
	var applyOutput bytes.Buffer
	var applyErr bytes.Buffer
	apply.Stdout = &applyOutput
	apply.Stderr = &applyErr
	if err := apply.Run(); err != nil {
		if applyOutput.Len() > 0 {
			log.Printf("kubectl apply stdout: %s", applyOutput.String())
		}
		if applyErr.Len() > 0 {
			log.Printf("kubectl apply stderr: %s", applyErr.String())
		}
		return fmt.Errorf("apply ns failed: %w", err)
	}

	if applyOutput.Len() > 0 {
		log.Printf("namespace secured: %s", applyOutput.String())
	}
	return nil
}

func DeleteNamespace(name string) error {
	deletecmd := exec.Command("kubectl", "delete", "ns", name, "--ignore-not-found")
	var deleteOutput bytes.Buffer
	var deleteErr bytes.Buffer
	deletecmd.Stdout = &deleteOutput
	deletecmd.Stderr = &deleteErr

	if err := deletecmd.Run(); err != nil {
		if deleteErr.Len() > 0 {
			log.Printf("kubectl delete stderr: %s", deleteErr.String())
		}
		return fmt.Errorf("delete ns failed: %w", err)
	}
	if deleteOutput.Len() > 0 {
		log.Printf("namespace deleted: %s", deleteOutput.String())
	}
	return nil
}
