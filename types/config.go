package types

type Config struct {
	Server  ServerConfig   `yaml:"server"`
	Devices []DeviceConfig `yaml:"devices"`
}

type ServerConfig struct {
	URL                     string `yaml:"url"`
	Timeout                 int    `yaml:"timeout"`
	MaxClients              int    `yaml:"max_clients"`
	AllowUndefinedRegisters bool   `yaml:"allow_undefined_registers"`
}

type DeviceConfig struct {
	Name     string     `yaml:"name"`
	ModbusId uint8      `yaml:"modbus_id"`
	Data     DataConfig `yaml:"data"`
}

type DataConfig struct {
	Coils            map[string]bool `yaml:"coils"`
	DiscreteInputs   map[string]bool `yaml:"discrete_inputs"`
	HoldingRegisters map[string]int  `yaml:"holding_registers"`
	InputRegisters   map[string]int  `yaml:"input_registers"`
}
