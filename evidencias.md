# Evidencias — TP1

## 1. Push directo a main rechazado
<img width="728" height="480" alt="cap1" src="https://github.com/user-attachments/assets/40ce1565-1396-4ec4-9cff-2e08741fa4b1" />
GitHub rechaza el intento de push directo a `main` porque la rama está protegida
y la regla "Do not allow bypassing" alcanza también al administrador del repositorio
(en este caso, yo mismo). El error muestra `protected branch hook declined`, confirmando
que ni siquiera el dueño del repo puede saltear la protección.

## 2. El PR de la rama B no se puede mergear: conflicto
<img width="1337" height="670" alt="cap2" src="https://github.com/user-attachments/assets/90324872-87a6-45e8-8318-6dbb84561ded" />
Al crear el Pull Request de la rama `feature/titulo-b`, GitHub avisa que no se puede
mergear automáticamente porque hay conflictos con `main`. Esto ocurre porque la rama A
(ya mergeada) y la rama B modificaron la misma línea del `README.md` a partir del mismo
punto de partida, y Git no puede decidir cuál versión es la correcta.

## 3. Marcadores del conflicto
<img width="1337" height="488" alt="cap3" src="https://github.com/user-attachments/assets/e3377d84-cd29-440b-ab34-2b554d802710" />
Al abrir "Resolve conflicts" en el PR de la rama B, se muestran los marcadores de Git
que delimitan el conflicto: `<<<<<<<` marca el inicio de la versión de mi rama actual,
`=======` separa las dos versiones en disputa, y `>>>>>>>` marca el final de la versión
que ya está en `main`. Resolver el conflicto implicó decidir qué contenido conservar y
eliminar estos marcadores.

## 4. Release v1.0.0 publicada
<img width="1405" height="677" alt="cap4" src="https://github.com/user-attachments/assets/b5ddf2c0-2ad6-441b-8ddd-dc2e6b991034" />
La release `v1.0.0` publicada en GitHub, generada a partir del tag creado sobre `main`,
con las notas describiendo qué incluye esta primera versión estable del TP.
