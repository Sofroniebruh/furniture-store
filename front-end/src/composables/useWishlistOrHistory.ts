import {ref} from "vue";
import {Product} from "@/lib/types";
import {useWishlistStore} from "@/stores/useWishlist";

export const useWishlistOrHistory = () => {
    // Move reactive refs INSIDE the composable function
    const productsHistoryData = ref<Product[]>([])
    const wishlistError = ref<boolean>(false);
    const wishlistLoading = ref(false);
    const historyLoading = ref<boolean>(false);
    const historyError = ref<boolean>(false);

    let {initFromBackend, isLoadingWishlist} = useWishlistStore()

    const fetchData = async (endpoint: string) => {
        try {
            switch (endpoint) {
                case "wishlist":
                    wishlistLoading.value = true
                    isLoadingWishlist = true // Also set store loading
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

            const data = await res.json()

            if (!res.ok) {
                switch (endpoint) {
                    case "wishlist":
                        wishlistError.value = true
                        break
                    case "history":
                        historyError.value = true
                        break
                }
                return
            }

            switch (endpoint) {
                case "history":
                    productsHistoryData.value = data?.products ?? []
                    break
                case "wishlist":
                    initFromBackend(data?.products ?? [])
                    break
                default:
                    return;
            }
        } catch (error) {
            console.error(`Error fetching ${endpoint}:`, error)
            switch (endpoint) {
                case "wishlist":
                    wishlistError.value = true
                    break
                case "history":
                    historyError.value = true
                    break
            }
        } finally {
            // Always reset loading states in finally block
            switch (endpoint) {
                case "wishlist":
                    wishlistLoading.value = false
                    isLoadingWishlist = false
                    break
                case "history":
                    historyLoading.value = false
                    break
            }
        }
    }

    const addToWishlist = async (productId: string) => {
        wishlistError.value = false

        try {
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
        } catch (error) {
            console.error('Error adding to wishlist:', error)
            wishlistError.value = true
            return 500
        }
    }

    const removeFromWishlist = async (productId: string) => {
        try {
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
        } catch (error) {
            console.error('Error removing from wishlist:', error)
            wishlistError.value = true
            return 500
        }
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