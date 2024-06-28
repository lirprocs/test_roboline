package main

import (
	"fmt"
)

type data struct {
	deviceAddress uint8
	functionCode  uint8
	memoryAddress uint8
	crc           uint8
	ctcResult     string
	result        [5]string
}

func crc8(data []byte) byte {
	var crc byte = 0
	for _, b := range data {
		crc ^= b
		for i := 0; i < 8; i++ {
			if crc&0x80 != 0 {
				crc = (crc << 1) ^ 0x07
			} else {
				crc <<= 1
			}
		}
	}
	return crc
}

func (d *data) parseMessage(mess []uint8) {
	d.ctcResult = "Некорректно"
	if ((d.functionCode == 0x01 || d.functionCode == 0x03) && len(mess) != 4) ||
		((d.functionCode == 0x05 || d.functionCode == 0x06) && len(mess) != 5) {
		d.result = [5]string{
			"Ошибка: Некорректный формат данных",
		}
	} else {

		switch d.functionCode {
		case 0x01:
			if crc8(mess[:3]) == d.crc {
				d.ctcResult = "Корректно"
			}
			d.result = [5]string{
				"Адрес устройства: " + fmt.Sprintf("%d", d.deviceAddress),
				"Код функции: " + fmt.Sprintf("%d (%s)", d.functionCode, "Чтение бит"),
				"Адрес ячейки памяти, откуда идёт чтение: " + fmt.Sprintf("%d", d.memoryAddress),
				"Результат проверки контрольной суммы: " + d.ctcResult,
			}
		case 0x03:
			if crc8(mess[:3]) == d.crc {
				d.ctcResult = "Корректно"
			}
			d.result = [5]string{
				"Адрес устройства: " + fmt.Sprintf("%d", d.deviceAddress),
				"Код функции: " + fmt.Sprintf("%d (%s)", d.functionCode, "Чтение байт"),
				"Адрес ячейки памяти, откуда идёт чтение: " + fmt.Sprintf("%d", d.memoryAddress),
				"Результат проверки контрольной суммы: " + d.ctcResult,
			}
		case 0x05:
			if crc8(mess[:4]) == d.crc {
				d.ctcResult = "Корректно"
			}
			d.result = [5]string{
				"Адрес устройства: " + fmt.Sprintf("%d", d.deviceAddress),
				"Код функции: " + fmt.Sprintf("%d (%s)", d.functionCode, "Запись бит"),
				"Адрес ячейки памяти, куда идёт запись: " + fmt.Sprintf("%d", d.memoryAddress),
				"Записываемое значение: " + fmt.Sprintf("%d", mess[3]),
				"Результат проверки контрольной суммы: " + d.ctcResult,
			}
		case 0x06:
			if crc8(mess[:4]) == d.crc {
				d.ctcResult = "Корректно"
			}
			d.result = [5]string{
				"Адрес устройства: " + fmt.Sprintf("%d", d.deviceAddress),
				"Код функции: " + fmt.Sprintf("%d (%s)", d.functionCode, "Запись байт"),
				"Адрес ячейки памяти, куда идёт запись: " + fmt.Sprintf("%d", d.memoryAddress),
				"Записываемое значение: " + fmt.Sprintf("%d", mess[3]),
				"Результат проверки контрольной суммы: " + d.ctcResult,
			}
		default:
			d.result = [5]string{
				"Ошибка: Неизвестная операция",
			}
		}
	}
}

func main() {
	var d data
	messages := [][]uint8{
		{0x01, 0x03, 0x01, 0x53},
		{0x3A, 0x01, 0x10, 0x03},
		{0xAB, 0x06, 0x4C, 0x13, 0xA8},
		{0x0C, 0x05, 0x04, 0x00, 0x7C},
		{0x01, 0x01, 0x01, 0x01, 0x79},
		{0x01, 0x06, 0x01, 0x01},
		{0x01},
		{0x01, 0x06},
		{0x01, 0x06, 0x01}}

	for _, v := range messages {
		if len(v) < 4 || len(v) > 5 {
			fmt.Println("Ошибка: Некорректный формат данных")
		} else {
			d = data{
				deviceAddress: v[0],
				functionCode:  v[1],
				memoryAddress: v[2],
				crc:           v[len(v)-1],
			}

			d.parseMessage(v)

			for _, value := range d.result {
				if value != "" {
					fmt.Printf("%s\n", value)
				}
			}
			fmt.Println("")
		}

	}
}
