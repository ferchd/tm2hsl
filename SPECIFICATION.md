# tm2hsl v0 - Supported Subset

## Supported
- Reglas `match` con regex básica
- Reglas `begin`/`end` con regex básica
- `contentName` para scope interior
- `captures` con nombres simples
- Includes: `$self`, `$base`
- Comentarios en línea y bloque

## Not Supported (v0)
- Repository (`#reference`)
- Captures en `begin`/`end`
- Reglas `while`
- Back-references en regex
- Lookahead/lookbehind complejo
- Patrones anidados profundos (>3 niveles)

## Comportamiento en Features No Soportados
1. Modo estricto: error en compilación
2. Modo permisivo: warning y conversión aproximada
3. Documentación clara de limitaciones

## Garantías
- Mismo input → mismo bytecode (determinismo)
- Bytecode estable entre versiones v0.x
- Error temprano en features no soportados