package encoding

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

func Unmarshal(data []byte, c *Config) error {
	if len(data) < 7 {
		return errors.New("binary data too short to contain a valid layout header")
	}

	if *c == nil {
		*c = make(Config)
	}

	reader := bytes.NewReader(data)

	// Magic
	var magic [4]byte
	if _, err := io.ReadFull(reader, magic[:]); err != nil {
		return fmt.Errorf("failed to read magic signature: %w", err)
	}
	if magic != Magic {
		return fmt.Errorf("invalid file format: expected magic %q, got %q", Magic, magic)
	}

	// Version
	var version uint8
	if err := binary.Read(reader, binary.BigEndian, &version); err != nil {
		return fmt.Errorf("failed to read file version: %w", err)
	}
	if version > CurrentVersion {
		return fmt.Errorf("unsupported file version: file is version %d, but application only supports up to version %d", version, CurrentVersion)
	}

	// Number of hosts
	var numHosts uint16
	if err := binary.Read(reader, binary.BigEndian, &numHosts); err != nil {
		return fmt.Errorf("failed to read map size: %w", err)
	}

	for i := uint16(0); i < numHosts; i++ {
		// Host name
		var nameLen uint8
		if err := binary.Read(reader, binary.BigEndian, &nameLen); err != nil {
			return fmt.Errorf("failed to read name length at record %d: %w", i, err)
		}
		nameBytes := make([]byte, nameLen)
		if _, err := io.ReadFull(reader, nameBytes); err != nil {
			return fmt.Errorf("failed to read name data at record %d: %w", i, err)
		}
		nameStr := string(nameBytes)

		// Number of records
		var numRecords uint8
		if err := binary.Read(reader, binary.BigEndian, &numRecords); err != nil {
			return fmt.Errorf("failed to read map size: %w", err)
		}

		for j := uint8(0); j < numRecords; j++ {
			// Key
			var keyLen uint8
			if err := binary.Read(reader, binary.BigEndian, &keyLen); err != nil {
				return fmt.Errorf("failed to read key length at record %d: %w", j, err)
			}

			// Value
			var valLen uint16
			if err := binary.Read(reader, binary.BigEndian, &valLen); err != nil {
				return fmt.Errorf("failed to read value length for key %d: %w", keyLen, err)
			}
			valBytes := make([]byte, valLen)
			if _, err := io.ReadFull(reader, valBytes); err != nil {
				return fmt.Errorf("failed to read value data for key %d: %w", keyLen, err)
			}

			if (*c)[nameStr] == nil {
				(*c)[nameStr] = make(map[uint8]string)
			}

			(*c)[nameStr][keyLen] = string(valBytes)
		}
	}

	return nil
}
