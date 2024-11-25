# MailingGo

#### MICROSERVICIO DE ENVÍO DE NOTIFICACIONES

#### Este microservicio forma parte de un e-commerce, es el encargado de el envío de notificaciones via e-mail en ciertas situaciones

## Versión: 2.0

**Contact Information:**  
Joaquín Cortez - joacomateocortez@gmail.com - legajo: 46915

**SMTP Service Used:**

- Usa el servicio SMTP de Google para enviar correos electrónicos personalizados.
- Documentación de Google SMTP: [https://support.google.com/a/answer/176600?hl=en](https://support.google.com/a/answer/176600?hl=en)

---

## Casos de Uso

### Almacenamiento de Usuarios

#### Descripción del Caso de Uso:

Este caso de uso permite almacenar y gestionar usuarios que están suscritos a las notificaciones en el sistema `MailingGo`. Dependiendo de su ID, se puede actualizar su estado (habilitado o deshabilitado) para determinar si está suscrito a las notificaciones de correo electrónico.

#### Endpoints Usadas:

- **POST**: `/mailinggo/userSubscription`

#### Descripción de los Pasos:

1. **Recepción de Datos**: Se recibe la información de un usuario en el formato adecuado, con su ID y estado de habilitación.
2. **Verificación del Usuario**: Si el usuario ya existe, se actualiza su estado (habilitado o deshabilitado).
3. **Almacenamiento**: Si el usuario no existe, se crea un nuevo registro en la base de datos de `MailingGo` con el ID proporcionado y su estado de habilitación.
4. **Resultado**: El usuario se almacena o actualiza correctamente en la base de datos.

---

### Almacenamiento de Artículos en Oferta

#### Descripción del Caso de Uso:

Este caso de uso permite almacenar y gestionar los artículos que están en oferta dentro de `MailingGo`. Dependiendo del ID del artículo, se puede habilitar o deshabilitar su disponibilidad para ser enviado en futuras campañas de correos electrónicos.

#### Endpoints Usadas:

- **POST**: `/mailinggo/articleOffer`

#### Descripción de los Pasos:

1. **Recepción de Datos**: Se recibe la información de un artículo en oferta, incluyendo su ID y estado de habilitación.
2. **Verificación del Artículo**: Si el artículo ya existe en la base de datos, se actualiza su estado (habilitado o deshabilitado).
3. **Almacenamiento**: Si el artículo no existe, se crea un nuevo registro en la base de datos de `MailingGo` con el ID proporcionado y su estado de habilitación.
4. **Resultado**: El artículo se almacena o actualiza correctamente en la base de datos.

---

### Envío de Ofertas

#### Descripción del Caso de Uso:

Este caso de uso permite enviar correos electrónicos con todas las ofertas actuales a los usuarios que están suscritos a las notificaciones.

#### Endpoints Usadas:

- **POST**: `/mailinggo/offers`
- **GET**: `v1/user/:userId` _AUTH SERVICE_
- **GET**: `v1/article/::articleId` CATALOGGO SERVICE\_

#### Descripción de los Pasos:

1. **Obtención de Suscriptores**: Se obtienen todos los usuarios habilitados para recibir notificaciones desde la base de datos interna de MailingGo.
2. **Obtención de data de suscriptores**: Se comunica con el servicio de Auth para buscar toda la información de valor restante de los usuarios suscritos utilizando el endpoint `v1/user/:userId`
3. **Obtención de Artículos en Oferta**: Se obtienen todos los artículos que están marcados como ofertas y habilitados en la base de datos de MailingGo.
4. **Obtención de data de articulos**: Se comunica con el servicio de CatalogGo para buscar toda la información de valor restante de los artiulos en oferta utilizando el endpoint: `v1/articles/articleId`
5. **Envío de Correos Electrónicos**: Se envían correos electrónicos a los usuarios con los detalles de los artículos en oferta (nombre, descripción, precio).
6. **Resultado**: Los correos electrónicos son enviados correctamente a los usuarios suscritos.

---

### Recordatorio de Carrito Abandonado

#### Descripción del Caso de Uso:

Este caso de uso permite enviar recordatorios a los usuarios sobre los carritos abandonados. Si el usuario está habilitado para recibir notificaciones, se le enviará un recordatorio de los artículos que dejó en su carrito.

#### Endpoints Usadas:

- **POST**: `v1/malinggp/openCart`
- **GET**: `/v1/cart/:userID/not-empty` _CARTGGO SERVICE_

#### Descripción de los Pasos:

1. **Obtención de Usuarios**: Se obtienen todos los usuarios registrados y habilitados desde la base de datos interna de `MailingGo`.
2. **Verificación de Carrito**: Para cada usuario, se verifica si su carrito tiene al menos un artículo utilizando el endpoint `/v1/cart/:userID/not-empty`. Este endpoint devuelve un booleano indicando si el carrito tiene algún artículo.
3. **Envío de Recordatorios**: Si el carrito del usuario no está vacío, se envía un correo electrónico de recordatorio utilizando la plantilla `_cart_open.tmpl_`.
4. **Resultado**: Los correos electrónicos de recordatorio son enviados correctamente a los usuarios.

---

### Aviso de Cambio de Precio

#### Descripción del Caso de Uso:

Este caso de uso permite enviar alertas a los usuarios que tienen un carrito abierto con artículos cuyo precio ha sido modificado.

#### Endpoints Usadas:

- **No utiliza un endpoint público directamente**, este proceso se activa mediante mensajes de `RabbitMQ`.

- **GET**: `v1/updated-carts?updatedArticle=` _CARTGGO SERVICE_
- **GET**: `v1/user/:userId` _AUTH SERVICE_

#### Descripción de los Pasos:

1. **Recepción del Mensaje**: Cuando un artículo tiene un cambio de precio, se publica un mensaje en la cola de `RabbitMQ`.
2. **Obtención de Carritos Afectados**: Se obtienen todos los carritos que contienen el artículo con el precio modificado. Se usa el endpoint `v1/updated-carts?updatedArticle=`
3. **Verificación de Usuario**: Para cada carrito recuperado, se verifica si el usuario asociado coincide con algún usuario almacenado localmente en la base de datos de `MailingGo` y que esté suscrito.
4. **Búsqueda de Información**: Para cada usuario seleccionado, se busca su información en el servicio de Auth utilizando el endpoint: `v1/user/:userId`
5. **Envío de Alerta**: Se le envía un correo electrónico notificándole el cambio de precio en los artículos de su carrito utilizando la plantilla `_price_updates.tmpl_`.
6. **Resultado**: Los correos electrónicos de alerta sobre el cambio de precio son enviados correctamente a los usuarios.

---

## Modelo de Datos

### Subscribers (`subscribers`)

Los usuarios representan a las personas que interactúan con el sistema y tienen la capacidad de habilitar o deshabilitar notificaciones.

- **Atributos:**
  - `ID` (string): Identificador único del usuario.
  - `Subscribed` (boolean): Indica si el usuario está habilitado para recibir notificaciones.

---

### Offer (`articles_offer`)

Los artículos representan los productos u ofertas que se envían a los usuarios.

- **Atributos:**
  - `ID` (string): Identificador único del artículo.
  - `Offer` (boolean): Indica si el artículo está habilitado como oferta.

---

## Interfaz REST

### usersSubscription Endpoint

#### URL

`POST /mailinggo/userSubscription`

#### Description

Este endoint es usado para manejar las suscripciones de los ususarios. Si el usuario no existe crea una nueva instancia, si este existe, la única acción es la actualización del estado.

#### Body

```json
{
  "id": "string",
  "subscribed": "boolean"
}
```

#### Responses

**Código de estado**: `200 OK`

- **Posibles respuestas**
  - **Usuario creado con éxito**
  ```json
  {
    "message": "User successfully created",
    "id": "string"
  }
  ```
  - **Usuario actualizado con éxito**
  ```json
  {
    "message": "User  state succesfully updated",
    "state": "enabled" // o disabled
  }
  ```
  - **No se realizaron cambios**
  ```json
  {
    "message": "No changes were made"
  }
  ```
  **Código de estado**: `400 Bad request`

```json
{
  "message": "Invalid data"
}
```

**Código de estado**: `500 Internal Server Error`

```json
{
  "error": "Failed to [find user | insert user | update user state]"
}
```

### articleOffer Endpoint

#### URL

`POST /mailinggo/userSubscription`

#### Description

Este endoint es usado para manejar las ofertas de los articulos, si el artículo no existe, crea una nueva instancia, si este existe, la única accion es la actualizacion del estado.

#### Body

```json
{
  "id": "string",
  "offer": "boolean"
}
```

#### Responses

**Código de estado**: `200 OK`

- **Posibles respuestas**
  - **Usuario creado con éxito**
  ```json
  {
    "message": "Article Offer successfully created",
    "id": "string"
  }
  ```
  - **Usuario actualizado con éxito**
  ```json
  {
    "message": "Article state succesfully updated",
    "state": "enabled" // o disabled
  }
  ```
  - **No se realizaron cambios**
  ```json
  {
    "message": "No changes were made"
  }
  ```
  **Código de estado**: `400 Bad request`

```json
{
  "message": "Invalid data"
}
```

**Código de estado**: `500 Internal Server Error`

```json
{
  "error": "Failed to [find article | insert article | update article state]"
}
```

### Offers Endpoint

#### URL

`POST /mailinggo/offers`

#### Description

Este endoint es usado para enviar una lista de las ofertas activas a los usuario suscritos a notificaciones.

#### Header

Bearer token

#### Body

- No recibe ningun Body

#### Responses

**Código de estado**: `200 OK`

- **Emails enviados con exito**

```json
{
  "message": "Emails sent successfully"
}
```

**Código de estado**: `400 Bad request`

```json
{
  "message": "Invalid data"
}
```

**Código de estado**: `500 Internal Server Error`

```json
{
  "error": "Failed to send emails"
}
```

### OpenCartNotification Endpoint

#### URL

`POST /mailinggo/openCart`

#### Description

Este endoint es usado para enviar un recordatorio a todos los usuarios suscritos de su carro abierto cuando contiene al menos un articulo.

#### Header

Bearer token

#### Header

`Authorization: bearer token`

#### Body

- No recibe ningun Body

#### Responses

**Código de estado**: `200 OK`

- **Emails enviados con exito**

```json
{
  "message": "Emails sent successfully"
}
```

**Código de estado**: `400 Bad request`

```json
{
  "message": "Invalid data"
}
```

**Código de estado**: `500 Internal Server Error`

```json
{
  "error": "Failed to send emails"
}
```
