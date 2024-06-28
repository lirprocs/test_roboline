package main

import (
	"testing"
)

func TestCrc8(t *testing.T) {
	tests := []struct {
		input    []byte
		expected byte
	}{
		{[]byte{0x01, 0x03, 0x01}, 0x53},
		{[]byte{0x3A, 0x01, 0x10}, 0x03},
		{[]byte{0x0C, 0x05, 0x04, 0x00}, 0x7C},
	}

	for _, tt := range tests {
		result := crc8(tt.input)
		if result != tt.expected {
			t.Errorf("crc8(%v) = %X; want %X", tt.input, result, tt.expected)
		}
	}
}

func TestParseMessage(t *testing.T) {
	tests := []struct {
		message  []uint8
		expected [5]string
	}{
		{
			[]uint8{0x01, 0x03, 0x01, 0x53},
			[5]string{
				"Адрес устройства: 1",
				"Код функции: 3 (Чтение байт)",
				"Адрес ячейки памяти, откуда идёт чтение: 1",
				"Результат проверки контрольной суммы: Корректно",
				"",
			},
		},
		{
			[]uint8{0x3A, 0x01, 0x10, 0x03},
			[5]string{
				"Адрес устройства: 58",
				"Код функции: 1 (Чтение бит)",
				"Адрес ячейки памяти, откуда идёт чтение: 16",
				"Результат проверки контрольной суммы: Корректно",
				"",
			},
		},
		{
			[]uint8{0xAB, 0x06, 0x4C, 0x13, 0xA8},
			[5]string{
				"Адрес устройства: 171",
				"Код функции: 6 (Запись байт)",
				"Адрес ячейки памяти, куда идёт запись: 76",
				"Записываемое значение: 19",
				"Результат проверки контрольной суммы: Некорректно",
			},
		},
		{
			[]uint8{0x0C, 0x05, 0x04, 0x00, 0x7C},
			[5]string{
				"Адрес устройства: 12",
				"Код функции: 5 (Запись бит)",
				"Адрес ячейки памяти, куда идёт запись: 4",
				"Записываемое значение: 0",
				"Результат проверки контрольной суммы: Корректно",
			},
		},
		{
			[]uint8{0x01, 0x01, 0x01, 0x01, 0x79},
			[5]string{
				"Ошибка: Некорректный формат данных",
			},
		},
		{
			[]uint8{0x01, 0x06, 0x01, 0x01},
			[5]string{
				"Ошибка: Некорректный формат данных",
			},
		},
		{
			[]uint8{0x01, 0x02, 0x01, 0x01},
			[5]string{
				"Ошибка: Неизвестная операция",
			},
		},
		{
			[]uint8{0x01, 0x03, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
			[5]string{
				"Ошибка: Некорректный формат данных",
			},
		},
	}

	for _, tt := range tests {
		var d data
		d.deviceAddress = tt.message[0]
		d.functionCode = tt.message[1]
		d.memoryAddress = tt.message[2]
		d.crc = tt.message[len(tt.message)-1]

		d.parseMessage(tt.message)
		for i, result := range d.result {
			if result != tt.expected[i] {
				t.Errorf("parseMessage(%v) = %v; want %v", tt.message, d.result, tt.expected)
				break
			}
		}
	}
}
