package ansible

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"

	corePort "starliner.app/internal/core/domain/port"
	"starliner.app/internal/provisioner/domain/port"
)

type K3sPlaybookOutput struct {
	Plays []struct {
		Tasks []struct {
			Hosts map[string]struct {
				Content string `json:"content,omitempty"`
			} `json:"hosts"`
		} `json:"tasks"`
	} `json:"plays"`
}

type Install struct {
	crypto       corePort.Crypto
	logPublisher port.LogPublisher
}

func NewInstall(
	crypto corePort.Crypto,
	logPublisher port.LogPublisher,
) port.Install {
	return &Install{
		crypto:       crypto,
		logPublisher: logPublisher,
	}
}

func (i *Install) InstallK3s(provisioningId string, ip string, privateKey []byte) (kubeconfig string, err error) {
	tmpPlaybook, err := os.CreateTemp("", "ansible-*.yml")
	if err != nil {
		return "", err
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			fmt.Printf("Failed to remove temp file: %v\n", err)
		}
	}(tmpPlaybook.Name())

	_, err = tmpPlaybook.Write([]byte(K3sPlaybook))
	if err != nil {
		return "", err
	}
	err = tmpPlaybook.Close()
	if err != nil {
		return "", err
	}

	tmpPrivateKey, err := os.CreateTemp("", "private-key-*.pem")
	if err != nil {
		return "", err
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			fmt.Printf("Failed to remove temp file: %v\n", err)
		}
	}(tmpPrivateKey.Name())

	pemBytes, err := i.crypto.EncodePrivateKeyToPEM(privateKey)
	if err != nil {
		return "", err
	}
	_, err = tmpPrivateKey.Write(pemBytes)
	if err != nil {
		return "", err
	}
	err = tmpPrivateKey.Close()
	if err != nil {
		return "", err
	}

	var (
		logs strings.Builder
		mu   sync.Mutex
	)
	appendLog := func(line string) {
		mu.Lock()
		logs.WriteString(line)
		mu.Unlock()

		if i.logPublisher != nil {
			if err := i.logPublisher.PublishLogChunk(provisioningId, []byte(line)); err != nil {
				fmt.Printf("failed to publish log chunk: %v", err)
			}
		}
	}

	defer func() {
		if i.logPublisher != nil {
			_ = i.logPublisher.PublishLogEnd(provisioningId)
		}
	}()

	args := []string{
		tmpPlaybook.Name(),
		"-i", fmt.Sprintf("%s,", ip),
		"-u", "root",
		"--private-key", tmpPrivateKey.Name(),
		"-e", "ansible_ssh_common_args='-o StrictHostKeyChecking=no'",
		"-e", fmt.Sprintf("target_host=%s", ip),
	}

	cmd := exec.Command("ansible-playbook", args...)
	cmd.Env = append(
		os.Environ(),
		"ANSIBLE_STDOUT_CALLBACK=json",
		"ANSIBLE_DEPRECATION_WARNINGS=False",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}

	if err := cmd.Start(); err != nil {
		return "", err
	}

	var stdoutBuf bytes.Buffer
	var wg sync.WaitGroup

	readPipe := func(r io.Reader, collect *bytes.Buffer) {
		defer wg.Done()

		scanner := bufio.NewScanner(r)
		scanner.Buffer(make([]byte, 0, 64*1024), 10*1024*1024)

		for scanner.Scan() {
			line := scanner.Text() + "\n"
			appendLog(line)

			if collect != nil {
				_, _ = collect.WriteString(line)
			}
		}

		if err := scanner.Err(); err != nil {
			appendLog(fmt.Sprintf("error reading ansible output: %v\n", err))
		}
	}

	wg.Add(2)
	go readPipe(stdout, &stdoutBuf)
	go readPipe(stderr, nil)

	wg.Wait()

	if err := cmd.Wait(); err != nil {
		return "", err
	}

	out := stdoutBuf.Bytes()

	var ansibleData K3sPlaybookOutput
	if err := json.Unmarshal(out, &ansibleData); err != nil {
		return "", err
	}

	var kubeconfigBase64 string
	for _, play := range ansibleData.Plays {
		for _, task := range play.Tasks {
			for _, host := range task.Hosts {
				if host.Content != "" {
					kubeconfigBase64 = host.Content
					break
				}
			}
		}
	}
	kubeconfigDecoded, err := base64.StdEncoding.DecodeString(kubeconfigBase64)
	if err != nil {
		fmt.Printf("failed to decode kubeconfig: %v\n", err)
	}

	kubeconfig = strings.ReplaceAll(
		string(kubeconfigDecoded),
		"https://127.0.0.1:",
		fmt.Sprintf("https://%s:", ip),
	)

	return kubeconfig, nil
}
