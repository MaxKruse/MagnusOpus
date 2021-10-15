import axios from "axios";
import storage from "./storageHandler"
import store from "./store"

const BASE_URL = "https://localhost/"

export default {
    get: (url: string, params?: any) => {
        return axios.get(BASE_URL + url, {
            params: params,
            headers: {
                "Authorization": "Bearer " + storage.sessionId()
            }
        })
    },
    post: (url: string, data: any) => {
        return axios.post(BASE_URL + url, data, {
            headers: {
                "Authorization": "Bearer " + storage.sessionId()
            }
        })
    },
    put: (url: string, data: any) => {
        return axios.put(BASE_URL + url, data, {
            headers: {
                "Authorization": "Bearer " + storage.sessionId()
            }
        })
    },
    delete: (url: string) => {
        return axios.delete(BASE_URL + url, {
            headers: {
                "Authorization": "Bearer " + storage.sessionId()
            }
        })
    }
}