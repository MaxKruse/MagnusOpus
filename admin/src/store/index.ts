import User from "@/models/user";
import { createStore } from "vuex";

export interface State {
  user: User;
}

export default createStore<State>({
  state: {
    user: null
  },
  mutations: {
    setUser(state, user: User) {
      state.user = user;
    },
  },
  actions: {},
  modules: {},
});
