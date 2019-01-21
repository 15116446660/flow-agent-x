package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/streadway/amqp"

	"github.com/flowci/flow-agent-x/domain"
	"github.com/flowci/flow-agent-x/util"
)

const (
	errSettingConnectFail = "Cannot get settings from server"
	defaultWorkspace      = "${HOME}/.flow.ci.agent"
	defaultLoggingDir     = defaultWorkspace + "/logs"
	defaultPluginDir      = defaultWorkspace + "/plugins"
)

var (
	singleton *Manager
	once      sync.Once
)

type QueueConfig struct {
	Conn          *amqp.Connection
	Channel       *amqp.Channel
	JobQueue      *amqp.Queue
	CallbackQueue *amqp.Queue
}

// Manager to handle server connection and config
type Manager struct {
	Settings   *domain.Settings
	Queue      *QueueConfig
	IsOffline  bool
	Workspace  string
	LoggingDir string
	PluginDir  string
}

// GetInstance get singleton of config manager
func GetInstance() *Manager {
	once.Do(func() {
		singleton = new(Manager)
		singleton.IsOffline = false
		singleton.Workspace = defaultWorkspace
		singleton.LoggingDir = defaultLoggingDir
		singleton.PluginDir = defaultPluginDir
	})
	return singleton
}

func (m *Manager) Init() error {
	var err = m.initSettings()
	if err != nil {
		util.LogInfo("Model: 'offline'")
		singleton.IsOffline = true
		return err
	}

	err = m.initRabbitMQ()
	if err != nil {
		return err
	}

	return nil
}

// HasQueue has rabbit mq connected
func (m *Manager) HasQueue() bool {
	return m.Queue != nil
}

// Close release resources and connections
func (m *Manager) Close() {
	if m.HasQueue() {
		m.Queue.Channel.Close()
		m.Queue.Conn.Close()
	}
}

func (m *Manager) initSettings() error {
	server, token, port := getVaraibles()

	uri := server + "/agents/connect"
	body, _ := json.Marshal(domain.AgentConnect{Token: token, Port: port})

	var message domain.SettingsResponse
	resp, errFromReq := http.Post(uri, util.HttpMimeJson, bytes.NewBuffer(body))
	if errFromReq != nil {
		return fmt.Errorf(errSettingConnectFail)
	}

	defer resp.Body.Close()
	raw, _ := ioutil.ReadAll(resp.Body)
	errFromJSON := json.Unmarshal(raw, &message)

	if errFromJSON != nil {
		return errFromJSON
	}

	if !message.IsOk() {
		return fmt.Errorf(message.Message)
	}

	m.Settings = message.Data
	return nil
}

func (m *Manager) initRabbitMQ() error {
	if m.Settings == nil {
		return ErrSettingsNotBeenLoaded
	}

	// get connection
	connStr := m.Settings.Queue.GetConnectionString()
	conn, err := amqp.Dial(connStr)
	if err != nil {
		return err
	}

	// get channel
	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	// init queue config
	qc := new(QueueConfig)
	qc.Conn = conn
	qc.Channel = ch

	// init queue to receive job
	jobQueue, err := ch.QueueDeclare(m.Settings.Agent.GetQueueName(), true, false, false, false, nil)
	qc.JobQueue = &jobQueue

	// init queue to send callback
	callbackQueue, err := ch.QueueDeclare(m.Settings.CallbackQueueName, true, false, false, false, nil)
	qc.CallbackQueue = &callbackQueue

	m.Queue = qc
	return nil
}

func (m *Manager) initZookeeper() error {
	if m.Settings == nil {
		return ErrSettingsNotBeenLoaded
	}

	zkConfig := m.Settings.Zookeeper

	// make connection of zk
	client := new(util.ZkClient)
	err := client.Connect(zkConfig.Host)

	if err != nil {
		return err
	}

	// register agent on zk
	agentPath := getZkPath(m.Settings)
	_, nodeErr := client.Create(agentPath, util.ZkNodeTypeEphemeral, string(domain.AgentIdle))

	return nodeErr
}

func getZkPath(s *domain.Settings) string {
	return s.Zookeeper.Root + "/" + s.Agent.ID
}

func getVaraibles() (server string, token string, port int) {
	server = os.Getenv("FLOWCI_SERVER_URL")
	token = os.Getenv("FLOWCI_AGENT_TOKEN")
	port, _ = strconv.Atoi(os.Getenv("FLOWCI_AGENT_PORT"))
	return
}
