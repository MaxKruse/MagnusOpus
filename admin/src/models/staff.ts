import User from "./user";

class Staff {
    id: number;
    User: User;
    Role: string;

    constructor(id: number, user: User, role: string) {
        this.id = id;
        this.User = user;
        this.Role = role;
    }
}

export default Staff;