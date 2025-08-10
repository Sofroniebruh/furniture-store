import {defineStore} from "pinia";
import {ref} from "vue";
import {Product} from "@/lib/types";

export const useWishlistStore = defineStore("wishlistStore", () => {
    const itemsInWishlist = ref<Product[]>([])

    const isItemInWishlist = (id: string): boolean => {
        return itemsInWishlist.value.some((item) => item.id === id)
    }

    const addToWishlist = (product: Product) => {
        if (!isItemInWishlist(product.id)) {
            itemsInWishlist.value = [...itemsInWishlist.value, product]
            console.log("Item: " + product + " was added")
        } else {
            console.warn("Item: " + product + " was already added")
        }
    }
    const removeFromWishlist = (productOrId: Product | string) => {
        const productId = typeof productOrId === "string" ? productOrId : productOrId.id
        if (isItemInWishlist(productId)) {
            itemsInWishlist.value = itemsInWishlist.value.filter((item) => item.id !== productId)
            console.log("Item with ID: " + productId + " was removed")
        } else {
            console.warn("Item with ID: " + productId + " was not found")
        }
    }

    const toggleWishlist = (product: Product) => {
        if (!isItemInWishlist(product.id)) {
            addToWishlist(product)
        } else {
            removeFromWishlist(product)
        }
    }

    const initFromLocalStorage = () => {
        try {
            const stored = localStorage.getItem("wishlistItem")
            if (stored) {
                itemsInWishlist.value = JSON.parse(stored)
                console.log("Wishlist was initialized: ", stored)
            }
        } catch (e) {
            console.error("Failed to initialize wishlist", e)
        }
    }

    const saveToLocalStorage = () => {
        try {
            localStorage.setItem("wishlistItem", JSON.stringify(itemsInWishlist.value))
        } catch (e) {
            console.error("Failed to save wishlist item to local storage: ", e)
        }
    }

    return {
        itemsInWishlist,
        isItemInWishlist,
        toggleWishlist,
        addToWishlistStore: addToWishlist,
        removeFromWishlistStore: removeFromWishlist,
        initFromLocalStorage,
        saveToLocalStorage,
    }
})