# WSSEI

Client Go para o módulo [WSSEI](https://pengovbr.github.io/mod-wssei/#/) do
SEI (Sistema Eletrônico de Informações), usado para automatizar processos,
documentos, blocos de assinatura e demais operações do SEI. Projetado para
ser compartilhado entre múltiplos projetos.

## Instalação

```bash
go get github.com/automatiza-mg/WSSEI
```

Requer Go 1.25 ou superior.

## Uso

O client autentica com usuário e senha uma única vez e reaproveita o token
gerado em cada requisição, renovando-o automaticamente quando o servidor
responde 401/403.

```go
package main

import (
    "context"
    "fmt"
    "os"
    "strconv"

    wssei "github.com/automatiza-mg/WSSEI"
)

func main() {
    orgao, _ := strconv.Atoi(os.Getenv("SEI_ORGAO"))

    client := wssei.NewClient(wssei.Config{
        BaseURL: os.Getenv("SEI_BASE_URL"),
        Usuario: os.Getenv("SEI_USUARIO"),
        Senha:   os.Getenv("SEI_SENHA"),
        Orgao:   orgao,
    })

    ctx := context.Background()

    processo, err := client.ConsultarProcesso(ctx, 1234567)
    if err != nil {
        panic(err)
    }

    fmt.Println(processo.ProtocoloFormatado)
}
```

### Consultar e listar processos

```go
processos, total, err := client.ListarProcessos(ctx, wssei.ListarProcessosParams{
    Limit: 10,
    Start: 0,
})
```

### Documentos

```go
doc, err := client.ConsultarDocumentoInterno(ctx, protocolo)

// Baixa o conteúdo de um documento externo (anexo).
body, contentType, err := client.BaixarAnexo(ctx, protocolo)
defer body.Close()
```

### Blocos de assinatura

```go
blocos, total, err := client.PesquisarBlocoAssinatura(ctx, wssei.PesquisarBlocoAssinaturaParams{
    Estado: wssei.EstadoSituacaoDisponibilizado,
    Limit:  20,
})

err = client.AssinarBlocoAssinatura(ctx, bloco, wssei.AssinarBlocoAssinaturaParams{
    Orgao:   orgao,
    Cargo:   "Servidor (a) Público (a)",
    Login:   usuario,
    Senha:   senha,
    Usuario: idUsuario,
})
```

### Marcadores

```go
marcador, err := client.ConsultarMarcador(ctx, protocolo)

err = client.MarcarProcesso(ctx, protocolo, wssei.MarcadorProcessoParams{
    Texto:    "Aguardando análise",
    Marcador: idMarcador,
})
```

## Configuração

### Config

O client é configurado com uma struct [`Config`](client.go):

| Campo             | Descrição                                                                   |
| ----------------- | --------------------------------------------------------------------------- |
| `BaseURL`         | URL base do SEI (ex: `https://www.sei.mg.gov.br`).                          |
| `Usuario`         | Login usado na autenticação.                                                |
| `Senha`           | Senha usada na autenticação.                                                |
| `Orgao`           | Id do órgão da autenticação.                                                |
| `Plataforma`      | Identificador da plataforma dona das credenciais (ex: `whatsapp`).          |
| `PlataformaID`    | Identificador do usuário dentro da plataforma (ex: número do WhatsApp).     |
| `OnAuthenticated` | Callback opcional, invocado após cada autenticação bem-sucedida (ver abaixo). |

### AuthCallback

`OnAuthenticated` recebe os dados retornados pelo WSSEI (`AuthResponse`) e
permite persistir ou observar o token, o `IdUsuario`, as unidades e os
perfis do usuário autenticado. É chamado no login inicial e a cada renovação
automática do token.

```go
cfg.OnAuthenticated = func(ctx context.Context, plataforma, plataformaID string, resp *wssei.AuthResponse) error {
    return cache.SaveToken(ctx, plataformaID, resp.Token)
}
```

### Autenticação isolada

Se você precisa apenas autenticar (sem executar chamadas subsequentes), use
[`Auth`](auth.go) diretamente:

```go
auth := wssei.NewAuth(baseURL)
resp, err := auth.Autenticar(ctx, usuario, senha, orgao)
```

## Envelope

Todas as respostas do WSSEI seguem o formato:

```json
{
    "sucesso": true,
    "mensagem": "",
    "total": "42",
    "data": { ... }
}
```

Os métodos do [`Client`](client.go) já extraem `data` e `total` para o
chamador, transformando `sucesso: false` em erro. O tipo genérico
[`Envelope[T]`](client.go) é exportado para uso em cenários onde a resposta
bruta é necessária.

## Compatibilidade Latin-1

Alguns endpoints do WSSEI (PHP legado) interpretam campos específicos como
Latin-1 mesmo recebendo JSON UTF-8 — o exemplo mais conhecido é o `cargo` em
assinaturas. O client transcodifica automaticamente esses campos usando
`jsonStringLatin1`, garantindo que acentos sejam aceitos pelo servidor.

## Referência

- [Documentação oficial do módulo WSSEI](https://pengovbr.github.io/mod-wssei/#/)

## Licença

[MIT](./LICENSE)
