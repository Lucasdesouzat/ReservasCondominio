package constants

import (
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
