
# Go Appointment Booking Microservices

Este proyecto es un sistema basado en microservicios para la gestión de citas (*appointment booking*) implementado en **Go**. Utiliza gRPC para la comunicación entre servicios y PostgreSQL como base de datos principal. El sistema permite a los usuarios registrar clientes y profesionales, programar citas, y enviar notificaciones por correo electrónico.

## Arquitectura

El proyecto sigue una arquitectura de microservicios, con los siguientes componentes principales:

1. **Authentication** (*Autenticación*):
   - Gestiona la autenticación de usuarios mediante JWT.
   - Responsable de login, registro y validación de tokens.

2. **Clients** (*Clientes*):
   - Maneja el registro, consulta y listado de clientes.
   - Almacena datos como nombre, email y teléfono en una base de datos PostgreSQL.

3. **Professionals** (*Profesionales*):
   - Administra la información de los profesionales (e.g., nombre, contacto).
   - Similar a *Clients*, usa PostgreSQL para persistencia.

4. **Agenda** (*Agenda*):
   - Gestiona la programación de citas y la disponibilidad de slots.
   - Permite crear slots, listar slots disponibles, reservar citas y consultar citas programadas.
   - Integra con el servicio de *Notificaciones* para enviar confirmaciones.

5. **Notifications** (*Notificaciones*):
   - Envía notificaciones por correo electrónico a clientes y profesionales tras reservar una cita.
   - Usa SMTP para el envío de emails y depende de *Clients* y *Professionals* para obtener datos.

## Tecnologías Utilizadas

- **Lenguaje**: Go (Golang)
- **Comunicación**: gRPC con Protobuf
- **Base de Datos**: PostgreSQL (usada en *Authentication*, *Clients*, *Professionals*, *Agenda*)
- **Notificaciones**: SMTP para envío de correos (*Notifications*)
- **Testing**: 
  - Unitarios: `testify/assert`, `testify/mock`, `DATA-DOG/go-sqlmock`
  - Integración: Pendiente
- **Dependencias**: GORM (ORM para PostgreSQL), net/smtp (para emails)

## Requisitos Previos

- Go 1.18 o superior
- PostgreSQL 13 o superior (para los servicios con base de datos)
- Protoc (Protocol Buffers Compiler) y plugins de Go para gRPC (`protoc-gen-go`, `protoc-gen-go-grpc`)
- Servidor SMTP configurado (e.g., Gmail SMTP para pruebas)

## Instalación

1. **Clonar el Repositorio**:
   ```bash
   git clone https://github.com/lpsaldana/go-appointment-booking-microservices.git
   cd go-appointment-booking-microservices

    Instalar Dependencias:
    Para cada microservicio:
    bash

    cd <microservice>  # e.g., cd authentication
    go mod tidy

    Configurar Variables de Entorno:
        Crea un archivo .env en cada microservicio con las configuraciones necesarias:
            DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME (para servicios con DB)
            SMTP_HOST, SMTP_PORT, SMTP_USER, SMTP_PASSWORD (para Notifications)
            Ejemplo para Notifications:

            SMTP_HOST=smtp.gmail.com
            SMTP_PORT=587
            SMTP_USER=tu-email@gmail.com
            SMTP_PASSWORD=tu-app-password

    Compilar Protos:
    Desde el directorio common/pb:
    bash

    protoc --go_out=. --go-grpc_out=. *.proto

    Ejecutar los Servicios:
    Para cada microservicio:
    bash

    cd <microservice>
    go run main.go

Pruebas
Cada microservicio incluye tests unitarios en test/unit/. Para ejecutarlos:
bash

cd <microservice>
go test ./test/unit -coverprofile=cover.out
go tool cover -html=cover.out -o cover.html

    Cobertura: ~80-90% para las capas service y repository (donde aplica).
    Mocks: Usamos sqlmock para simular la base de datos y testify/mock para dependencias externas (SMTP, gRPC).

Uso

    Autenticación:
        Registra un usuario y obtén un JWT para autenticar solicitudes a otros servicios.
    Clientes y Profesionales:
        Registra clientes y profesionales mediante sus respectivos endpoints gRPC.
    Agenda:
        Crea slots de disponibilidad para profesionales.
        Reserva citas vinculando clientes y slots.
        Lista citas programadas.
    Notificaciones:
        Al reservar una cita en Agenda, se envían correos al cliente y al profesional con los detalles.

Ejemplo de Flujo

    Un cliente se autentica y registra sus datos en Clients.
    Un profesional registra su información en Professionals.
    El profesional crea slots disponibles en Agenda.
    El cliente reserva una cita en un slot disponible.
    Agenda llama a Notifications para enviar confirmaciones por email al cliente y al profesional.

Licencia
Este proyecto está bajo la Licencia MIT (LICENSE).
