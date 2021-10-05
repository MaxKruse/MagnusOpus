# Todos

- [x] Support for signup using ripple
- [x] Save user data
- [ ] Save beatmap 

## Spec

### Response Types

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
           "UserID": "int",
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

### User Spec

- [x] [GET] /oauth: Redirect user to oauth of the server, e.g. https://localhost/oauth/ripple redirects to ripple's oauth page
  - Valid server: ripple

- [x] [GET] /oauth/{server}/logout: Logout user for the server
  - Valid server: ripple

- [x] [GET] /api/v1/tournaments: Get all tournaments.

- [x] [GET] /api/v1/tournament/{id}: Get tournament by id.

- [x] [GET] /api/v1/me: Get personal data, including all tournament data, beatmaps and tokens. Literally everything

- [x] [GET] /api/v1/self: Get personal data, including all tournament data, beatmaps and tokens. Literally everything

- [ ] [POST] /api/v1/{tournament}/signup: Signup for a tournament

- [ ] [POST] /api/v1/upload: Upload a beatmap
  - Required fields:
    - file: .osu file
    - name: Beatmap name
    - tournament_id: Tournament id

### Admin Spec

All of the below endpoints require admin privileges. This will be enforced by middlewares.

If you're not an admin, you won't be able to access these endpoints. If you believe you are missing an admin privilege, please contact the tournament host.

- [ ] [POST] /api/v1/tournaments: Create a tournament
  - Required fields:
    - name: Tournament name
    - start_time: Start time of the tournament
    - end_time: End time of the tournament
    - server: Server to run the tournament on
  - Returns:
    - id: Id of the tournament

- [ ] [PATCH] /api/v1/tournaments/{id}: Update a tournament
  - Required fields:
    - Any one of the above fields provided will be updated

- [ ] [DELETE] /api/v1/tournaments/{id}: Delete a tournament

- [ ] [POST] /api/v1/tournaments/{id}/start: Start a tournament

- [ ] [POST] /api/v1/tournaments/{id}/end: End a tournament

- [ ] [POST] /api/v1/tournaments/{id}/round/create: Create a round in the tournament
   - Required fields:
     - file: Example .osz file, including timed .osu file and .mp3 file
     - name: Name of the round
     - description: Description of the round
     - start_time: Start time of the round
     - end_time: End time of the round

- [ ] [POST] /api/v1/tournaments/{id}/round/{round_id}/start: Start a round

- [ ] [POST] /api/v1/tournaments/{id}/round/{round_id}/end: End a round

- [ ] [POST] /api/v1/tournaments/{id}/staff: Add a staff member
  - Required fields:
    One of the following:
      - ripple_id: Ripple id of the staff member
      - bancho_id: Bancho id of the staff member
    - role: Role of the staff member
      - Valid roles:
        - owner (delete tournament, + everything else)
        - admin (full edits)
        - mod (rounds)
        - judge (judge maps)

- [ ] [DELETE] /api/v1/tournaments/{id}/staff: Remove a staff member
  - Required fields:
    One of the following:
      - ripple_id: Ripple id of the staff member
      - bancho_id: Bancho id of the staff member

- [ ] [POST] /api/v1/tournaments/{id}/round/{round_id}/map/rating: Judge a map
  - Required fields:
    - map_id: Map entry id
    - scores: Scores of the map


    