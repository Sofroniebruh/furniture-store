import {ref} from "vue";
import {Product} from "@/lib/types";

const productsHistoryData = ref<Product[]>([])
const wishlistLoading = ref<boolean>(false);
const wishlistError = ref<boolean>(false);
const historyLoading = ref<boolean>(false);
const historyError = ref<boolean>(false);

export const useWishlistOrHistory = () => {
    const fetchData = async (endpoint: string) => {
        switch (endpoint) {
            case "wishlist":
                wishlistLoading.value = true
                wishlistError.value = false
                break;
            case "history":
                historyLoading.value = true
                historyError.value = false
                break;
            default:
                return
        }

        const res = await fetch(`${import.meta.env.VITE_USER_RELATED_SERVICE_URL}/user/${endpoint}`, {
            method: "GET",
            credentials: "include",
        })

        wishlistLoading.value = false
        historyLoading.value = false

        const data = await res.json()

        if (!res.ok) {
            endpoint == "wishlist" ? wishlistError.value = true :
                endpoint == "history" ? historyError.value = true : null;
            return
        }

        switch (endpoint) {
            case "history":
                productsHistoryData.value = data?.products ?? []
                break
            default:
                return;
        }
    }
    const addToWishlist = async (productId: string) => {
        wishlistError.value = false

        const res = await fetch(`${import.meta.env.VITE_USER_RELATED_SERVICE_URL}/user/wishlist`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                product_id: productId,
            }),
            credentials: "include",
        })

        const data = await res.json() as { product: Product }

        if (!res.ok) {
            wishlistError.value = true;
            return res.status
        }

        return res.status
    }
    const removeFromWishlist = async (productId: string) => {
        const res = await fetch(`${import.meta.env.VITE_USER_RELATED_SERVICE_URL}/user/wishlist`, {
            method: "DELETE",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                product_id: productId,
            }),
            credentials: "include",
        })

        const data = await res.json()

        if (!res.ok) {
            wishlistError.value = true;
            console.error(data?.error || "Internal server error")
            return res.status
        }

        return res.status
    }

    return {
        productsHistoryData,
        wishlistLoading,
        wishlistError,
        historyLoading,
        historyError,
        fetchWishlistData: fetchData,
        addToWishlist,
        removeFromWishlist,
    }
}