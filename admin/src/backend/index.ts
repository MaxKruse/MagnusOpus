import axios, { AxiosRequestConfig, Method } from "axios";
import cookies from "@/cookies";

import User from "../models/user"
import Tournament from "../models/tournament"

import CustomError from "@/models/CustomError";

const BASE_URL = "https://localhost";

export default {
    async GetSelf(callback: (data: any, error: CustomError | null) => void) {
        return await apiRequest<User>("/api/v1/self", "GET", {}, callback);
    },
    async GetUsers(callback: (data: any, error: CustomError | null) => void) {
        return await apiRequest<User[]>("/api/v1/users", "GET", {}, callback);
    },
    async GetTournaments(callback: (data: any, error: CustomError | null) => void) {
        return await apiRequest<Tournament[]>("/api/v1/tournaments", "GET", {}, callback);
    },
    async GetTournament(id: number, callback: (data: any, error: CustomError | null) => void) {
        return await apiRequest(`/api/v1/tournaments/${id}`, "GET", {}, callback);
    },
    async PostTournament(data: any, callback: (data: any, error: CustomError | null) => void) {
        return await apiRequest("/api/v1/tournaments", "POST", data, callback);
    },
    async PutTournament(id: number, data: any, callback: (data: any, error: CustomError | null) => void) {
        return await apiRequest(`/api/v1/tournaments/${id}`, "PUT", data, callback);
    },
    async DeleteTournament(id: number, callback: (data: any, error: CustomError | null) => void) {
        return await apiRequest(`/api/v1/tournaments/${id}`, "DELETE", {}, callback);
    },
    async AddStaff(id: number, data: any, callback: (data: any, error: CustomError | null) => void) {
        return await apiRequest(`/api/v1/tournaments/${id}/staff`, "POST", data, callback);
    },
    async AddRound(id: number, data: any, callback: (data: any, error: CustomError | null) => void) {
        return await apiRequest(`/api/v1/tournaments/${id}/rounds`, "POST", data, callback);
    },
    async ActivateRound(id: number, data: any, callback: (data: any, error: CustomError | null) => void) {
        return await apiRequest(`/api/v1/tournaments/${id}/rounds/activate`, "POST", data, callback);
    }
}

// stub for authenticated CRUD requests
async function apiRequest<Type>(url: string, method: Method, data: any, callback: (data: any, error: CustomError | null) => void) {
    const token = cookies.sessionToken();
    const headers = {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${token}`
    };

    const req: AxiosRequestConfig<Type> = {
        method,
        url: `${BASE_URL}${url}`,
        headers,
        data,
        responseType: "json"
    };

    axios.request<Type>(req).then((response) => {
        callback(response.data, null)
    }).catch((error) => {
        callback(null, {
            message: error
        });
    })
}
