package core

import (
	"encoding/binary"
	"errors"
)

type Instruction byte

const (
	InstrPushInt  Instruction = 0x0a // 10
	InstrAdd      Instruction = 0x0b // 11
	InstrPushByte Instruction = 0x0c // 12
	InstrPack     Instruction = 0x0d // 13
	InstrSub      Instruction = 0x0e // 14
	InstrStore    Instruction = 0x0f // 15
)

type Stack struct {
	data []interface{}
	sp   int // stack pointer
}

func NewStack(size int) *Stack {
	return &Stack{
		data: make([]interface{}, 0, size),
		sp:   0,
	}
}

func (s *Stack) Push(v interface{}) {
	s.data = append(s.data, v)
	s.sp++
}

func (s *Stack) Pop() interface{} {
	if s.sp <= 0 {
		panic("stack underflow")
	}
	value := s.data[s.sp-1]
	s.data = s.data[:s.sp-1]
	s.sp--
	return value
}

type VM struct {
	data          []byte   // VM bytecode
	ip            int      // instruction pointer
	stack         *Stack   // VM stack
	contractState *State   // contract state
}

func NewVM(data []byte, contractState *State) *VM {
	return &VM{
		data:          data,
		ip:            0,
		stack:         NewStack(128),
		contractState: contractState,
	}
}

func (vm *VM) Run() error {
	for vm.ip < len(vm.data) {
		instr := Instruction(vm.data[vm.ip])

		if err := vm.Exec(instr); err != nil {
			return err
		}

		vm.ip++
	}

	return nil
}

func (vm *VM) Exec(instr Instruction) error {
	switch instr {
	case InstrStore:
		key := vm.stack.Pop().([]byte)
		value := vm.stack.Pop()

		switch v := value.(type) {
		case int:
			serializedValue := serializeInt64(int64(v))
			vm.contractState.Put(key, serializedValue)
		default:
			return errors.New("unsupported type for storage")
		}

	case InstrPushInt:
		vm.stack.Push(int(vm.data[vm.ip-1]))

	case InstrPushByte:
		vm.stack.Push(byte(vm.data[vm.ip-1]))

	case InstrPack:
		n := vm.stack.Pop().(int)
		if n < 0 {
			return errors.New("negative length in InstrPack")
		}

		b := make([]byte, n)
		for i := 0; i < n; i++ {
			b[i] = vm.stack.Pop().(byte)
		}

		vm.stack.Push(b)

	case InstrSub:
		a := vm.stack.Pop().(int)
		b := vm.stack.Pop().(int)
		vm.stack.Push(a - b)

	case InstrAdd:
		a := vm.stack.Pop().(int)
		b := vm.stack.Pop().(int)
		vm.stack.Push(a + b)

	default:
		return fmt.Errorf("unknown instruction: %v", instr)
	}

	return nil
}

func serializeInt64(value int64) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(value))
	return buf
}
