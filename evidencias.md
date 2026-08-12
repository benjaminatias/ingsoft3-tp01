# Evidencias — TP1

## 1. Push directo a main rechazado
![push rechazado](img/cap1.png)
GitHub rechaza el intento de push directo a `main` porque la rama está protegida
y la regla "Do not allow bypassing" alcanza también al administrador del repositorio
(en este caso, yo mismo). El error muestra `protected branch hook declined`, confirmando
que ni siquiera el dueño del repo puede saltear la protección.

## 2. El PR de la rama B no se puede mergear: conflicto
![aviso de conflicto](img/cap2.png)
Al crear el Pull Request de la rama `feature/titulo-b`, GitHub avisa que no se puede
mergear automáticamente porque hay conflictos con `main`. Esto ocurre porque la rama A
(ya mergeada) y la rama B modificaron la misma línea del `README.md` a partir del mismo
punto de partida, y Git no puede decidir cuál versión es la correcta.

## 3. Marcadores del conflicto
![marcadores de conflicto](img/cap3.png)
Al abrir "Resolve conflicts" en el PR de la rama B, se muestran los marcadores de Git
que delimitan el conflicto: `<<<<<<<` marca el inicio de la versión de mi rama actual,
`=======` separa las dos versiones en disputa, y `>>>>>>>` marca el final de la versión
que ya está en `main`. Resolver el conflicto implicó decidir qué contenido conservar y
eliminar estos marcadores.

## 4. Release v1.0.0 publicada
![release publicada](img/cap4.png)
La release `v1.0.0` publicada en GitHub, generada a partir del tag creado sobre `main`,
con las notas describiendo qué incluye esta primera versión estable del TP.
