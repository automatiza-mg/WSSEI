package wssei

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// AnotacaoParams reúne os dados necessários para cadastrar uma anotação em um protocolo.
type AnotacaoParams struct {
	Descricao  string         `json:"descricao"`
	Protocolo  int            `json:"protocolo"`
	Unidade    int            `json:"unidade"`
	Usuario    int            `json:"usuario"`
	Prioridade TipoPrioridade `json:"prioridade"`
}

// TipoPrioridade representa os valores aceitos para definir a prioridade da anotação.
type TipoPrioridade string

const (
	Sim TipoPrioridade = "S"
	Nao TipoPrioridade = "N"
)

// CadastrarAnotacao cria uma anotação vinculada a um protocolo.
func (c *Client) CadastrarAnotacao(ctx context.Context, params AnotacaoParams) error {
	if strings.TrimSpace(params.Descricao) == "" {
		return fmt.Errorf("descricao required: %s", params.Descricao)
	}
	if params.Protocolo <= 0 {
		return fmt.Errorf("protocolo inválido: %d", params.Protocolo)
	}

	jsonBody, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("json body: %w", err)
	}

	endpoint := c.endpoint + "/anotacao/"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("http do: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d: %s", res.StatusCode, strings.TrimSpace(string(body)))
	}

	var env Envelope[struct{}]
	if err := json.Unmarshal(body, &env); err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}

	if !env.Sucesso {
		return fmt.Errorf("invalid response: %s", env.Mensagem)
	}

	return nil
}
