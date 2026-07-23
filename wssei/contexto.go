package wssei

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Contexto representa os dados retornados na listagem de contextos de um órgão.
type Contexto struct {
	ID         string `json:"id"`
	Nome       string `json:"nome"`
	Descricao  string `json:"descricao"`
	BaseDNLDAP string `json:"base_dn_ldap"`
}

// ListarContextoOrgao retorna a lista de contextos vinculados ao órgão informado.
func (c *Client) ListarContextoOrgao(ctx context.Context, orgao int) ([]Contexto, error) {
	if orgao <= 0 {
		return nil, fmt.Errorf("orgao invalido: %d", orgao)
	}
	endpoint := fmt.Sprintf("%s/contexto/listar/%d", c.endpoint, orgao)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http do: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", res.StatusCode, strings.TrimSpace(string(body)))
	}

	var env Envelope[[]Contexto]
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	if !env.Sucesso {
		return nil, fmt.Errorf("invalid response: %s", env.Mensagem)
	}

	return env.Data, nil
}
