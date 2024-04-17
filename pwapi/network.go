package pwapi

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"reflect"
	"strconv"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type PwCuint int32
type Cuint PwCuint

func (c *Cuint) UnmarshalBinary(reader *bytes.Reader) (int, error) {
	// Read the first byte to determine the size
	var b byte
	b, err := reader.ReadByte()
	if err != nil {
		return 0, err // Handle error gracefully
	}

	//volta o reader para o inicio
	reader.Seek(-1, io.SeekCurrent)

	size := 1
	var min int64
	min = 0
	if b >= 0x80 {
		if b < 0xC0 {
			size = 2
			min = 0x8000
		} else if b < 0xE0 {
			size = 4
			min = 0xC0000000

		} else {
			size = 5
		}
	}

	// Ensure sufficient bytes are available
	if reader.Len() < size {
		return 0, errors.New("invalid cuint data: insufficient bytes")
	}

	// Read the remaining bytes of the Cuint value
	data := make([]byte, size)

	_, err = reader.Read(data) // Read the remaining bytes
	if err != nil {
		if err == io.EOF {
			return 0, err // Handle error gracefully
		}
		fmt.Printf("err1: %v\n", err)
		return 0, err // Handle error gracefully
	}
	buf := bytes.NewReader(data)

	var value int32
	switch size {
	case 1:
		var tmp int8
		err = binary.Read(buf, binary.BigEndian, &tmp)
		if err != nil {
			fmt.Printf("err1: %v\n", err)
		}

		value = int32(tmp)
	case 2:
		var tmp uint16
		err = binary.Read(buf, binary.BigEndian, &tmp)
		if err != nil {
			fmt.Printf("err2: %v\n", err)
		}
		value = int32(tmp)

	case 4:
		err = binary.Read(buf, binary.BigEndian, &value)
		if err != nil {
			fmt.Printf("err4: %v\n", err)
		}
	case 5:
		fmt.Println("Tamanho de 5 bytes não suportado neste exemplo.")
		return 0, nil
	default:
		fmt.Println("Tamanho inválido.")
		return 0, nil
	}

	// Assign the integer value to the Cuint
	*c = Cuint(value - int32(min))
	return size, nil
}

func createPack(args ...interface{}) []byte {
	var packBuffer bytes.Buffer

	for _, arg := range args {
		v := reflect.ValueOf(arg)
		if v.Kind() != reflect.Struct {
			fmt.Println("Erro: fornecido um valor inválido para a função. Deve ser um ponteiro não nulo.")
			os.Exit(1)
		}
		// Iterate through the fields of the struct
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			switch field.Kind() {
			case reflect.Uint8:
				binary.Write(&packBuffer, binary.BigEndian, field.Interface().(uint8))
			case reflect.Int:
				binary.Write(&packBuffer, binary.BigEndian, int32(field.Interface().(int)))
			case reflect.String:
				packString(&packBuffer, field.Interface().(string))
			case reflect.Slice:
				size := field.Len()

				var sizebytes []byte

				if size < 64 {
					sizebytes = []byte{byte(size)}
				} else if size < 16384 {
					sizebytes = []byte{byte((size >> 8) | 0x80), byte(size)}
				} else if size < 536870912 {
					sizebytes = []byte{byte((size >> 24) | 0xC0), byte((size >> 16) & 0xFF), byte((size >> 8) & 0xFF), byte(size & 0xFF)}
				} else {
					sizebytes = []byte{byte(0xE0), byte((size >> 24) & 0xFF), byte((size >> 16) & 0xFF), byte((size >> 8) & 0xFF), byte(size & 0xFF)}
				}

				binary.Write(&packBuffer, binary.BigEndian, sizebytes)
				binary.Write(&packBuffer, binary.BigEndian, field.Interface())
			case reflect.Struct:
				var bytes []byte
				bytes = createPack(field.Interface())
				binary.Write(&packBuffer, binary.BigEndian, bytes)
			default:
				fmt.Printf("Tipo de campo inválido: %v\n", field.Kind())
				os.Exit(1)

			}
		}
	}

	return packBuffer.Bytes()
}

