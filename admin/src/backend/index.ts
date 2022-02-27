import axios, { AxiosRequestConfig, Method } from "axios";
import cookies from "@/cookies";

import User from "../models/user";
import Tournament from "../models/tournament";

import CustomError from "@/models/CustomError";

const BASE_URL = "https://localhost";

export default {
  async GetSelf() {
    return await apiRequest<User>("/api/v1/self", "GET", {});
  },
  async GetUsers() {
    return await apiRequest<User[]>("/api/v1/users", "GET", {});
  },
  async GetTournaments() {
    return await apiRequest<Tournament[]>("/api/v1/tournaments", "GET", {});
  },
  async GetTournament(id: number) {
    return await apiRequest<Tournament>(`/api/v1/tournaments/${id}`, "GET", {});
  },
  async PostTournament(data: any) {
    return await apiRequest("/api/v1/tournaments", "POST", data);
  },
  async PutTournament(id: number, data: any) {
    return await apiRequest(`/api/v1/tournaments/${id}`, "PUT", data);
  },
  async DeleteTournament(id: number) {
    return await apiRequest(`/api/v1/tournaments/${id}`, "DELETE", {});
  },
  async AddStaff(id: number, data: any) {
    return await apiRequest(`/api/v1/tournaments/${id}/staff`, "POST", data);
  },
  async AddRound(id: number, data: any) {
    return await apiRequest(`/api/v1/tournaments/${id}/rounds`, "POST", data);
  },
  async ActivateRound(id: number, data: any) {
    return await apiRequest(
      `/api/v1/tournaments/${id}/rounds/activate`,
      "POST",
      data
    );
  },
};

// stub for authenticated CRUD requests
async function apiRequest<Type>(url: string, method: Method, data: any) {
  const token = cookies.sessionToken();
  const headers = {
    "Content-Type": "application/json",
    Authorization: `Bearer ${token}`,
  };

  const req: AxiosRequestConfig<Type> = {
    method,
    url: `${BASE_URL}${url}`,
    headers,
    data,
    responseType: "json",
  };

  return await (
    await axios.request<Type>(req)
  ).data;
}
