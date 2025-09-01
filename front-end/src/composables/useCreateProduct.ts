import {Color} from "@/lib/types";


type ColorsAPIResponse = {
    colors: Color[]
}

type ColorsComposableResponse = {
    status: number,
    colors?: Color[]
}

export const useCreateProduct = () => {

    const fetchColors = async (): Promise<ColorsComposableResponse> => {
        const res = await fetch(`${import.meta.env.VITE_BACKEND_URL}/colors`, {
            method: "GET",
        })

        const data = await res.json() as ColorsAPIResponse

        if (!res.ok) {
            console.error(res)
            return {
                status: res.status,
            }
        }

        return {
            status: res.status,
            colors: data.colors,
        }
    }
    const addProduct = async (product: FormData): Promise<number> => {
        const res = await fetch(`${import.meta.env.VITE_BACKEND_URL}/products`, {
            method: "POST",
            body: product,
            credentials: 'include',
        })

        if (!res.ok) {
            console.error(res)
            return res.status
        }

        return res.status
    }
    const editProduct = async (product: FormData): Promise<number> => {
        const res = await fetch(`${import.meta.env.VITE_BACKEND_URL}/products`, {
            method: "PUT",
            body: product,
            credentials: 'include',
        })

        if (!res.ok) {
            console.error(res)
            return res.status
        }

        return res.status
    }

    return {
        fetchColors,
        addProduct,
        editProduct,
    }
}