func packString(buf *bytes.Buffer, s string) {

	utf16le := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
	transformed, _, _ := transform.String(utf16le, s)
	//transforma em bytes
	bytes := []byte(transformed)
	size := len(bytes)
	var sizebytes []byte

	if size < 64 {
		sizebytes = []byte{byte(size)}
	} else if size < 16384 {
		sizebytes = []byte{byte((size >> 8) | 0x80), byte(size)}

	} else if size < 536870912 {
		sizebytes = []byte{byte((size >> 24) | 0xC0), byte((size >> 16) & 0xFF), byte((size >> 8) & 0xFF), byte(size & 0xFF)}
	} else {
		sizebytes = []byte{byte(0xE0), byte((size >> 24) & 0xFF), byte((size >> 16) & 0xFF), byte((size >> 8) & 0xFF), byte(size & 0xFF)}
	}
	buf.Write(sizebytes)
	buf.Write(bytes)
}

func unpackData(data []byte, structure interface{}) []byte {
	buf := bytes.NewReader(data)
	destValue := reflect.ValueOf(structure)
	if destValue.Kind() != reflect.Ptr || destValue.IsNil() {
		fmt.Println("Erro: fornecido um valor inválido para a função. Deve ser um ponteiro não nulo.")
		return nil
	}

	destValue = destValue.Elem()

	var cycleCount int
	cycleCount = -1
	for i := 0; i < destValue.NumField(); i++ {
		field := destValue.Field(i)
		fieldType := field.Type()

		if cycleCount > -1 {
			if cycleCount > 0 {
				for j := 0; j < cycleCount; j++ {
					//posicao do ponteiro
					item := reflect.New(fieldType.Elem()).Elem()
					outradata := unpackData(data, item.Addr().Interface())
					data = outradata
					field.Set(reflect.Append(field, item))
				}

				cycleCount = -1

			} else {

				//posicao do ponteiro
				item := reflect.New(fieldType.Elem()).Elem()
				outradata := unpackData(data, item.Addr().Interface())
				data = outradata
				field.Set(reflect.Append(field, item))

			}
		} else {
			switch fieldType {
			case reflect.TypeOf(int(0)):
				value := int(binary.BigEndian.Uint32(data[:4]))
				field.SetInt(int64(value))
				data = data[4:]
			case reflect.TypeOf(UserID(0)):
				value := int(binary.BigEndian.Uint32(data[:4]))
				field.SetInt(int64(value))
				data = data[4:]
			case reflect.TypeOf(int64(0)):
				var value int64
				error := binary.Read(buf, binary.BigEndian, &value)
				if error != nil {
					fmt.Printf("Erro ao ler: %v\n", error)
				}

				field.SetInt(value)
				data = data[int(binary.Size(value)):]
			case reflect.TypeOf(byte(0)):
				var value byte
				value = data[0]
				data = data[1:]

				field.SetUint(uint64(value))

			case reflect.TypeOf(Cuint(0)):
				buf = bytes.NewReader(data)
				var cui Cuint
				size, error := cui.UnmarshalBinary(buf)
				if error != nil {
					fmt.Printf("Erro ao ler: %v\n", error)
					os.Exit(1)
				}
				// Save the position of the pointer after reading the value
				field.Set(reflect.ValueOf(cui))

				data = data[size:]
				if cui > 0 {
					cycleCount = int(cui)
				} else {
					cycleCount = -1
				}
			case reflect.TypeOf(float32(0)):
				var value float32
				error := binary.Read(buf, binary.BigEndian, &value)
				if error != nil {
					fmt.Printf("Erro ao ler: %v\n", error)
				}
				field.SetFloat(float64(value))

				data = data[4:]
			case reflect.TypeOf([]byte{}):
				buf = bytes.NewReader(data)
				var length Cuint
				size, error := length.UnmarshalBinary(buf)
				if error != nil {
					fmt.Printf("Erro ao ler: %v\n", error)
					os.Exit(1)
				}

				// Read the remaining bytes of the string
				value := make([]byte, length)

				err := binary.Read(buf, binary.BigEndian, &value)
				if err != nil {
					fmt.Printf("Erro ao ler: %v\n", err)

				}

				field.Set(reflect.ValueOf(value))
				data = data[size+int(length):]
			case reflect.TypeOf(uint16(0)):
				var value uint16
				error := binary.Read(buf, binary.BigEndian, &value)
				if error != nil {
					fmt.Printf("Erro ao ler: %v\n", error)
				}
				field.SetUint(uint64(value))
				data = data[2:]
			case reflect.TypeOf(uint(0)):
				var value uint32
				error := binary.Read(buf, binary.BigEndian, &value)
				if error != nil {
					fmt.Printf("Erro ao ler: %v\n", error)
				}
				field.SetUint(uint64(value))
				data = data[4:]

			case reflect.TypeOf(string("")):
				size := 1
				if data[0] >= 0x80 {
					if data[0] < 0xC0 {
						size = 2
					} else if data[0] < 0xE0 {
						size = 4
					} else {
						size = 5
					}
				}
				var octetlen int
				octetlen = 0
				if size == 1 {
					octetlen = int(data[0])
				} else {
					octetlen = int(data[0]) - 0x80 + int(data[1])
				}

				//discart size
				data = data[size:]
				// Read the remaining bytes of the string
				str := string(data[:octetlen])
				//change encoding
				utf16le := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
				transformed, _, _ := transform.String(utf16le, str)
				str = transformed
				// Save the position of the pointer after reading the value

				field.SetString(str)
				data = data[octetlen:]

			}
		}
	}

	return data
}

