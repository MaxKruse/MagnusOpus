# Todos

Any specific types are in reference to golangs default types.

- tournament:
    ```json
       {
           "id": "int",
           "name": "string",
           "description": "string",
           "start_time": "time.Time",
           "end_time": "time.Time",
           "rounds": [
                "round",
                "round",
                "round",
           ],
           "staffs": [
                "staff",
                "staff",
                "staff",
           ],
           "registration_start_time": "time.Time",
           "registartion_end_time": "time.Time",
           "registrations": [
                "user"
           ]
       }
    ```

- round:
    ```json
       {
           "id": "int",
           "name": "string",
           "description": "string",
           "start_time": "time.Time",
           "end_time": "time.Time",
           "download_path": "Download path to file of this round",
        }
    ```

- staff:
    ```json
       {
           "id": "int",
           "User": "User",
           "role": "string",
       }
    ```

- rating:
    ```json
       {
           "id": "int",
           "scores": [
               "score",
               "score",
               "score",
           ]
        }
    ```

- score:
    ```json
       {
           "id": "int",
           "judge_id": "int",
           "score": "int",
           "category": "string",
        }
    ```


- beatmap:
    ```json
       {
           "id": "int",
           "hash": "string",
           "to_use": "bool",
        }
    ```

### User Spec

- [x] [GET] /oauth: Redirect user to oauth of the server, e.g. https://localhost/oauth/ripple redirects to ripple's oauth page
  - Valid server: ripple

- [x] [GET] /oauth/{server}/logout: Logout user for the server
  - Valid server: ripple

- [x] [GET] /api/v1/tournaments: Get all tournaments.

- [x] [GET] /api/v1/tournament/{id}: Get tournament by id.

- [x] [GET] /api/v1/me: Get personal data, including all tournament data, beatmaps and tokens. Literally everything

- [x] [GET] /api/v1/self: Get personal data, including all tournament data, beatmaps and tokens. Literally everything

- [X] [POST] /api/v1/{tournament}/registration: Signup for a tournament

- [X] [DELETE] /api/v1/{tournament}/registration: Unsignup for a tournament

- [X] [GET] /api/v1/{tournament}/rounds/{round_name}/beatmaps: Receive uploaded maps, or judgable maps for judges

- [ ] [DELETE] /api/v1/{tournament}/rounds/{round_name}/beatmaps/{beatmap_hash}: Delete my uploaded map

- [ ] [POST] /api/v1/{tournament}/rounds/{round_name}/beatmaps: Upload a beatmap
  - Required fields:
    - file: .osu file
  - This will replace the oldest available file if the limit of 5 submittions is reached

### Admin Spec

All of the below endpoints require admin privileges. This will be enforced by middlewares.

If you're not an admin, you won't be able to access these endpoints. If you believe you are missing an admin privilege, please contact the tournament host.

- [X] [POST] /api/v1/tournaments: Create a tournament
  - Required fields:
    - same fields as tournament
  - Returns:
    - tournament

- [x] [PUT] /api/v1/tournaments/{id}: in-place replace an tournament by id
  - Required fields:
    - same fields as tournament

- [X] [DELETE] /api/v1/tournaments/{id}: Delete a tournament

- [x] [POST] /api/v1/tournaments/{id}/rounds/create: Create a round in the tournament
  - Required fields:
    - same fields as round
  - Returns:
    - round

- [x] [POST] /api/v1/tournaments/{id}/round/{round_name}/activate: Start a round

- [X] [POST] /api/v1/tournaments/{id}/staff: Add a staff member
  - Required fields:
    - user_id: Internal user id of the user
    - role: "admin", "judge" or "mod"

- [ ] [DELETE] /api/v1/tournaments/{id}/staff: Remove a staff member
  - Required fields:
    One of the following:
      - ripple_id: Ripple id of the staff member
      - bancho_id: Bancho id of the staff member

- [ ] [POST] /api/v1/tournaments/{id}/round/{round_id}/map/rating: Judge a map
  - Required fields:
    - map_id: Map entry id
    - scores: Scores of the map


    