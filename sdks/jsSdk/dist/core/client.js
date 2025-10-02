"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetMeClient = void 0;
const axios_1 = __importDefault(require("axios"));
class GetMeClient {
    constructor(socketPath) {
        this.socketPath = socketPath;
        this.get = (key) => __awaiter(this, void 0, void 0, function* () {
            try {
                const response = yield this.axiosClient.get(`/get`, {
                    params: { key },
                });
                return response.data.value;
            }
            catch (error) {
                console.error("Error fetching value:", error);
                throw error;
            }
        });
        this.put = (key, value) => __awaiter(this, void 0, void 0, function* () {
            try {
                yield this.axiosClient.post(`/put`, { Key: key, Value: value });
            }
            catch (error) {
                console.error("Error storing value:", error);
                throw error;
            }
        });
        this.delete = (key) => __awaiter(this, void 0, void 0, function* () {
            try {
                yield this.axiosClient.delete(`/delete`, {
                    params: { key },
                });
            }
            catch (error) {
                console.error("Error deleting value:", error);
                throw error;
            }
        });
        this.axiosClient = axios_1.default.create({
            baseURL: "http://unix",
            socketPath: process.env.GETME_SOCKET_PATH || '/tmp/getMeStore/getMe.sock',
            headers: {
                "Content-Type": "application/json",
            },
        });
        console.log("GetMeClient initialized");
    }
}
exports.GetMeClient = GetMeClient;