// cuint converte um número inteiro não negativo em um slice de bytes
// no formato compact uint.
func cuint(data uint32) []byte {
	if data < 64 {
		return []byte{byte(data)}
	} else if data < 16384 {
		return []byte{byte((data >> 8) | 0x80), byte(data)}
	} else if data < 536870912 {
		return []byte{byte((data >> 24) | 0xC0), byte((data >> 16) & 0xFF), byte((data >> 8) & 0xFF), byte(data & 0xFF)}
	}

	return []byte{byte(0xE0), byte((data >> 24) & 0xFF), byte((data >> 16) & 0xFF), byte((data >> 8) & 0xFF), byte(data & 0xFF)}
}

func createHeader(opcodeHex string, data []byte) []byte {
	// Converte a string hexadecimal para um inteiro
	opcode, err := strconv.ParseInt(opcodeHex, 16, 32)
	if err != nil {
		fmt.Printf("Erro ao converter o opcode: %v\n", err)
		os.Exit(1)
	}
	// Empacota o opcode como um inteiro de 32 bits
	opcodeBytes := cuint(uint32(opcode))

	// Empacota o comprimento do pacote como um inteiro de 16 bits
	lengthBytes := cuint(uint32(len(data)))

	// Cria um buffer para armazenar os bytes do cabeçalho
	var headerBuffer bytes.Buffer

	// Escreve os bytes do opcode
	headerBuffer.Write(opcodeBytes)
	// Escreve os bytes do comprimento do pacote
	headerBuffer.Write(lengthBytes)
	// Escreve os bytes do pacote
	headerBuffer.Write(data)

	// Retorna os bytes do cabeçalho empacotado
	return headerBuffer.Bytes()
}

func SendToDelivery(data []byte, recvAfterSend bool, justSend bool) ([]byte, error) {

	port := AppConfig.Ports["gdeliveryd"]
	return SendToSocket(data, port, recvAfterSend, nil, justSend)
}

func SendToProvider(data []byte, recvAfterSend bool, justSend bool) ([]byte, error) {

	port := AppConfig.Ports["provider"]
	buf := make([]byte, 8196)

	return SendToSocket(data, port, recvAfterSend, buf, justSend)
}
func SendToGamedBD(data []byte, recvAfterSend bool, justSend bool) ([]byte, error) {
	port := AppConfig.Ports["gamedbd"]
	buf := make([]byte, 8196)

	return SendToSocket(data, port, recvAfterSend, buf, justSend)
}

