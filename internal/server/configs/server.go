package configs

// Wiregost - Post-Exploitation & Implant Framework
// Copyright © 2020 Para
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

import (
	"encoding/hex"
	"encoding/json"
	insecureRand "math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/maxlandon/wiregost/internal/server/assets"
	"github.com/maxlandon/wiregost/internal/server/log"
)

const (
	serverConfigFileName = "server.json"
)

var serverConfigLog = log.NamedLogger("config", "server")

// GetServerConfigPath - File path to config.json
func GetServerConfigPath() string {
	appDir := assets.GetRootAppDir()
	serverConfigPath := filepath.Join(appDir, "configs", serverConfigFileName)
	serverConfigLog.Debugf("Loading config from %s", serverConfigPath)
	return serverConfigPath
}

// LogConfig - Server logging config
type LogConfig struct {
	Level              int  `json:"level"`
	GRPCUnaryPayloads  bool `json:"grpc_unary_payloads"`
	GRPCStreamPayloads bool `json:"grpc_stream_payloads"`
	TLSKeyLogger       bool `json:"tls_key_logger"`
}

// DaemonConfig - Configure daemon mode
type DaemonConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// JobConfig - Restart Jobs on Load
type JobConfig struct {
	Multiplayer []*MultiplayerJobConfig `json:"multiplayer"`
}

type MultiplayerJobConfig struct {
	Host  string `json:"host"`
	Port  uint16 `json:"port"`
	JobID string `json:"job_id"`
}

// ServerConfig - Server config
type ServerConfig struct {
	DaemonMode   bool          `json:"daemon_mode"`
	DaemonConfig *DaemonConfig `json:"daemon"`
	Logs         *LogConfig    `json:"logs"`
	Jobs         *JobConfig    `json:"jobs,omitempty"`
	GoProxy      string        `json:"go_proxy"`
}

// Save - Save config file to disk
func (c *ServerConfig) Save() error {
	configPath := GetServerConfigPath()
	configDir := filepath.Dir(configPath)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		serverConfigLog.Debugf("Creating config dir %s", configDir)
		err := os.MkdirAll(configDir, 0o700)
		if err != nil {
			return err
		}
	}
	data, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return err
	}
	serverConfigLog.Infof("Saving config to %s", configPath)
	err = os.WriteFile(configPath, data, 0o600)
	if err != nil {
		serverConfigLog.Errorf("Failed to write config %s", err)
	}
	return nil
}

// AddMultiplayerJob - Add Job Configs
func (c *ServerConfig) AddMultiplayerJob(config *MultiplayerJobConfig) error {
	if c.Jobs == nil {
		c.Jobs = &JobConfig{}
	}
	config.JobID = getRandomID()
	c.Jobs.Multiplayer = append(c.Jobs.Multiplayer, config)
	return c.Save()
}

// RemoveJob - Remove Job by ID
func (c *ServerConfig) RemoveJob(jobID string) {
	if c.Jobs == nil {
		return
	}
	defer c.Save()
}

// GetServerConfig - Get config value
func GetServerConfig() *ServerConfig {
	configPath := GetServerConfigPath()
	config := getDefaultServerConfig()
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		data, err := os.ReadFile(configPath)
		if err != nil {
			serverConfigLog.Errorf("Failed to read config file %s", err)
			return config
		}
		err = json.Unmarshal(data, config)
		if err != nil {
			serverConfigLog.Errorf("Failed to parse config file %s", err)
			return config
		}
	} else {
		serverConfigLog.Warnf("Config file does not exist, using defaults")
	}

	if config.Logs.Level < 0 {
		config.Logs.Level = 0
	}
	if 6 < config.Logs.Level {
		config.Logs.Level = 6
	}
	log.RootLogger.SetLevel(log.LevelFrom(config.Logs.Level))

	err := config.Save() // This updates the config with any missing fields
	if err != nil {
		serverConfigLog.Errorf("Failed to save default config %s", err)
	}
	return config
}

func getDefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		DaemonMode: false,
		DaemonConfig: &DaemonConfig{
			Host: "",
			Port: 31337,
		},
		Logs: &LogConfig{
			Level:              int(logrus.InfoLevel),
			GRPCUnaryPayloads:  false,
			GRPCStreamPayloads: false,
		},
		Jobs: &JobConfig{},
	}
}

func getRandomID() string {
	seededRand := insecureRand.New(insecureRand.NewSource(time.Now().UnixNano()))
	buf := make([]byte, 32)
	seededRand.Read(buf)
	return hex.EncodeToString(buf)
}
