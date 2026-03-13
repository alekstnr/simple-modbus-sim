# Simple Modbus Simulator

A simple Modbus TCP server simulator that reads configuration from a YAML file and serves defined coils, discrete inputs, holding registers, and input registers.

## Installation

```bash
go build -o modbus-sim .
```

## Usage

```bash
./modbus-sim [config.yaml]
```

If no config file is specified, it defaults to `config.yaml`.

## Configuration

### Example Config

```yaml
server:
  url: "tcp://localhost:5502" # Server listen address
  timeout: 30 # Connection timeout in seconds
  max_clients: 5 # Max concurrent connections
  allow_undefined_registers: true # Return 0 for undefined addresses
devices:
  - name: "device1"
    modbus_id: 1
    data:
      coils:
        "0": true
        "1": false
        "10": true
      discrete_inputs:
        "0": true
        "1": false
        "100": true
      holding_registers:
        "0": 100
        "1": 200
        "100": 1234
      input_registers:
        "0": 1000
        "1": 2000
        "100": 5000
  - name: "device2"
    modbus_id: 2
    data:
      coils:
        "0": false
        "1": true
      holding_registers:
        "0": 1000
        "1": 2000
      input_registers:
        "0": 5000
        "1": 6000
```

### Server Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `url` | string | `tcp://localhost:5502` | Listen address |
| `timeout` | int | 30 | Connection timeout (seconds) |
| `max_clients` | int | 5 | Max concurrent connections |
| `allow_undefined_registers` | bool | false | Return 0 for undefined addresses |

### Device Options

| Option | Type | Description |
|--------|------|-------------|
| `name` | string | Device name (for display) |
| `modbus_id` | uint8 | Modbus unit ID |
| `data.coils` | map | Coil addresses (bool) |
| `data.discrete_inputs` | map | Discrete input addresses (bool) |
| `data.holding_registers` | map | Holding register addresses (int) |
| `data.input_registers` | map | Input register addresses (int) |

## Testing

You can test the server using the modbus-cli tool:

```bash
# Build the CLI
go build -o modbus-cli github.com/simonvetter/modbus/cmd/modbus-cli

# Read coils
./modbus-cli --target tcp://localhost:5502 rc:0+10

# Read holding registers
./modbus-cli --target tcp://localhost:5502 rhr:0+10

# Read input registers
./modbus-cli --target tcp://localhost:5502 rir:0+10

# Read discrete inputs
./modbus-cli --target tcp://localhost:5502 rdi:0+10

# Write a coil
./modbus-cli --target tcp://localhost:5502 wc:0:true

# Write a holding register
./modbus-cli --target tcp://localhost:5502 wh:0:1234
```

## Notes

- Address keys in the YAML are strings (e.g., "0", "100") to avoid issues with leading zeros
- Modbus addresses are 0-indexed in this simulator
- Multiple devices are supported with different unit IDs

## TODOs

[ ] add logging
