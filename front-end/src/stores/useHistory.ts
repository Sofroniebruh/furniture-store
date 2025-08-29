import {defineStore} from "pinia";
import {ref} from "vue";
import {Product} from "@/lib/types";

export const useHistoryStore = defineStore("historyStore", () => {
    const itemsInHistory = ref<Product[]>([])
    const isLoadingHistory = ref<boolean>(false)
    const isInitializedHistory = ref<boolean>(false)

    const isItemInHistory = (id: string): boolean => {
        return itemsInHistory.value.some((item) => item.id === id)
    }

    const initHistoryFromBackend = (values: Product[]) => {
        itemsInHistory.value = values
        isInitializedHistory.value = true
        if (itemsInHistory.value.length !== 0) {
            saveToLocalStorage()
        }
    }

    const addToHistory = (product: Product) => {
        if (!isItemInHistory(product.id)) {
            itemsInHistory.value = [...itemsInHistory.value, product]
            saveToLocalStorage()
            console.log("Item: " + product.name + " was added to history")
        } else {
            console.warn("Item: " + product.name + " was already in history")
        }
    }

    const removeFromHistory = async (productOrId: Product | string) => {
        const productId = typeof productOrId === "string" ? productOrId : productOrId.id
        if (isItemInHistory(productId)) {
            itemsInHistory.value = itemsInHistory.value.filter((item) => item.id !== productId)
            saveToLocalStorage()
            console.log("Item: " + productId + " was removed from history")
        } else {
            console.warn("Item with ID: " + productId + " was not found in history")
        }
    }

    const toggleHistory = async (product: Product) => {
        if (!isItemInHistory(product.id)) {
            addToHistory(product)
        } else {
            await removeFromHistory(product)
        }
    }

    const initFromLocalStorage = () => {
        try {
            const stored = localStorage.getItem("historyItem")
            if (stored && stored !== "{}") {
                const parsedData = JSON.parse(stored)
                if (Array.isArray(parsedData)) {
                    itemsInHistory.value = parsedData
                    console.log("History initialized from localStorage: ", parsedData)
                }
            }
            isInitializedHistory.value = true
        } catch (e) {
            console.error("Failed to initialize history from localStorage", e)
            itemsInHistory.value = []
            isInitializedHistory.value = true
        }
    }

    const saveToLocalStorage = () => {
        try {
            localStorage.setItem("historyItem", JSON.stringify(itemsInHistory.value))
        } catch (e) {
            console.error("Failed to save history item to local storage: ", e)
        }
    }

    const clearHistory = () => {
        try {
            localStorage.setItem("historyItem", JSON.stringify([]))
            itemsInHistory.value = []
        } catch (e) {
            console.error("Failed to clear history", e)
        }
    }

    return {
        itemsInHistory,
        isLoadingHistory,
        isInitializedHistory,
        isItemInHistory,
        toggleHistory,
        initHistoryFromBackend,
        addToHistoryStore: addToHistory,
        clearHistory,
        removeFromHistoryStore: removeFromHistory,
        initFromLocalStorage,
        saveToLocalStorage,
    }
})