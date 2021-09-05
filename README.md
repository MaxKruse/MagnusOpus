# Magnus opus

[The Magnus opus](https://en.wikipedia.org/wiki/Magnum_opus_(alchemy)) is a collection of alchemical recipes. In this context, it is the codename for the project.

It's stages are:

- **Nigredo**, the blackening (frontend)
- **Albedo**, the whitening (backend)
- **Citrinitas**, the yellowing (admin panel)
- **Rubedo**, the reddening (tbd)

The Magnus opus is a full-stack service aimed at providing the interface and backend for organizing mapping tournaments in a controlled fashion. This includes features such as:

- [ ] **Registration**, allowing users to register for tournaments
- [ ] **Tournament creation**, allowing users to create tournaments
- [ ] **Tournament management**, allowing moderators to moderate and edit tournaments
- [ ] **Tournament rounds**, allowing moderators to create rounds in tournaments
- [ ] **Map Uploads**, allowing users to upload and view their maps
- [ ] **Map management**, allowing moderators to see submittions
- [ ] **Map judging**, allowing judges to see and rate submittions anonymously
- [ ] **Staff management**, allowing admins to manage staff members

## Nigredo

The frontend of this service includes all nesseccary frontend components to allow users to register, create tournaments, and manage tournaments.

It will also include an OAuth2 login system to allow users to login with their [ripple](https://ripple.moe) accounts.

## Albedo

The backend of this service includes all nesseccary backend components to allow moderators to create tournaments, manage tournaments, and manage rounds.

It will also provide for any needed backend components to allow users to upload and manage maps, as well as saving judgements in a database.

The database will be multi-client compliant, meaning users will only ever see what they are supposed to see. The same goes for judges. Only staff gets clear views.

Also noteworthy, the backend serves all files.

## Citrinitas

The admin panel will have a seperate view and requires a [ripple](https://ripple.moe) login.

The initial config will allow for 1 admin, that can then add more sub-admins and moderators to the system. These can be many2many relationships.

This admin panel is purely for creating tournaments, changing the active rounds and adding staff to tournaments.

## Technical Specifications

The following ports will be exposed by docker-compose:

    5000: The backend

The follow paths will be exposed:

    / (Root): frontend
    /admin  : admin
    /api/v1 : backend api

The port for the database will not be exposed.

It is suggested to use [docker-compose](https://docs.docker.com/compose/overview/) to manage the docker containers.

For https, we recommend using [caddy](https://caddyserver.com/). An example `Caddyfile` is provided.

### Frontend

The frontend has no additional configurations needed.

### Backend

The backend includes a `backend.env` file that contains the backend-specific configs.

### Admin Panel

The admin panel has no additional configurations needed.

### Database

The database is a postgres database. The documentation can be found [here](https://hub.docker.com/_/postgres).

The `postgres.env` file contains the database-specific configs used in the actual database and backend containers.

It is suggested to use `$ pwgen 32 1 -sB` to generate a secure password for the database.

## Usage (development)

    git clone 
    cd magnus-opus
    docker-compose up

## Usage (production)

    git clone
    cd magnus-opus
    docker-compose -f docker-compose.production.yml up -d

## Credits

Magnus Opus is a project by [@MaxKruse](https://github.com/MaxKruse) developed mainly for the [Ripple](https://ripple.moe) osu! server.