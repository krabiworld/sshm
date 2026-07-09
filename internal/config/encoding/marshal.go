package encoding

import (
	"bytes"
	"encoding/binary"
)

func Marshal(c Config) ([]byte, error) {
	buf := new(bytes.Buffer)

	// Magic
	if _, err := buf.Write(Magic[:]); err != nil {
		return nil, err
	}
	// Version
	if err := binary.Write(buf, binary.BigEndian, CurrentVersion); err != nil {
		return nil, err
	}
	// Number of hosts
	if err := binary.Write(buf, binary.BigEndian, uint16(len(c))); err != nil {
		return nil, err
	}

	for host, records := range c {
		// Host name
		nameBytes := []byte(host)
		if err := binary.Write(buf, binary.BigEndian, uint8(len(nameBytes))); err != nil {
			return nil, err
		}
		if _, err := buf.Write(nameBytes); err != nil {
			return nil, err
		}

		// Number of records
		if err := binary.Write(buf, binary.BigEndian, uint8(len(records))); err != nil {
			return nil, err
		}

		for key, value := range records {
			// Record type
			if err := binary.Write(buf, binary.BigEndian, key); err != nil {
				return nil, err
			}

			// Record value
			valBytes := []byte(value)
			if err := binary.Write(buf, binary.BigEndian, uint16(len(valBytes))); err != nil {
				return nil, err
			}
			if _, err := buf.Write(valBytes); err != nil {
				return nil, err
			}
		}
	}

	return buf.Bytes(), nil
}
