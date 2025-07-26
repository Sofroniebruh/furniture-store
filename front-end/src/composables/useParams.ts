import {useRoute, useRouter} from "vue-router";

export const useParams = () => {
    const route = useRoute()
    const router = useRouter()

    const updatePage = async (page: number) => {
        await router.push({
            query: {
                ...route.query,
                page: page.toString(),
            }
        })
    }
    const updatePriceRange = async (priceFrom: number, priceTo: number) => {
        await router.push({
            query: {
                ...route.query,
                price_from: priceFrom.toString(),
                price_to: priceTo.toString(),
            }
        })
    }
    const updateEvent = async (event: string) => {
        await router.push({
            query: {
                ...route.query,
                event: event,
            }
        })
    }
    const updateModel = async (models: []) => {
        await router.push({
            query: {
                ...route.query,
                model: models,
            }
        })
    }
    const updateSorting = async (sortingType: string) => {
        await router.push({
            query: {
                ...route.query,
                sorting: sortingType,
            }
        })
    }

    return {
        updateEvent,
        updatePage,
        updatePriceRange,
        updateModel,
        updateSorting,
    }
}