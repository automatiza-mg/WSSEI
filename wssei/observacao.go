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

// CadastrarObservacaoParams reúne os parâmetros necessários
// para cadastrar uma observação no WSSEI.
type CadastrarObservacaoParams struct {
	Descricao string `json:"descricao"`
	Unidade   int    `json:"unidade"`
	Protocolo int    `json:"protocolo"`
}

// CadastrarObservacao cadastra uma observação no WSSEI.
func (c *Client) CadastrarObservacao(ctx context.Context, params CadastrarObservacaoParams) error {
	if strings.TrimSpace(params.Descricao) == "" {
		return fmt.Errorf("descricao required: %d", params.Descricao)
	}
	if params.Unidade <= 0 {
		return fmt.Errorf("unidade invalida: %d", params.Unidade)
	}
	if params.Protocolo <= 0 {
		return fmt.Errorf("protocolo invalido: %d", params.Protocolo)
	}

	jsonBody, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("json body: %w", err)
	}

	endpoint := c.endpoint + "/observacao/"

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
