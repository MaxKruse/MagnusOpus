import CustomError from "./CustomError";
import Session from "./session";

type User = null | {
  username: string;
  id: number;
  ripple_id: number;
  sessions: Session[];
};

export default User;
