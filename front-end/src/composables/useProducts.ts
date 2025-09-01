import {computed, ref, watch} from "vue";
import {Product} from "@/lib/types";
import {useRoute} from "vue-router";
import {useParams} from "@/composables/useParams";

export function useProducts() {
    const route = useRoute()
    const limit = import.meta.env.VITE_PRODUCTS_PER_PAGE_LIMIT
    const backendUrl = import.meta.env.VITE_BACKEND_URL
    const products = ref<Product[]>([])
    const product = ref<Product | null>(null)
    const error = ref<boolean>(false)
    const loading = ref<boolean>(false)
    const totalPages = ref<Number>(0)

    const {updatePage} = useParams()

    const filters = computed(() => ({
        page: Number(route.query.page) || 1,
        priceFrom: Number(route.query.price_from) || null,
        priceTo: Number(route.query.price_to) || null,
        event: route.query.event || null,
        model: route.query.model || null,
        sorting: route.query.sorting || null,
        productName: route.query.product_name || null,
    }))

    const setPage = async (pageNumber: number) => {
        updatePage(pageNumber)
    }

    const fetchProducts = async () => {
        const queryParams = new URLSearchParams()
        loading.value = true
        error.value = false

        queryParams.set("page", filters.value.page.toString())
        queryParams.set("limit", limit)
        if (filters.value.productName) queryParams.set("product_name", filters.value.productName.toString())
        if (filters.value.priceFrom) queryParams.set("price_from", filters.value.priceFrom.toString())
        if (filters.value.priceTo) queryParams.set("price_to", filters.value.priceTo.toString())
        if (filters.value.event) queryParams.set("event", filters.value.event.toString())
        if (filters.value.model) queryParams.set("model", filters.value.model.toString())
        if (filters.value.sorting) queryParams.set("sorting", filters.value.sorting.toString())

        const baseUrl = computed(() => `${backendUrl}/products?${queryParams.toString()}`)

        try {
            const res = await fetch(baseUrl.value, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                }
            })

            if (!res.ok) {
                throw new Error(`Failed to fetch ${backendUrl}. Status: ` + res.status)
            }

            const data = await res.json() as { products: Product[], totalPages: Number }
            console.log(data)

            products.value = data.products
            totalPages.value = data.totalPages
        } catch (e) {
            if (e instanceof Error) {
                error.value = true
                console.error("Failed to fetch: ", e.message)
            }
            console.error(e)
        } finally {
            loading.value = false
        }
    }
    const fetchProductById = async (id: string) => {
        try {
            loading.value = true
            error.value = false
            const res = await fetch(`${import.meta.env.VITE_BACKEND_URL}/products/${id}`, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                }
            })

            if (!res.ok) {
                product.value = null
                error.value = true
                console.error(`Failed to fetch ${backendUrl}. Status: ` + res.status)
                return {status: res.status}
            }

            const data = await res.json() as { product: Product }

            product.value = data.product
            return {status: 200}
        } catch (e) {
            error.value = true
            console.error(e)
            return {status: 500}
        } finally {
            loading.value = false
        }
    }

    const handleNext = async (pageNumber: number) => {
        await setPage(pageNumber)
    }

    const handlePrevious = async (pageNumber: number) => {
        await setPage(pageNumber)
    }

    watch(() => route.query, async () => {
        await fetchProducts()
    }, {immediate: true, deep: true})

    return {
        products,
        product,
        error,
        loading,
        totalPages,
        currentPage: computed(() => filters.value.page),
        handleNext,
        handlePrevious,
        fetchProducts,
        fetchProductById,
        setPage,
    }
}