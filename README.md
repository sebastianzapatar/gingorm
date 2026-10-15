# CRUD pequeño: Go + GORM + SQLite

API REST local de productos. Usa Gin para las rutas HTTP, GORM para las
consultas y SQLite para guardar los datos en un archivo.

## Estructura

```
crud-go-sqlite/
├── main.go                         # Punto de entrada: conecta la base y arranca el servidor
├── go.mod / go.sum                 # Módulo y dependencias
├── pruebas.http                    # Solicitudes de ejemplo
└── internal/
    ├── models/producto.go          # Structs Producto y Entrada, y la validación
    ├── database/database.go        # Abre SQLite y ejecuta AutoMigrate
    ├── handlers/productos.go       # Handlers del CRUD (crear, listar, obtener, ...)
    ├── handlers/respuestas.go      # Helper para responder errores en JSON
    └── routes/routes.go            # Enrutador Gin y registro de rutas
```

Cada paquete tiene una sola responsabilidad. `internal/` impide que otros módulos
importen estos paquetes; solo este proyecto puede usarlos.

## Requisitos

- Go 1.25 o superior (lo exige Gin v1.12): https://go.dev/doc/install
- Compilador C disponible y CGO habilitado. El driver oficial de SQLite usa CGO.
  En macOS: `xcode-select --install` si no tienes las herramientas de desarrollo.
  En Ubuntu: `sudo apt install build-essential`.
  En Windows x64: instala MSYS2 (https://www.msys2.org/), abre UCRT64 y ejecuta
  `pacman -S mingw-w64-ucrt-x86_64-gcc`; agrega `C:\msys64\ucrt64\bin` al PATH
  si usaste la ubicación predeterminada. Abre de nuevo tu terminal.
- Verifica `go version`, `gcc --version` y `go env CGO_ENABLED` (debe mostrar 1).
  Si muestra 0, usa `go env -w CGO_ENABLED=1`.

No necesitas instalar un servidor SQLite, crear una base manualmente ni usar Docker.

## Ejecutar el código descargado

Descomprime el ZIP, entra en `crud-go-sqlite` y ejecuta:

```bash
go mod download
go run .
```

Mantén esa terminal abierta y prueba desde Postman o desde una segunda terminal.
La API escucha en `http://localhost:8080`. Detén el proceso con Ctrl+C.

## Crearlo desde cero

```bash
mkdir crud-go-sqlite
cd crud-go-sqlite
go mod init ejemplo/crud-go-sqlite
go get github.com/gin-gonic/gin@v1.12.0 gorm.io/gorm@v1.31.1 gorm.io/driver/sqlite@v1.6.0
```

Crea las carpetas y archivos de la sección [Estructura](#estructura) con el
contenido de este repositorio. Luego ejecuta:

```bash
go mod tidy
go run .
```

`go.mod` declara el módulo y las dependencias; `go.sum` guarda hashes de módulos.
Las versiones fijadas hacen que el ejemplo use las mismas dependencias.

## Rutas

| Método | Ruta | Resultado esperado |
| --- | --- | --- |
| POST | /productos | Crear: 201 y producto con ID generado |
| GET | /productos | Listar: 200 y arreglo JSON |
| GET | /productos/{id} | Consultar uno: 200 o 404 |
| PUT | /productos/{id} | Reemplazar nombre y precio: 200 o 404 |
| DELETE | /productos/{id} | Eliminar: 204 sin cuerpo, o 404 |

El precio se expresa en pesos enteros: `80000` representa 80.000 pesos.
POST y PUT requieren `nombre` no vacío y `precio` entero mayor o igual a cero.
PUT requiere los dos campos. El ID lo genera la base de datos y no se envía
en el cuerpo del POST o PUT. Un cuerpo inválido o un ID no positivo devuelve 400.

## Probar con Postman

1. POST `http://localhost:8080/productos`.
2. Body -> raw -> JSON. Envía `{"nombre":"Teclado","precio":80000}`.
3. Copia el `id` de la respuesta. Con una base nueva será 1.
4. GET `/productos`, luego GET `/productos/1`.
5. PUT `/productos/1` con `{"nombre":"Teclado mecanico","precio":95000}`.
6. DELETE `/productos/1`; debe devolver 204 sin cuerpo.
7. GET `/productos/1`; ahora debe devolver 404.

En Postman, usa la URL completa para cada ruta. Seleccionar JSON agrega el
header `Content-Type: application/json`. `pruebas.http` incluye estas solicitudes
y casos de validación como referencia.

## Pruebas con curl (macOS, Linux o Bash)

```bash
curl -i -X POST http://localhost:8080/productos \
  -H 'Content-Type: application/json' \
  -d '{"nombre":"Teclado","precio":80000}'

curl -i http://localhost:8080/productos
curl -i http://localhost:8080/productos/1

curl -i -X PUT http://localhost:8080/productos/1 \
  -H 'Content-Type: application/json' \
  -d '{"nombre":"Teclado mecanico","precio":95000}'

curl -i -X DELETE http://localhost:8080/productos/1
curl -i http://localhost:8080/productos/1
```

Usa el ID devuelto al crear. En Windows puedes usar Postman para evitar diferencias
en las comillas y continuaciones de línea de PowerShell.

## Entender el código

- `main.go` solo une las piezas: `database.Conectar`, `handlers.NuevoProductos`
  y `routes.Nuevo`. Los paquetes se importan con la ruta del módulo, por ejemplo
  `ejemplo/crud-go-sqlite/internal/models`.
- `Producto` (en `models`) representa una fila. GORM crea la tabla `productos` con `id`, `nombre`
  y `precio`. Los tags `json` definen los nombres en la API.
- `Entrada` define los campos que el cliente puede enviar. `Precio *int64` permite
  distinguir un precio omitido (`nil`) de un precio cero. `*e.Precio` lee su valor.
- `Entrada.Validar` limpia espacios del nombre y comprueba los campos obligatorios.
- `database.Conectar` usa `gorm.Open` para abrir o crear `productos.db` y
  `AutoMigrate` para crear la tabla que falte.
- `handlers.Productos` guarda la conexión `*gorm.DB` y sus métodos (`Crear`,
  `Listar`, `Obtener`, `Actualizar`, `Eliminar`) atienden cada ruta. Así no hace
  falta una variable global para la base de datos.
- Solo los nombres que empiezan con mayúscula (`Producto`, `Crear`, ...) son
  visibles desde otros paquetes; `leer`, `buscar` o `fallo` son internos.
- `leer` convierte el JSON en una estructura y llama a `Validar`.
- `buscar` convierte el ID a número, consulta por clave primaria y distingue
  un registro inexistente (404) de un error de base de datos (500).
- `Create`, `Find`, `First`, `Updates` y `Delete` generan el SQL del CRUD.
- `Select("Nombre", "Precio").Updates(p)` incluye explícitamente el precio cero.
  `Updates` con un struct omite valores cero por defecto si no se seleccionan.
- `routes.Nuevo` crea el motor con `gin.Default()`, que ya trae los middleware
  `Logger` (imprime cada solicitud en la terminal) y `Recovery` (convierte un
  `panic` en un 500). `r.Group("/productos")` agrupa las rutas con ese prefijo y
  `/:id` declara un parámetro de ruta.
- Cada handler recibe `c *gin.Context`, que reúne la solicitud y la respuesta:
  `c.Param("id")` lee el ID, `c.ShouldBindJSON(&e)` lee el cuerpo y
  `c.JSON(estado, datos)` responde en JSON. `c.Status(204)` responde sin cuerpo.
- `fallo` responde los errores con la forma `{"error": "..."}` usando `gin.H`,
  un atajo de `map[string]any`.
- `&p` pasa la dirección del producto: GORM puede rellenarlo, por ejemplo con el ID.
- `err != nil` comprueba un error; `ok` indica si un helper terminó correctamente.

`productos.db` se crea en la carpeta desde la que ejecutas el programa. Los datos
se conservan al reiniciar la API. La eliminación es física: este modelo no tiene
`DeletedAt`. El ejemplo es didáctico y se ejecuta solamente en el equipo local.

## Problemas comunes

- `go: command not found`: instala Go y vuelve a abrir la terminal.
- `gcc not found` o error de CGO: instala el compilador C y comprueba CGO_ENABLED.
- `requires go >= 1.25`: actualiza Go; Gin v1.12 no compila con versiones anteriores.
- Mensajes `[GIN-debug] [WARNING]` al arrancar: son normales en modo debug. Para
  ocultarlos ejecuta con `GIN_MODE=release go run .`.
- `address already in use`: detén otra API que use 8080 o cambia el puerto
  en `r.Run` (en `main.go`) y en las URLs de prueba.
- 400: revisa JSON, campos obligatorios, precio entero no negativo e ID positivo.
- 404: usa el ID recibido al crear; puede haber sido eliminado o no existir.
- `[]`: la lista está vacía. Crea un producto primero.

## Documentación oficial

- Go, instalación: https://go.dev/doc/install
- Gin, documentación: https://gin-gonic.com/en/docs/
- Gin, repositorio: https://github.com/gin-gonic/gin
- GORM, SQLite: https://gorm.io/docs/connecting_to_the_database.html#SQLite
- GORM, migraciones: https://gorm.io/docs/migration.html
- GORM, consultas: https://gorm.io/docs/query.html
- GORM, actualizaciones: https://gorm.io/docs/update.html
- Driver SQLite y CGO: https://github.com/mattn/go-sqlite3#installation
- MSYS2, instalación: https://www.msys2.org/
