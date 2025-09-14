import {defineStore} from "pinia";
import {ref} from "vue";

export type User = {
    id: string;
    email: string;
    roles: string[];
}

type AuthResponse = {
    created?: User;
    user?: User;
    message?: string;
    error?: string;
}

export type UserCredentialsInput = {
    email: string;
    password: string;
}

type AuthBody = {
    status: number;
    body?: AuthResponse;
}

export const useAuthStore = defineStore("authStore", () => {
    const isAuthenticated = ref<boolean>(false);
    const user = ref<User | null>(null);
    const error = ref<string>("")
    const loading = ref<boolean>(false);
    const isAdmin = ref<boolean>(false);

    const handleSendingAuthData = async (endpoint: string, body: UserCredentialsInput): Promise<AuthBody> => {
        try {
            error.value = ""
            loading.value = true;
            const res = await fetch(`${import.meta.env.VITE_AUTH_SERVICE_URL}/${endpoint}`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(body),
                credentials: "include"
            })
            loading.value = false;

            const data = await res.json() as AuthResponse;

            if (data?.user?.roles?.some(role => role === "admin")) {
                isAdmin.value = true;
            }

            if (!res.ok) {
                error.value = res.status === 401 && endpoint.endsWith("login") ? "Password and/or email are incorrect" : data?.error || "Unknown error";
                console.error("Error: ", error.value);
                return {
                    status: res.status,
                    body: data
                };
            }

            if (endpoint.endsWith("login") && data?.user) {
                isAuthenticated.value = true;
                user.value = data.user
            }

            return {
                status: res.status,
                body: data
            }
        } catch (e) {
            loading.value = false;
            error.value = "Network error or server unavailable";
            console.error("Network error: ", e);
            return {
                status: 500,
                body: { error: "Network error" }
            };
        }
    }
    const fetchUser = async () => {
        try {
            error.value = ""
            const refreshRes = await fetch(`${import.meta.env.VITE_AUTH_SERVICE_URL}/refresh`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                credentials: "include"
            })

            if (!refreshRes.ok) {
                isAuthenticated.value = false;
                user.value = null;
                return
            }

            const res = await fetch(`${import.meta.env.VITE_USER_RELATED_SERVICE_URL}/user`, {
                method: "GET",
                credentials: "include"
            })

            if (!res.ok) {
                isAuthenticated.value = false;
                user.value = null;
                return
            }

            const data = await res.json() as { user: User }

            if (data?.user?.roles?.some(role => role === "admin")) {
                isAdmin.value = true;
            }

            isAuthenticated.value = true
            user.value = data.user
        } catch (e) {
            console.error(e)
            isAuthenticated.value = false;
            user.value = null;
        }
    }
    const verifyCode = async (code: string, email: string) => {
        try {
            error.value = ""
            loading.value = true;
            const res = await fetch(`${import.meta.env.VITE_AUTH_SERVICE_URL}/verify`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    email: email,
                    code: code,
                }),
                credentials: "include"
            })

            const data = await res.json()

            if (!res.ok) {
                error.value = data?.error || "Unknown error";
                console.error("Error: ", error.value);
                return false
            }

            if (data?.verified) {
                isAuthenticated.value = true
                user.value = data.verified
            }

            return true
        } catch (e) {
            console.error(e)
            error.value = "Network error or server unavailable";
            return false
        } finally {
            loading.value = false
        }
    }
    const resendCode = async (email: string) => {
        try {
            error.value = ""
            const res = await fetch(`${import.meta.env.VITE_AUTH_SERVICE_URL}/resend`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    email: email,
                })
            })

            const data = await res.json()

            if (!res.ok) {
                error.value = data?.error || "Unknown error";
                console.error("Error: ", error.value);
                return res.status
            }

            return res.status
        } catch (e) {
            console.error(e)
            return 500
        }
    }
    const logout = async () => {
        try {
            error.value = ""
            const res = await fetch(`${import.meta.env.VITE_AUTH_SERVICE_URL}/logout`, {
                method: "POST",
                credentials: "include"
            })

            const data = await res.json()

            if (!res.ok) {
                error.value = data?.error || "Unknown error";
                console.error("Error: ", error.value);
                return
            }

            isAuthenticated.value = false
            isAdmin.value = false;
            user.value = null
        } catch (e) {
            console.error(e)
        }
    }
    const registration = async (body: UserCredentialsInput): Promise<AuthBody> => {
        try {
            return await handleSendingAuthData("registration", body)
        } catch (e) {
            console.error(e)
            return {
                status: 500,
                body: { error: "Network error" }
            }
        }
    }
    const login = async (body: UserCredentialsInput): Promise<AuthBody> => {
        try {
            return await handleSendingAuthData("login", body)
        } catch (e) {
            console.error(e)
            return {
                status: 500,
                body: { error: "Network error" }
            }
        }
    }

    return {
        user,
        isAuthenticated,
        error,
        loading,
        isAdmin,
        fetchUser,
        registration,
        resendCode,
        verifyCode,
        logout,
        login,
    }
})