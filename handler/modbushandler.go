package handler

import (
	"strconv"
	"sync"

	"github.com/alekstnr/simple-modbus-sim/types"
	"github.com/simonvetter/modbus"
)

type deviceData struct {
	coils            map[uint16]bool
	discreteInputs   map[uint16]bool
	holdingRegisters map[uint16]uint16
	inputRegisters   map[uint16]uint16
}

type ModbusHandler struct {
	lock                    sync.RWMutex
	devices                 map[uint8]*deviceData
	allowUndefinedRegisters bool
}

func NewModbusHandler(devices []types.DeviceConfig, allowUndefined bool) *ModbusHandler {
	h := &ModbusHandler{
		devices:                 make(map[uint8]*deviceData),
		allowUndefinedRegisters: allowUndefined,
	}

	for _, d := range devices {
		data := &deviceData{
			coils:            make(map[uint16]bool),
			discreteInputs:   make(map[uint16]bool),
			holdingRegisters: make(map[uint16]uint16),
			inputRegisters:   make(map[uint16]uint16),
		}

		for addrStr, value := range d.Data.Coils {
			addr, err := strconv.ParseUint(addrStr, 10, 16)
			if err == nil {
				data.coils[uint16(addr)] = value
			}
		}

		for addrStr, value := range d.Data.DiscreteInputs {
			addr, err := strconv.ParseUint(addrStr, 10, 16)
			if err == nil {
				data.discreteInputs[uint16(addr)] = value
			}
		}

		for addrStr, value := range d.Data.HoldingRegisters {
			addr, err := strconv.ParseUint(addrStr, 10, 16)
			if err == nil {
				data.holdingRegisters[uint16(addr)] = uint16(value)
			}
		}

		for addrStr, value := range d.Data.InputRegisters {
			addr, err := strconv.ParseUint(addrStr, 10, 16)
			if err == nil {
				data.inputRegisters[uint16(addr)] = uint16(value)
			}
		}

		h.devices[d.ModbusId] = data
	}

	return h
}

func (h *ModbusHandler) getDevice(unitId uint8) (*deviceData, error) {
	h.lock.RLock()
	defer h.lock.RUnlock()

	if device, ok := h.devices[unitId]; ok {
		return device, nil
	}
	return nil, modbus.ErrBadUnitId
}

func (h *ModbusHandler) HandleCoils(req *modbus.CoilsRequest) ([]bool, error) {
	device, err := h.getDevice(req.UnitId)
	if err != nil {
		return nil, err
	}

	if req.IsWrite {
		h.lock.Lock()
		defer h.lock.Unlock()
		for i, val := range req.Args {
			addr := req.Addr + uint16(i)
			device.coils[addr] = val
		}
		return nil, nil
	}

	h.lock.RLock()
	defer h.lock.RUnlock()

	res := make([]bool, req.Quantity)
	for i := uint16(0); i < req.Quantity; i++ {
		addr := req.Addr + i
		if val, ok := device.coils[addr]; ok {
			res[i] = val
		} else if h.allowUndefinedRegisters {
			res[i] = false
		} else {
			return nil, modbus.ErrIllegalDataAddress
		}
	}
	return res, nil
}

func (h *ModbusHandler) HandleDiscreteInputs(req *modbus.DiscreteInputsRequest) ([]bool, error) {
	device, err := h.getDevice(req.UnitId)
	if err != nil {
		return nil, err
	}

	h.lock.RLock()
	defer h.lock.RUnlock()

	res := make([]bool, req.Quantity)
	for i := uint16(0); i < req.Quantity; i++ {
		addr := req.Addr + i
		if val, ok := device.discreteInputs[addr]; ok {
			res[i] = val
		} else if h.allowUndefinedRegisters {
			res[i] = false
		} else {
			return nil, modbus.ErrIllegalDataAddress
		}
	}
	return res, nil
}

func (h *ModbusHandler) HandleHoldingRegisters(req *modbus.HoldingRegistersRequest) ([]uint16, error) {
	device, err := h.getDevice(req.UnitId)
	if err != nil {
		return nil, err
	}

	if req.IsWrite {
		h.lock.Lock()
		defer h.lock.Unlock()
		for i, val := range req.Args {
			addr := req.Addr + uint16(i)
			device.holdingRegisters[addr] = val
		}
		return nil, nil
	}

	h.lock.RLock()
	defer h.lock.RUnlock()

	res := make([]uint16, req.Quantity)
	for i := uint16(0); i < req.Quantity; i++ {
		addr := req.Addr + i
		if val, ok := device.holdingRegisters[addr]; ok {
			res[i] = val
		} else if h.allowUndefinedRegisters {
			res[i] = 0
		} else {
			return nil, modbus.ErrIllegalDataAddress
		}
	}
	return res, nil
}

func (h *ModbusHandler) HandleInputRegisters(req *modbus.InputRegistersRequest) ([]uint16, error) {
	device, err := h.getDevice(req.UnitId)
	if err != nil {
		return nil, err
	}

	h.lock.RLock()
	defer h.lock.RUnlock()

	res := make([]uint16, req.Quantity)
	for i := uint16(0); i < req.Quantity; i++ {
		addr := req.Addr + i
		if val, ok := device.inputRegisters[addr]; ok {
			res[i] = val
		} else if h.allowUndefinedRegisters {
			res[i] = 0
		} else {
			return nil, modbus.ErrIllegalDataAddress
		}
	}
	return res, nil
}
