package wssei

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// AcompanhamentoEspecial representa um acompanhamento especial retornado pelo
// WSSEI.
type AcompanhamentoEspecial struct {
	IDAcompanhamento        string                          `json:"idAcompanhamento"`
	IDGrupoAcompanhamento   string                          `json:"idGrupoAcompanhamento"`
	NomeGrupoAcompanhamento string                          `json:"nomeGrupoAcompanhamento"`
	IDProtocolo             string                          `json:"idProtocolo"`
	IDUsuarioGerador        string                          `json:"idUsuarioGerador"`
	DataGeracao             string                          `json:"dataGeracao"`
	Observacao              string                          `json:"observacao"`
	SiglaUsuario            string                          `json:"siglaUsuario"`
	NomeUsuario             string                          `json:"nomeUsuario"`
	TipoVisualizacao        string                          `json:"tipoVisualizacao"`
	NomeTipoProcedimento    string                          `json:"nomeTipoProcedimento"`
	ProtocoloFormatado      string                          `json:"protocoloFormatado"`
	Atributos               AcompanhamentoEspecialAtributos `json:"atributos"`
}

// AcompanhamentoEspecialAtributos reúne os atributos de um
// [AcompanhamentoEspecial].
type AcompanhamentoEspecialAtributos struct {
	Anotacao                   []string `json:"anotacao"`
	ProcessoBloqueado          bool     `json:"processoBloquado"`
	RemocaoSobrestamento       bool     `json:"remocaoSobrestamento"`
	DocumentoAssinadoProcesso  bool     `json:"documentoAssinadoProcesso"`
	DocumentoPublicadoProcesso bool     `json:"documentoPublicadoProcesso"`
	RetornoProgramado          []string `json:"retornoProgramado"`
	AndamentoSituacao          []string `json:"andamentoSituacao"`
	AndamentoMarcador          []string `json:"andamentoMarcador"`
}

// AcompanhamentoEspecialParams reúne os parâmetros opcionais da listagem de
// acompanhamentos especiais.
//
// Campos com valor zero (0 ou "") são omitidos da requisição.
type AcompanhamentoEspecialParams struct {
	// Limit é o limite de registros da paginação.
	Limit int
	// Start é a página de início da paginação.
	Start int
	// GrupoAcompanhamento é o identificador do grupo de acompanhamento para filtro da pesquisa.
	GrupoAcompanhamento int
}

// values converte os parâmetros da pesquisa em query params,
// omitindo campos que possuem valor zero.
func (p AcompanhamentoEspecialParams) values() url.Values {
	q := make(url.Values)
	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Start != 0 {
		q.Set("start", strconv.Itoa(p.Start))
	}
	if p.GrupoAcompanhamento != 0 {
		q.Set("grupoAcompanhamento", strconv.Itoa(p.GrupoAcompanhamento))
	}
	return q
}

// ListarAcompanhamentoEspecial retorna os acompanhamentos especiais
// cadastrados.
func (c *Client) ListarAcompanhamentoEspecial(ctx context.Context, params AcompanhamentoEspecialParams) ([]AcompanhamentoEspecial, int, error) {
	endpoint := c.endpoint + "/acompanhamentoespecial/listar"
	if q := params.values().Encode(); q != "" {
		endpoint += "?" + q
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("http do: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("read body: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("unexpected status %d: %s", res.StatusCode, strings.TrimSpace(string(body)))
	}

	var env Envelope[[]AcompanhamentoEspecial]
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, 0, fmt.Errorf("json unmarshal: %w", err)
	}

	if !env.Sucesso {
		return nil, 0, fmt.Errorf("invalid response: %s", env.Mensagem)
	}

	total, err := env.getTotal()
	if err != nil {
		return nil, 0, fmt.Errorf("parse total %q: %w", env.Total, err)
	}

	return env.Data, total, nil
}