func SendToSocket(data []byte, port int, recvAfterSend bool, buf []byte, justSend bool) ([]byte, error) {

	conn, err := net.Dial("tcp", net.JoinHostPort(AppConfig.IP, strconv.Itoa(port)))
	if err != nil {
		fmt.Printf("Erro ao conectar ao socket: %v\n", err)
		return nil, err
	}
	defer conn.Close()

	_, err = conn.Write(data)
	if err != nil {
		fmt.Printf("Erro ao enviar para o socket: %v\n", err)
		return nil, err
	}

	if justSend {
		return nil, nil
	}
	var tmp []byte

	// Leia e descarte um número fixo de bytes após o envio
	var readLen int
	readLen = 23
	tmp = make([]byte, readLen)
	_, err23 := io.ReadFull(conn, tmp)
	if err23 != nil {
		return nil, err23
	}

	// Reproduz o comportamento de leitura da função PHP, atendendo aos diferentes casos
	switch 3 {
	case 1:
		// Leia um número fixo de bytes
		tmp := make([]byte, 8192)

		_, err = io.ReadFull(conn, tmp)
		if err != nil {
			return nil, err
		}
	case 2:
		// Leia até encontrar um buffer com menos de 1024 bytes
		var buffer []byte
		for {
			buffer = make([]byte, 1024)
			n, err := conn.Read(buffer)
			if err != nil {
				return nil, err
			}
			buf = append(buf, buffer[:n]...)
			if n < 1024 {
				break
			}
		}
	case 3:
		// Leia o cabeçalho de 8 bytes contendo o tamanho da mensagem
		var buf []byte
		if !recvAfterSend {
			buf = tmp
		} else {
			buf = make([]byte, 16)
			_, err = conn.Read(buf)
			if err != nil {
				return nil, err
			}
			readLen = 16
		}

		bufReader := bytes.NewReader(buf)

		var cui Cuint
		var size int
		size, error := cui.UnmarshalBinary(bufReader)
		if error != nil {
			fmt.Printf("Erro ao ler: %v\n", error)
		}

		var length Cuint
		var size2 int
		size2, error = length.UnmarshalBinary(bufReader)

		if error != nil {
			fmt.Printf("Erro ao ler: %v\n", error)
		}
		restante := int64(length) - int64(readLen) + int64(size) + int64(size2)
		if restante > 0 {
			restobuf := make([]byte, restante)
			_, err := conn.Read(restobuf)
			if err != nil {
				return nil, err
			}

			//juntar os buffers
			buf = append(buf, restobuf...)
		}

		return buf, nil
	}

	return buf, nil
}

func deleteHeader(data []byte) []byte {

	length := 8
	var cuint1 Cuint
	var size int
	size, error := cuint1.UnmarshalBinary(bytes.NewReader(data))
	if error != nil {
		fmt.Printf("Erro ao ler: %v\n", error)
	}

	//remove size bytes
	data = data[size:]

	var cuint2 Cuint
	size, error = cuint2.UnmarshalBinary(bytes.NewReader(data))
	if error != nil {
		fmt.Printf("Erro ao ler: %v\n", error)
	}

	//remove size bytes
	data = data[size:]

	data = data[length:]
	return data
}

func ConvertToBytes(inputString string) ([]byte, error) {
	// Convertendo a string para bytes (big-endian)
	var bytesRepresentation []byte

	// Iterando sobre a string para converter cada número
	for i := 0; i < len(inputString); i += 8 {
		numStr := inputString[i : i+8]
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return nil, err
		}

		// Convertendo cada número decimal para uma representação de bytes (big-endian)
		numBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(numBytes, uint64(num))

		// Adicionando os bytes à representação total
		bytesRepresentation = append(bytesRepresentation, numBytes...)
	}

	return bytesRepresentation, nil
}
