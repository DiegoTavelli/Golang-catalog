# go.sum — explicación

Equivalente a `package-lock.json` en npm o `yarn.lock` en yarn.
Go lo genera y actualiza automáticamente — nunca lo editás a mano.

---

## ¿Qué contiene?

Cada línea tiene este formato:
```
github.com/gin-gonic/gin v1.12.0 h1:HASH=
github.com/gin-gonic/gin v1.12.0/go.mod h1:HASH=
```

- El nombre del paquete y su versión exacta
- Un hash criptográfico (SHA-256) del contenido del paquete
- Otro hash del go.mod de ese paquete

---

## ¿Para qué sirve el hash?

Seguridad e integridad. Cuando alguien más clona tu proyecto y corre `go mod download`,
Go verifica que el paquete descargado tenga exactamente el mismo hash.
Si alguien modificó el paquete en el servidor (supply chain attack), el hash no coincide y Go rechaza la descarga.

Equivalente a: la verificación de integridad que hace yarn con los checksums en yarn.lock.

---

## ¿Lo commiteás al repositorio?

Sí, siempre. Igual que yarn.lock o package-lock.json.
Garantiza que todos en el equipo usen exactamente las mismas versiones con el mismo contenido.

---

## Lo que NO tenés en Go

En Node tenés `node_modules/` en el proyecto — una carpeta pesada que no se commitea.
En Go las dependencias se guardan en un cache global en `C:\Users\Diego\go\pkg\mod\`.
Tu proyecto no tiene carpeta de dependencias — está todo centralizado en tu máquina.
Esto significa que Go es mucho más liviano por proyecto.
