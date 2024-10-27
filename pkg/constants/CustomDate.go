package constants

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// Tipo CustomDate para manipulação de datas no formato yyyy-mm-dd
type CustomDate struct {
	time.Time
}

// Implementação do UnmarshalJSON para o tipo CustomDate
func (cd *CustomDate) UnmarshalJSON(data []byte) error {
	// Definindo o formato esperado
	const layout = "2006-01-02"
	str := string(data)
	str = str[1 : len(str)-1] // Remove aspas do JSON

	parsedTime, err := time.Parse(layout, str)
	if err != nil {
		return fmt.Errorf("data inválida: %v", err)
	}

	cd.Time = parsedTime
	return nil
}

// Implementação opcional do MarshalJSON, caso você precise serializar a data novamente
func (cd CustomDate) MarshalJSON() ([]byte, error) {
	const layout = "2006-01-02"
	return json.Marshal(cd.Time.Format(layout))
}

// Método para retornar o time.Time para inserção no banco
func (cd CustomDate) ToTime() time.Time {
	return cd.Time
}

// Implementação da interface Scanner do pacote database/sql
// Permite converter valores do banco de dados para CustomDate
func (cd *CustomDate) Scan(value interface{}) error {
	if value == nil {
		*cd = CustomDate{Time: time.Time{}}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*cd = CustomDate{Time: v}
	case string:
		// Para o caso de o valor vir como string
		const layout = "2006-01-02"
		parsedTime, err := time.Parse(layout, v)
		if err != nil {
			return fmt.Errorf("não foi possível converter string para CustomDate: %v", err)
		}
		*cd = CustomDate{Time: parsedTime}
	default:
		return fmt.Errorf("não foi possível converter %T para CustomDate", value)
	}
	return nil
}

// Implementação da interface Valuer do pacote database/sql
// Permite converter CustomDate para valores suportados pelo banco de dados
func (cd CustomDate) Value() (driver.Value, error) {
	// Retorna o valor no formato que o banco de dados espera (time.Time)
	return cd.Time, nil
}
