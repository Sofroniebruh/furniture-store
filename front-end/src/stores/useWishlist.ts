import {defineStore} from "pinia";
import {ref, watch} from "vue";
import {Product} from "@/lib/types";

export const useWishlistStore = defineStore("wishlistStore", () => {
    const itemsInWishlist = ref<Product[]>([])
    const isLoadingWishlist = ref<boolean>(false)
    const isInitialized = ref<boolean>(false)

    const isItemInWishlist = (id: string): boolean => {
        return itemsInWishlist.value.some((item) => item.id === id)
    }

    const initFromBackend = (values: Product[]) => {
        itemsInWishlist.value = values
        isInitialized.value = true
        if (itemsInWishlist.value.length !== 0) {
            saveToLocalStorage()
        }
    }

    const addToWishlist = (product: Product) => {
        if (!isItemInWishlist(product.id)) {
            itemsInWishlist.value = [...itemsInWishlist.value, product]
            saveToLocalStorage() // Auto-save when adding
            console.log("Item: " + product.name + " was added")
        } else {
            console.warn("Item: " + product.name + " was already added")
        }
    }

    const removeFromWishlist = (productOrId: Product | string) => {
        const productId = typeof productOrId === "string" ? productOrId : productOrId.id
        if (isItemInWishlist(productId)) {
            itemsInWishlist.value = itemsInWishlist.value.filter((item) => item.id !== productId)
            saveToLocalStorage() // Auto-save when removing
            console.log("Item with ID: " + productId + " was removed")
        } else {
            console.warn("Item with ID: " + productId + " was not found")
        }
    }

    const toggleWishlist = (product: Product) => {
        console.log("Items: ", itemsInWishlist.value)
        console.log(!isItemInWishlist(product.id))
        if (!isItemInWishlist(product.id)) {
            addToWishlist(product)
        } else {
            removeFromWishlist(product)
        }
    }

    const initFromLocalStorage = () => {
        try {
            const stored = localStorage.getItem("wishlistItem")
            if (stored && stored !== "{}") {
                const parsedData = JSON.parse(stored)
                // Ensure we're dealing with an array
                if (Array.isArray(parsedData)) {
                    itemsInWishlist.value = parsedData
                    console.log("Wishlist initialized from localStorage: ", parsedData)
                }
            }
            isInitialized.value = true
        } catch (e) {
            console.error("Failed to initialize wishlist from localStorage", e)
            itemsInWishlist.value = []
            isInitialized.value = true
        }
    }

    const saveToLocalStorage = () => {
        try {
            localStorage.setItem("wishlistItem", JSON.stringify(itemsInWishlist.value))
        } catch (e) {
            console.error("Failed to save wishlist item to local storage: ", e)
        }
    }

    const clearWishlist = () => {
        try {
            localStorage.setItem("wishlistItem", JSON.stringify([]))
            itemsInWishlist.value = []
        } catch (e) {
            console.error("Failed to clear wishlist", e)
        }
    }

    // Note: Initialization is handled in App.vue

    return {
        itemsInWishlist,
        isLoadingWishlist,
        isInitialized,
        isItemInWishlist,
        toggleWishlist,
        initFromBackend,
        addToWishlistStore: addToWishlist,
        clearWishlist,
        removeFromWishlistStore: removeFromWishlist,
        initFromLocalStorage,
        saveToLocalStorage,
    }
})