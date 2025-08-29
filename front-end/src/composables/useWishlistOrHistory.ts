import {computed, ref} from "vue";
import {Product} from "@/lib/types";
import {useWishlistStore} from "@/stores/useWishlist";
import {useHistoryStore} from "@/stores/useHistory";
import {useRoute} from "vue-router";
import {useParams} from "@/composables/useParams";

export const useWishlistOrHistory = () => {
    const wishlistError = ref<boolean>(false);
    const wishlistLoading = ref(false);
    const historyLoading = ref<boolean>(false);
    const historyError = ref<boolean>(false);
    const totalPages = ref<number>(0);

    const route = useRoute()
    const {updatePage} = useParams()
    const limit = import.meta.env.VITE_PRODUCTS_PER_PAGE_LIMIT

    const filters = computed(() => ({
        page: Number(route.query.page) || 1,
    }))

    const wishlistStore = useWishlistStore()
    const historyStore = useHistoryStore()
    const {initFromBackend} = wishlistStore
    const {initHistoryFromBackend} = historyStore

    const setPage = async (page: number) => {
        updatePage(page)
    }

    const handleNext = async (pageNumber: number) => {
        await setPage(pageNumber)
    }

    const handlePrevious = async (pageNumber: number) => {
        await setPage(pageNumber)
    }

    const fetchData = async (endpoint: string) => {
        try {
            switch (endpoint) {
                case "wishlist":
                    wishlistLoading.value = true
                    wishlistStore.isLoadingWishlist = true
                    wishlistError.value = false
                    break;
                case "history":
                    historyLoading.value = true
                    historyStore.isLoadingHistory = true
                    historyError.value = false
                    break;
                default:
                    return
            }

            const baseUrl = computed(() => `${import.meta.env.VITE_USER_RELATED_SERVICE_URL}/user/${endpoint}?page=${filters.value.page}&limit=${limit}`)

            console.log(baseUrl.value)

            const res = await fetch(`${baseUrl.value}`, {
                method: "GET",
                credentials: "include",
            })

            const data = await res.json() as { products: Product[], total: number }

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

            totalPages.value = data.total

            switch (endpoint) {
                case "history":
                    initHistoryFromBackend(data?.products ?? [])
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
            switch (endpoint) {
                case "wishlist":
                    wishlistLoading.value = false
                    wishlistStore.isLoadingWishlist = false
                    break
                case "history":
                    historyLoading.value = false
                    historyStore.isLoadingHistory = false
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
        wishlistLoading,
        wishlistError,
        historyLoading,
        historyError,
        totalPages,
        currentPage: computed(() => filters.value.page),
        handlePrevious,
        handleNext,
        setPage,
        fetchWishlistData: fetchData,
        addToWishlist,
        removeFromWishlistComposable: removeFromWishlist,
    }
}