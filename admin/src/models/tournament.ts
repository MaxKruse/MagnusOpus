import User from "./user";
import Staff from "./staff";
import Round from "./round";

class Tournament {
    id: number;
    name: string;
    description: string;
    start_time: Date;
    end_time: Date;
    rounds: Round[];
    staffs: Staff[];
    registrations: User[];
    registration_start_time: Date;
    registration_end_time: Date;

    constructor(id: number, name: string, description: string, start_time: Date, end_time: Date, rounds: Round[], staffs: Staff[], registrations: User[], registration_start_time: Date, registration_end_time: Date) {
        this.id = id;
        this.name = name;
        this.description = description;
        this.start_time = start_time;
        this.end_time = end_time;
        this.rounds = rounds;
        this.staffs = staffs;
        this.registrations = registrations;
        this.registration_start_time = registration_start_time;
        this.registration_end_time = registration_end_time;
    }
}

export default Tournament;