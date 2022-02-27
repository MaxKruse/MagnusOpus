import User from "@/models/user";
import { defineStore } from "pinia";

export type UserState = {
  user: User | null;
}

export const useUserStore = defineStore("User", {
  state: () => {
    return {
      user: null
    } as UserState
  },
  actions: {
    setUser(user: User) {
      if (!user) return;
      this.user = user;
    }
  }
});
