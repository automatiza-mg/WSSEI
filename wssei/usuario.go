package wssei

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// ListarUsuariosParams seleciona a query do ListasUsuarios
type ListarUsuariosParams struct {
	// Limit é o limite de registros da paginação.
	Limit int
	// Start é a página de início da paginação.
	Start int
	//Procedimento é o ID do processo. OBRIGATORIO
	Unidade int
}

// Converte os parâmetros em [url.Values], omitindo os campos zerados.
func (p ListarUsuariosParams) values() url.Values {
	q := make(url.Values)
	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Start != 0 {
		q.Set("start", strconv.Itoa(p.Start))
	}
	if p.Unidade != 0 {
		q.Set("unidade", strconv.Itoa(p.Unidade))
	}
	return q
}

// ListarUsuarios retorna a lista de Usuários
func (c *Client) ListarUsuarios(ctx context.Context, params ListarUsuariosParams) ([]Usuarios, int, error) {
	url := fmt.Sprintf("%s/usuario/listar", c.endpoint)
	if q := params.values().Encode(); q != "" {
		url += "?" + q
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("erro request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("erro response: %w", err)
	}
	defer resp.Body.Close()

	var result Envelope[[]Usuarios]

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, 0, fmt.Errorf("erro json decoder: %w", err)
	}

	if result.Sucesso != true {
		return nil, 0, fmt.Errorf("erro listar ususarios : %s", result.Mensagem)
	}

	total, err := result.getTotal()
	if err != nil {
		return nil, 0, fmt.Errorf("total invalido")
	}

	return result.Data, total, nil

}

// Usuarios tipo utilizado na funcao "ListarUsuarios"
type Usuarios struct {
	IDUsuario string `json:"id_usuario"`
	Sigla     string `json:"sigla"`
	Nome      string `json:"nome"`
	IDContato string `json:"id_contato"`
	Total     string `json:"total"`
}

// ListarUsuariosParams seleciona a query do ListasUsuarios
type PesquisarUsuariosParams struct {
	// Limit é o limite de registros da paginação.
	Limit int
	// Start é a página de início da paginação.
	Start int
	//Procedimento é o ID do processo. OBRIGATORIO
	Unidade int
}

// Converte os parâmetros em [url.Values], omitindo os campos zerados.
func (p PesquisarUsuariosParams) values() url.Values {
	q := make(url.Values)
	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Start != 0 {
		q.Set("start", strconv.Itoa(p.Start))
	}
	if p.Unidade != 0 {
		q.Set("unidade", strconv.Itoa(p.Unidade))
	}
	return q
}

// PesquiarUsuarios retorna a pesquisa de Usuários
func (c *Client) PesquiarUsuarios(
	ctx context.Context,
	limit int,
	start int,
	unidade int,
) (*Usuarios, int, error) {
	url := fmt.Sprintf(
		"%s/usuario/listar",
		c.endpoint,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("erro request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("erro response: %w", err)
	}
	defer resp.Body.Close()

	var result Envelope[Usuarios]

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, 0, fmt.Errorf("erro json decoder: %w", err)
	}

	if result.Sucesso != true {
		return nil, 0, fmt.Errorf("erro listar ususarios : %s", result.Mensagem)
	}

	total, err := result.getTotal()
	if err != nil {
		return nil, 0, fmt.Errorf("total invalido")
	}

	return &result.Data, total, nil

}

// UsuariosPesquisa tipo utilizado na funcao "PesquiarUsuarios"
type UsuariosPesquisa struct {
	IDContato string `json:"id_contato"`
	IDUsuario string `json:"id_usuario"`
	Sigla     string `json:"sigla"`
	Nome      string `json:"nome"`
}
