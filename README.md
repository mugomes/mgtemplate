# MGTemplate

MGTemplate é um mecanismo de templates leve escrito em Go, focado em **simplicidade**, **performance** e **controle explícito por código**.

Ele permite separar HTML da lógica da aplicação sem depender de parsing complexo, ASTs ou fases de compilação. O template é interpretado de forma incremental conforme o código executa.

---

## ✨ Características

- Interpolação simples de variáveis: `{{title}}`
- Suporte a propriedades de structs: `{{user.name}}`
- Modificadores encadeáveis: `{{title|upper|trim}}`
- Blocos reutilizáveis com repetição controlada por código
- Zero reflexão dinâmica de funções
- Ideal para servidores HTTP, ferramentas CLI e projetos embarcados

---

## 📦 Instalação

```bash
go get github.com/mugomes/mgtemplate
````

---

## 🚀 Uso básico

### Template (`template.html`)

```html
<!DOCTYPE html>
<html>
<head>
    <title>{{title}}</title>
</head>
<body>

<h1>{{title|upper}}</h1>

[[ITEM]]
<div>
    <strong>{{user.name}}</strong><br>
    <small>{{info}}</small>
</div>
[[/ITEM]]

</body>
</html>
```

---

### Código Go

```go
package main

import (
	"net/http"
	"strconv"

	"github.com/mugomes/mgtemplate"
)

type User struct {
	Name string
}

func handler(w http.ResponseWriter, r *http.Request) {
	tpl, err := mgtemplate.ReadFile("template.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	tpl.Var("title", "MGTemplate Example")

	for i := 0; i < 5; i++ {
		tpl.Var("user", User{
			Name: "Usuário " + strconv.Itoa(i),
		})

		tpl.Var("info", "Informação " + strconv.Itoa(i))
		tpl.Section("ITEM")
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(tpl.Render()))
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8000", nil)
}
```

---

## 🧩 Blocos (Sections)

Blocos são definidos diretamente no HTML:

```html
[[ITEM]]
<p>{{user.name}}</p>
[[/ITEM]]
```

E ativados no código sempre que necessário:

```go
tpl.Section("ITEM")
```

Cada chamada adiciona uma nova instância do bloco renderizado.

---

## 🔧 Variáveis

### Variável simples

```html
{{title}}
```

```go
tpl.Var("title", "Exemplo")
```
Também é possível verificar se a variável existe no código HTML.

```go
if tpl.VarExist("user") {
	tpl.Var("title", "Exemplo")
}
```

---

### Structs

```html
{{user.name}}
```

```go
tpl.Var("user", User{Name: "Usuário"})
```

O acesso a campos é **case-insensitive** e ignora `_`.

---

### Incluir Outros Arquivos HTML

```html
{{menutopo}}
```

```go
tpl.IncludeFile("menutopo", "public/menutopo.html")
```

## 🔁 Modificadores

Modificadores podem ser encadeados usando `|`.

| Modificador | Descrição                |
| ----------- | ------------------------ |
| `upper`     | Converte para maiúsculas |
| `lower`     | Converte para minúsculas |
| `trim`      | Remove espaços           |

Exemplo:

```html
{{title|upper|trim}}
```

---

## 👤 Autor

**Murilo Gomes**

🔗 [https://mugomes.github.io](https://mugomes.github.io)

📺 https://youtube.com/@mugomesoficial

---

## Support

- GitHub: https://github.com/sponsors/mugomes
- More: https://mugomes.github.io/apoie.html

---

## License

Copyright (c) 2026 Murilo Gomes Julio

Licensed under the [MIT](https://github.com/mugomes/mgtemplate/blob/main/LICENSE) license.

All contributions to the MGTemplate are subject to this license.